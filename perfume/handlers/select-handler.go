package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/handlers/responses"
)

// @description Get info about perfumes. Can accept brand and name parameters
// @tags Perfumes
// @summary Get list of perfumes
// @produce json
// @param brand query string false "Brand of the perfume"
// @param name query string false "Name of the perfume"
// @success 200 {object} responses.PerfumeResponse "Found perfumes"
// @success 204 "No perfumes found"
// @failure 500 "Something went wrong while processing request"
// @router /get [get]
func SelectHandler(w http.ResponseWriter, r *http.Request) {
	p := getSelectionParameters(r)
	perfumes, status := core.Select(p)
	response := responses.PerfumeResponse{Perfumes: perfumes, State: status}
	if !status.Success {
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}
	log.Printf("Found perfumes: %d\n", len(perfumes))
	if len(perfumes) == 0 {
		writeResponse(w, http.StatusNoContent, response)
		return
	}
	writeResponse(w, http.StatusOK, response)
}

func getSelectionParameters(r *http.Request) *core.SelectParameters {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")

	return core.NewSelectParameters().WithBrand(brand).WithName(name)
}

func writeResponse(w http.ResponseWriter, code int, body responses.PerfumeResponse) {
	w.WriteHeader(code)
	writeResponseBody(w, body)
}

func writeResponseBody(w http.ResponseWriter, body responses.PerfumeResponse) {
	w.Header().Set("Content-Type", "application/json")
	encodedPerfumes, _ := json.Marshal(body)
	w.Write(encodedPerfumes)
}
