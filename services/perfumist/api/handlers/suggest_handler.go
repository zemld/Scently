package handlers

import (
	"context"
	"math"
	"net/http"
	"os"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/app"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/advising"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/matching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type suggestionsContextKey string

const SuggestionsContextKey suggestionsContextKey = "suggestions"

const aiSuggestUrl = "http://ai_advisor:8000/v1/advise"

const (
	getPerfumesUrl          = "http://perfume:8000/v1/perfumes/get"
	perfumeInternalTokenEnv = "PERFUME_INTERNAL_TOKEN"
)

const (
	familyWeight = 0.4
	notesWeight  = 0.55
	typeWeight   = 0.05
)

const (
	upperNotesWeight  = 0.15
	middleNotesWeight = 0.45
	baseNotesWeight   = 0.4
)

const (
	threadsCount = 5
)

const (
	suggestCount = 4
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	var suggestResponse SuggestResponse
	params := r.Context().Value(parameters.ParamsKey).(parameters.RequestPerfume)

	var advisor advising.Advisor
	dbFetcher := fetching.NewDB(getPerfumesUrl, os.Getenv(perfumeInternalTokenEnv))
	if params.UseAI {
		advisor = advising.NewAI(fetching.NewAI(aiSuggestUrl), dbFetcher)
	} else {
		advisor = advising.NewBase(dbFetcher, matching.NewOverlay(
			familyWeight,
			notesWeight,
			typeWeight,
			upperNotesWeight,
			middleNotesWeight,
			baseNotesWeight,
			threadsCount,
		), suggestCount)
	}
	suggested, err := advisor.Advise(params)
	if err != nil {
		WriteResponse(w, err.Error(), http.StatusInternalServerError)
	}

	suggestResponse.Suggested = suggested
	status := http.StatusNoContent
	if len(suggested) > 0 {
		status = http.StatusOK
	}
	WriteResponse(w, suggestResponse, status)
}

func OldSuggest(w http.ResponseWriter, r *http.Request) {
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
