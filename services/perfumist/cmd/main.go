package main

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/api/handlers"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/api/middleware"
)

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /v1/suggest/perfume", middleware.Cors(middleware.ParseAndValidateQuery(middleware.Cache(handlers.Suggest))))

	http.ListenAndServe(":8000", r)
}
