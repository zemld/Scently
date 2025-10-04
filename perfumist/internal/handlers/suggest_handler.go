package handlers

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/app"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/app/similarity"
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
func SuggestHandler(w http.ResponseWriter, r *http.Request) {
	var suggestResponse SuggestResponse
	params, ok := parseQuery(r, &suggestResponse)
	if !ok {
		WriteResponse(w, suggestResponse, http.StatusBadRequest)
		return
	}
	// TODO: эти запросы нужно посылать параллельно
	favouriteRawPerfumes, ok := app.GetPerfumes(params)
	if !ok {
		suggestResponse.Success = false
		WriteResponse(w, suggestResponse, http.StatusInternalServerError)
		return
	}
	favouritePerfume := app.Glue(favouriteRawPerfumes)[0]
	allRawPerfumes, ok := app.GetPerfumes(*util.NewGetParameters())
	if !ok {
		suggestResponse.Success = false
		WriteResponse(w, suggestResponse, http.StatusInternalServerError)
		return
	}

	allPerfumes := app.Glue(allRawPerfumes)
	mostSimilar := make([]gluedPerfumeWithScore, suggestsCount)
	for _, perfume := range allPerfumes {
		if favouritePerfume.Equal(perfume) {
			continue
		}
		similarityScore := similarity.GetPerfumeSimilarityScore(favouritePerfume.Properties, perfume.Properties)
		updateMostSimilarIfNeeded(mostSimilar, perfume, similarityScore)
	}
	fillResponseWithSuggestions(&suggestResponse, mostSimilar)
	if suggestResponse.Success {
		WriteResponse(w, suggestResponse, http.StatusOK)
	} else {
		WriteResponse(w, suggestResponse, http.StatusNoContent)
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
			rankedPerfumeWithProps{
				Rank:    i + 1,
				Perfume: suggestion.GluedPerfume,
			})
	}
	if len(response.Suggested) > 0 {
		response.Success = true
	} else {
		response.Success = false
	}
}
