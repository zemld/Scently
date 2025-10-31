package handlers

import (
	"context"
	"math"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/app"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

type suggestionsContextKey string

const SuggestionsContextKey suggestionsContextKey = "suggestions"

// @description Get suggests for perfumes. Accept brand and name and recommends 4+- perfumes which user probably will like.
// @tags Perfumes
// @summary Suggests some perfumes
// @produce json
// @param brand query string true "Brand of the perfume which you like"
// @param name query string true "Name of the perfume which you like"
// @param sex query string false "For her or for him"
// @param use_ai query boolean false "Use AI to suggest perfumes"
// @success 200 {object} SuggestResponse "Suggested perfumes"
// @success 204 "No content"
// @failure 400 "Incorrect parameters"
// @failure 403 "Forbidden"
// @failure 500
// @router /perfume [get]
func Suggest(w http.ResponseWriter, r *http.Request) {
	var suggestResponse SuggestResponse
	params := r.Context().Value(parameters.ParamsKey).(parameters.RequestPerfume)

	var mostSimilar []models.GluedPerfumeWithScore
	if params.UseAI {
		mostSimilar = app.GetAIEnrichedSuggestions(r.Context(), params)
	}
	if mostSimilar == nil {
		mostSimilar = app.GetComparisionSuggestions(r.Context(), params)
	}

	suggestResponse.Suggested = rankSuggestions(mostSimilar)
	status := http.StatusNoContent
	if len(mostSimilar) > 0 {
		status = http.StatusOK
	}
	WriteResponse(w, suggestResponse, status)

	ctxWithSuggestions := context.WithValue(r.Context(), SuggestionsContextKey, suggestResponse)
	*r = *r.WithContext(ctxWithSuggestions)
}

func rankSuggestions(suggestions []models.GluedPerfumeWithScore) []models.RankedPerfumeWithProps {
	rankedSuggestions := make([]models.RankedPerfumeWithProps, 0, len(suggestions))
	for i, suggestion := range suggestions {
		rankedSuggestions = append(
			rankedSuggestions,
			models.RankedPerfumeWithProps{
				Rank:    i + 1,
				Perfume: suggestion.GluedPerfume,
				Score:   math.Round(suggestion.Score*100) / 100,
			})
	}
	return rankedSuggestions
}
