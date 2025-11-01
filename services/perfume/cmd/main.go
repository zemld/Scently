package main

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/api/handlers"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/api/middleware"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
)

// @title Perfume DB API
// @version 1.0
// @description Microservice which operates over perfume database
// @BasePath /v1/perfumes
func main() {
	core.Initiate()
	r := http.NewServeMux()

	r.Handle("/v1/perfumes/get", middleware.Auth(middleware.ParseQuery(http.HandlerFunc(handlers.Select))))
	r.Handle("/v1/perfumes/update", middleware.Auth(middleware.ParseQuery(http.HandlerFunc(handlers.Update))))

	fs := http.FileServer(http.Dir("./docs"))
	r.Handle("/docs/", http.StripPrefix("/docs/", fs))
	r.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/swagger.json")))

	http.ListenAndServe(":8000", r)
}
