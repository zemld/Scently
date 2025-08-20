package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/handlers"
)

func main() {
	core.Initiate()
	r := chi.NewRouter()

	r.Route("/v1/perfumes", func(r chi.Router) {
		r.Get("/get", handlers.SelectHandler)
		r.Post("/update", handlers.UpdateHandler)
	})

	http.ListenAndServe(":8089", r)
}
