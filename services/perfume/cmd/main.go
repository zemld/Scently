package main

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/api/handlers"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/api/middleware"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
)

func main() {
	core.Initiate()
	r := http.NewServeMux()

	r.Handle("/v1/perfumes/get", middleware.Auth(middleware.ParseQuery(http.HandlerFunc(handlers.Select))))
	r.Handle("/v1/perfumes/update", middleware.Auth(middleware.ParseQuery(http.HandlerFunc(handlers.Update))))

	http.ListenAndServe(":8000", r)
}
