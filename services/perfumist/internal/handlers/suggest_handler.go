package handlers

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/app"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/rdb"
)

const suggestsCount = 4

const (
	defaultTimeout = time.Second
	aiTimeout      = time.Second * 6
)

// @description Get suggests for perfumes. Accept brand and name and recommends 4+- perfumes which user probably will like.
// @tags Perfumes
// @summary Suggests some perfumes
// @produce json
// @param brand query string true "Brand of the perfume which you like"
// @param name query string true "Name of the perfume which you like"
// @param sex query string false "For her or for him"
// @param use_ai query boolean false "Use AI to suggest perfumes"
// @success 200 {object} SuggestResponse "Suggested perfumes"
// @success 204 {object} SuggestResponse "No perfumes found for suggestion"
// @failure 400 {object} SuggestResponse "Incorrect parameters"
// @failure 403 {object} SuggestResponse "Forbidden"
// @failure 500 {object} SuggestResponse
// @router /perfume [get]
func Suggest(w http.ResponseWriter, r *http.Request) {
	params := parseQuery(r)
	log.Printf("params: %+v", params)
	var suggestResponse SuggestResponse
	setInputToResponse(&suggestResponse, params)
	if !isValidQuery(params) {
		WriteResponse(w, suggestResponse, http.StatusBadRequest)
		return
	}

	var timeout time.Duration = defaultTimeout
	if params.UseAI {
		timeout = aiTimeout
	}
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	suggests, err := app.LookupCache(ctx, rdb.PerfumeCacheKey{
		Brand:      params.Brand,
		Name:       params.Name,
		AdviseType: suggestResponse.Input.AdviseType,
		Sex:        params.Sex,
	})
	if err == nil && suggests != nil {
		suggestResponse.Suggested = suggests
		suggestResponse.Success = true
		WriteResponse(w, suggestResponse, http.StatusOK)
		return
	}

	var mostSimilar []models.GluedPerfumeWithScore
	if params.UseAI {
		aiSuggests, err := app.AISuggest(ctx, params)
		if err == nil && aiSuggests != nil {
			mostSimilar = aiSuggests
		}
		if mostSimilar != nil {
			enrichmentParams := make([]parameters.RequestPerfume, len(mostSimilar))
			for i, suggestion := range mostSimilar {
				enrichmentParams[i] = *parameters.NewGet().WithBrand(suggestion.GluedPerfume.Brand).WithName(suggestion.GluedPerfume.Name).WithSex(params.Sex)
			}
			enrichedSuggests, ok := app.FetchPerfumes(ctx, enrichmentParams)
			if ok && enrichedSuggests != nil {
				for i, suggestion := range enrichedSuggests {
					mostSimilar[i].GluedPerfume = suggestion
				}
			}
		}
	}
	if mostSimilar == nil {
		favouritePerfumes, ok := app.FetchPerfumes(ctx, []parameters.RequestPerfume{params})
		if !ok || favouritePerfumes == nil || len(favouritePerfumes) == 0 {
			suggestResponse.Success = false
			WriteResponse(w, suggestResponse, http.StatusNoContent)
			return
		}
		allPerfumes, ok := app.FetchPerfumes(ctx, []parameters.RequestPerfume{*parameters.NewGet().WithSex(params.Sex)})
		if !ok {
			suggestResponse.Success = false
			WriteResponse(w, suggestResponse, http.StatusNoContent)
			return
		}

		mostSimilar = app.FoundSimilarities(favouritePerfumes[0], allPerfumes, suggestsCount)
	}

	fillResponseWithSuggestions(&suggestResponse, mostSimilar)
	if suggestResponse.Success {
		WriteResponse(w, suggestResponse, http.StatusOK)
	} else {
		WriteResponse(w, suggestResponse, http.StatusNoContent)
	}

	if !suggestResponse.Success || len(suggestResponse.Suggested) == 0 {
		return
	}
	err = app.Cache(
		ctx,
		rdb.PerfumeCacheKey{Brand: params.Brand, Name: params.Name, AdviseType: suggestResponse.Input.AdviseType},
		suggestResponse.Suggested,
	)
	if err != nil {
		log.Printf("Cannot cache: %v\n", err)
	}
}

func parseQuery(r *http.Request) parameters.RequestPerfume {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")
	sex := r.URL.Query().Get("sex")
	useAI := r.URL.Query().Get("use_ai")
	useAIBool, err := strconv.ParseBool(useAI)
	if err != nil {
		useAIBool = false
	}
	return *parameters.NewGet().WithBrand(brand).WithName(name).WithSex(sex).WithUseAI(useAIBool)
}

func isValidQuery(params parameters.RequestPerfume) bool {
	return params.Brand != "" && params.Name != ""
}

func setInputToResponse(response *SuggestResponse, params parameters.RequestPerfume) {
	response.Input = inputPerfume{Brand: params.Brand, Name: params.Name, Ok: true, Sex: params.Sex}
	if params.UseAI {
		response.Input.AdviseType = "AI"
	} else {
		response.Input.AdviseType = "Comparision"
	}
	if !isValidQuery(params) {
		response.Input.Ok = false
		response.Success = false
	}
}

func fillResponseWithSuggestions(response *SuggestResponse, suggestions []models.GluedPerfumeWithScore) {
	for i, suggestion := range suggestions {
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
