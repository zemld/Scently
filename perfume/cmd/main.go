package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/handlers"
)

// @title Perfume DB API
// @version 1.0
// @description Microservice which operates over perfume database
// @BasePath /v1/perfumes
func main() {
	core.Initiate()
	r := chi.NewRouter()

	r.Route("/v1/perfumes", func(r chi.Router) {
		r.Get("/get", handlers.SelectHandler)
		// /v1/perfumes/update?hard:bool&password=string
		r.Post("/update", handlers.UpdateHandler)
	})

	fs := http.FileServer(http.Dir("./docs"))
	r.Handle("/docs/*", http.StripPrefix("/docs/", fs))
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8089/docs/swagger.json")))
	http.ListenAndServe(":8089", r)
}
