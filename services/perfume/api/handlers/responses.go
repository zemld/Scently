package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

type PerfumeResponse struct {
	Perfumes []models.Perfume      `json:"perfumes"`
	State    models.ProcessedState `json:"state"`
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
