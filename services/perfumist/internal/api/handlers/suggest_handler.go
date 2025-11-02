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
