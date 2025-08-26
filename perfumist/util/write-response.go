package util

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, response any, status int) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
