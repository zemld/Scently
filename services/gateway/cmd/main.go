package main

import (
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/api/handlers"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/api/middleware"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /perfume/suggest", middleware.Cors(middleware.Cache(handlers.Suggest)))

	log.Printf("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
