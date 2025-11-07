package handlers

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

func Update(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(models.UpdateParametersContextKey).(*models.UpdateParameters)
	if len(params.Perfumes) == 0 {
		WriteResponse(w, http.StatusBadRequest, models.NewProcessedState())
		return
	}
	updateStatus := core.Update(r.Context(), params)
	if !updateStatus.Success {
		WriteResponse(w, http.StatusInternalServerError, updateStatus)
		return
	}
	WriteResponse(w, http.StatusOK, updateStatus)
}
