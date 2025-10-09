package handlers

import (
	"context"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/app"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/util"
)

const suggestsCount = 4

type gluedPerfumeWithScore struct {
	models.GluedPerfume
	Score float64
}

// @description Get suggests for perfumes. Accept brand and name and recommends 4+- perfumes which user probably will like.
// @tags Perfumes
// @summary Suggests some perfumes
// @produce json
// @param brand query string true "Brand of the perfume which you like"
// @param name query string true "Name of the perfume which you like"
// @success 200 {object} SuggestResponse "Suggested perfumes"
// @success 204 {object} SuggestResponse "No perfumes found for suggestion"
// @failure 400 {object} SuggestResponse "Incorrect parameters"
// @failure 500 {object} SuggestResponse
// @router /perfume [get]
func Suggest(w http.ResponseWriter, r *http.Request) {
	var suggestResponse SuggestResponse
	params, ok := parseQuery(r, &suggestResponse)
	if !ok {
		WriteResponse(w, suggestResponse, http.StatusBadRequest)
		return
	}

	suggests, err := app.LookupCache(models.Perfume{Brand: params.Brand, Name: params.Name})
	if err == nil && suggests != nil {
		suggestResponse.Suggested = suggests
		WriteResponse(w, suggestResponse, http.StatusOK)
		return
	}

	favouritePerfumesChan := make(chan perfumesFetchAndGlueResult)
	allPerfumesChan := make(chan perfumesFetchAndGlueResult)
	go getAndGluePerfumesAsync(params, favouritePerfumesChan)
	go getAndGluePerfumesAsync(*util.NewGetParameters(), allPerfumesChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	favoritePerfumes, allPerfumes, ok := fetchPerfumeResults(ctx, favouritePerfumesChan, allPerfumesChan)
	if !ok {
		suggestResponse.Success = false
		WriteResponse(w, suggestResponse, http.StatusInternalServerError)
		return
	}
	favouritePerfume := favoritePerfumes[0]

	mostSimilar := make([]gluedPerfumeWithScore, suggestsCount)
	for _, perfume := range allPerfumes {
		if favouritePerfume.Equal(perfume) {
			continue
		}
		similarityScore := app.GetPerfumeSimilarityScore(favouritePerfume.Properties, perfume.Properties)
		updateMostSimilarIfNeeded(mostSimilar, perfume, similarityScore)
	}

	fillResponseWithSuggestions(&suggestResponse, mostSimilar)
	if suggestResponse.Success {
		WriteResponse(w, suggestResponse, http.StatusOK)
	} else {
		WriteResponse(w, suggestResponse, http.StatusNoContent)
	}

	err = app.Cache(models.Perfume{Brand: params.Brand, Name: params.Name}, suggestResponse.Suggested)
	if err != nil {
		log.Printf("Cannot cache: %v\n", err)
	}
}

func parseQuery(r *http.Request, suggestResponse *SuggestResponse) (util.GetParameters, bool) {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")
	suggestResponse.Input = inputPerfume{Brand: brand, Name: name, Ok: true}
	suggestResponse.Success = true
	if brand == "" || name == "" {
		suggestResponse.Input.Ok = false
		suggestResponse.Success = false
		return util.GetParameters{}, false
	}
	return *util.NewGetParameters().WithBrand(brand).WithName(name), true
}

type perfumesFetchAndGlueResult struct {
	Perfumes []models.GluedPerfume
	Ok       bool
}

func getAndGluePerfumesAsync(params util.GetParameters, results chan<- perfumesFetchAndGlueResult) {
	defer close(results)
	perfumes, ok := app.GetPerfumes(params)
	if !ok {
		results <- perfumesFetchAndGlueResult{Perfumes: nil, Ok: false}
		return
	}
	results <- perfumesFetchAndGlueResult{Perfumes: app.Glue(perfumes), Ok: true}
}

func fetchPerfumeResults(ctx context.Context, favChan <-chan perfumesFetchAndGlueResult, allChan <-chan perfumesFetchAndGlueResult) ([]models.GluedPerfume, []models.GluedPerfume, bool) {
	var favs []models.GluedPerfume
	var all []models.GluedPerfume

	select {
	case favResult := <-favChan:
		favs = favResult.Perfumes
		select {
		case allResult := <-allChan:
			all = allResult.Perfumes
		case <-ctx.Done():
			return favs, all, false
		}
	case allResult := <-allChan:
		all = allResult.Perfumes
		select {
		case favResult := <-favChan:
			favs = favResult.Perfumes
		case <-ctx.Done():
			return favs, all, false
		}
	case <-ctx.Done():
		return favs, all, false
	}
	return favs, all, true
}

func updateMostSimilarIfNeeded(mostSimilar []gluedPerfumeWithScore, perfume models.GluedPerfume, similarityScore float64) {
	current := perfume
	for i := range mostSimilar {
		if similarityScore > mostSimilar[i].Score {
			tmp := mostSimilar[i]
			mostSimilar[i].Score = similarityScore
			mostSimilar[i].GluedPerfume = current
			current = tmp.GluedPerfume
			similarityScore = tmp.Score
		}
	}
}

func fillResponseWithSuggestions(response *SuggestResponse, suggestions []gluedPerfumeWithScore) {
	for i, suggestion := range suggestions {
		if suggestion.Score == 0 {
			break
		}
		response.Suggested = append(
			response.Suggested,
			models.RankedPerfumeWithProps{
				Rank:    i + 1,
				Perfume: suggestion.GluedPerfume,
				Score:   math.Round(suggestion.Score*100) / 100,
			})
	}
	if len(response.Suggested) > 0 {
		response.Success = true
	} else {
		response.Success = false
	}
}
