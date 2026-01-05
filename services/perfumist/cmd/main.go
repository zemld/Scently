package main

import (
	"log"
	"net/http"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/api/handlers"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/api/middleware"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
)

func main() {
	config.Manager().StartLoading(10 * time.Second)
	defer config.Manager().StopLoading()

	r := http.NewServeMux()

	r.HandleFunc("GET /v1/perfume/suggest", middleware.Auth(handlers.Suggest))

	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
