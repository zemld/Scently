package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zemld/Scently/models"
)

type SuggestResponse struct {
	Suggested []models.Ranked `json:"suggested"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteResponse(w http.ResponseWriter, response any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
