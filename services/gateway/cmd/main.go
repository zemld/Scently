package main

import (
	"log"
	"net/http"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/api/handlers"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/api/middleware"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/config"
)

func main() {
	config.Manager().StartLoading(10 * time.Second)
	defer config.Manager().StopLoading()

	router := http.NewServeMux()

	router.HandleFunc("GET /perfume/suggest", middleware.Cors(middleware.Cache(handlers.Suggest)))
	router.HandleFunc("GET /perfume/suggest-by-tags", handlers.SuggestByTags)

	log.Printf("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
