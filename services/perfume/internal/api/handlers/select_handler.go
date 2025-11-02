package handlers

import (
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
)

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
