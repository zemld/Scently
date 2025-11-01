package handlers

import (
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
)

// @description Get info about perfumes. Can accept brand and name parameters
// @tags Perfumes
// @summary Get list of perfumes
// @produce json
// @param brand query string false "Brand of the perfume"
// @param name query string false "Name of the perfume"
// @param sex query string false "For him or for her"
// @success 200 {object} PerfumeResponse "Found perfumes"
// @success 204 {object} PerfumeResponse "No perfumes found"
// @failure 500 {object} PerfumeResponse "Something went wrong while processing request"
// @router /get [get]
func Select(w http.ResponseWriter, r *http.Request) {
	perfumes, status := core.Select(r.Context(), r.Context().Value(core.SelectParametersContextKey).(*core.SelectParameters))
	response := PerfumeResponse{Perfumes: perfumes, State: status}
	if !status.Success {
		WriteResponse(w, http.StatusInternalServerError, response)
		return
	}
	log.Printf("Found perfumes: %d\n", len(perfumes))
	if len(perfumes) == 0 {
		WriteResponse(w, http.StatusNoContent, response)
		return
	}
	WriteResponse(w, http.StatusOK, response)
}
