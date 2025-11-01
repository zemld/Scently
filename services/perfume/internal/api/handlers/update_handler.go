package handlers

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
)

// @description Update database state about perfumes. Can truncate table
// @tags Perfumes
// @summary Update perfumes
// @accept json
// @param hard query boolean false "Is hard update"
// @param password query string false "Password for hard update"
// @param perfumes body PerfumeCollection true "List of perfumes to update"
// @success 200 {object} UpdateResponse "Update successful"
// @failure 400 {object} UpdateResponse "Wrong perfumes format"
// @failure 500 {object} UpdateResponse "Something went wrong while update perfumes state"
// @router /update [post]
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
