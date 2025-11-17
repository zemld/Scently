package main

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/api/handlers"
)

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /v1/perfume/suggest", handlers.Suggest)

	http.ListenAndServe(":8000", r)
}
