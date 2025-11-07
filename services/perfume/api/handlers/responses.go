package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

type PerfumeResponse struct {
	Perfumes []models.Perfume    `json:"perfumes"`
	State    core.ProcessedState `json:"state"`
}

type UpdateResponse struct {
	SuccessfulPerfumes []models.Perfume    `json:"successful_perfumes"`
	FailedPerfumes     []models.Perfume    `json:"failed_perfumes"`
	State              core.ProcessedState `json:"state"`
}

func WriteResponse(w http.ResponseWriter, code int, body any) {
	w.WriteHeader(code)
	writeResponseBody(w, body)
}

func writeResponseBody(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	encodedPerfumes, _ := json.Marshal(body)
	w.Write(encodedPerfumes)
}
