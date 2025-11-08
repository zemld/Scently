package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
)

type SuggestResponse struct {
	Suggested []models.RankedPerfumeWithProps `json:"suggested"`
}

func WriteResponse(w http.ResponseWriter, response any, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
