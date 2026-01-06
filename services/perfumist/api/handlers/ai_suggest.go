package handlers

import (
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/advising"
)

func AISuggest(w http.ResponseWriter, r *http.Request) {
	log.Println("AISuggest request received")
	ctx := r.Context()
	params, err := generalParseSimilarParameters(r)
	if err != nil {
		log.Printf("Error parsing parameters: %v\n", err)
		handleError(w, err)
		return
	}

	perfumeHubFetcher, err := createPerfumeHubFetcher(config.Manager())
	if err != nil {
		log.Printf("Error creating perfume hub fetcher: %v\n", err)
		handleError(w, err)
		return
	}
	aiFetcher := createAIFetcher(config.Manager())
	log.Printf("AI fetcher created: %+v\n", aiFetcher)
	aiAdvisor := advising.NewAI(
		aiFetcher,
		perfumeHubFetcher,
		config.Manager(),
	)

	suggested, err := aiAdvisor.Advise(ctx, params)
	if err != nil {
		log.Printf("Error advising: %v\n", err)
		handleError(w, err)
		return
	}

	WriteResponse(w, SuggestResponse{Suggested: suggested}, http.StatusOK)
}
