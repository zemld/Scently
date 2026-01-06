package main

import (
	"log"
	"net/http"
	"time"

	"github.com/zemld/Scently/perfumist/api/handlers"
	"github.com/zemld/Scently/perfumist/api/middleware"
	"github.com/zemld/Scently/perfumist/internal/config"
)

func main() {
	config.Manager().StartLoading(10 * time.Second)
	defer config.Manager().StopLoading()

	r := http.NewServeMux()

	r.HandleFunc("GET /v2/perfume/suggest", middleware.Auth(handlers.Suggest))
	r.HandleFunc("GET /v2/perfume/ai-suggest", middleware.Auth(handlers.AISuggest))
	r.HandleFunc("GET /v2/perfume/suggest-by-tags", handlers.SuggestByTags)

	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
