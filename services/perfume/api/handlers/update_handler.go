package handlers

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
)

func Update(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(core.UpdateParametersContextKey).(*core.UpdateParameters)
	if len(params.Perfumes) == 0 {
		WriteResponse(w, http.StatusBadRequest, UpdateResponse{State: core.NewProcessedState()})
		return
	}
	updateStatus := core.Update(r.Context(), params)
	if !updateStatus.State.Success {
		WriteResponse(w, http.StatusInternalServerError, updateStatus)
		return
	}
	WriteResponse(w, http.StatusOK, updateStatus)
}
