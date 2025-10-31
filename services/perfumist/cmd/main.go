package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/handlers"
)

// @title Perfume Suggestion Service
// @version 1.0
// @description Mircoservice with bussiness logic for perfume suggestions
// @BasePath /v1/suggest
func main() {
	r := chi.NewRouter()

	r.Route("/v1/suggest", func(r chi.Router) {
		r.Get("/perfume", handlers.Suggest)
	})

	fs := http.FileServer(http.Dir("./docs"))
	r.Handle("/docs/*", http.StripPrefix("/docs/", fs))
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/swagger.json")))

	http.ListenAndServe(":8000", r)
}
