package util

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, code int, body any) {
	w.WriteHeader(code)
	writeResponseBody(w, body)
}

func writeResponseBody(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	encodedPerfumes, _ := json.Marshal(body)
	w.Write(encodedPerfumes)
}
