package middleware

import (
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/api/handlers"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/app"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

func Cache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, ok := r.Context().Value(parameters.ParamsKey).(parameters.RequestPerfume)
		if !ok {
			next(w, r)
			return
		}

		cached, err := app.LookupCache(r.Context(), params)
		if err == nil && cached != nil {
			handlers.WriteResponse(w, cached, http.StatusOK)
			return
		}

		next(w, r)

		suggestionsValue := r.Context().Value(handlers.SuggestionsContextKey)
		suggestions, _ := suggestionsValue.([]perfume.RankedWithProps)
		if err := app.Cache(r.Context(), params, suggestions); err != nil {
			log.Printf("Cannot cache: %v\n", err)
		}
	}
}
