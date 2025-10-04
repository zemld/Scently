package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
)

type rankedPerfumeWithProps struct {
	Perfume models.GluedPerfume `json:"perfume"`
	Rank    int                 `json:"rank"`
}

type inputPerfume struct {
	Brand string `json:"brand"`
	Name  string `json:"name"`
	Ok    bool   `json:"ok"`
}

type SuggestResponse struct {
	Input     inputPerfume             `json:"input"`
	Suggested []rankedPerfumeWithProps `json:"suggested"`
	Success   bool                     `json:"success"`
}

func WriteResponse(w http.ResponseWriter, response any, status int) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
