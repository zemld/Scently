package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

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
	p := getSelectionParameters(r)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	perfumes, status := core.Select(ctx, p)
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

func getSelectionParameters(r *http.Request) *core.SelectParameters {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")
	sex := r.URL.Query().Get("sex")

	return core.NewSelectParameters().WithBrand(brand).WithName(name).WithSex(sex)
}
