package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/suggest/perfume", handlers.SuggestHandler)
	})

	http.ListenAndServe(":8088", r)
}
