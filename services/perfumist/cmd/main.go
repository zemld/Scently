package main

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/api/handlers"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/api/middleware"
)

// @title Perfume Suggestion Service
// @version 1.0
// @description Mircoservice with bussiness logic for perfume suggestions
// @BasePath /v1/suggest
func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /v1/suggest/perfume", middleware.Cors(middleware.ParseAndValidateQuery(middleware.Cache(handlers.Suggest))))

	fs := http.FileServer(http.Dir("./docs"))
	r.Handle("/docs/", http.StripPrefix("/docs/", fs))
	r.Handle("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8088/docs/swagger.json")))

	http.ListenAndServe(":8000", r)
}
