package handlers

import (
	"context"
	"math"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/app"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type suggestionsContextKey string

const SuggestionsContextKey suggestionsContextKey = "suggestions"

func Suggest(w http.ResponseWriter, r *http.Request) {
	var suggestResponse SuggestResponse
	params := r.Context().Value(parameters.ParamsKey).(parameters.RequestPerfume)

	var mostSimilar []perfume.WithScore
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

func rankSuggestions(suggestions []perfume.WithScore) []perfume.Ranked {
	rankedSuggestions := make([]perfume.Ranked, 0, len(suggestions))
	for i, suggestion := range suggestions {
		rankedSuggestions = append(
			rankedSuggestions,
			perfume.Ranked{
				Rank:    i + 1,
				Perfume: suggestion.Perfume,
				Score:   math.Round(suggestion.Score*100) / 100,
			})
	}
	return rankedSuggestions
}
