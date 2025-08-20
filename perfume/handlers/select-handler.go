package handlers

import (
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/handlers/util"
)

// @description Get info about perfumes. Can accept brand and name parameters
// @tags Perfumes
// @summary Get list of perfumes
// @produce json
// @param brand query string false "Brand of the perfume"
// @param name query string false "Name of the perfume"
// @success 200 {object} util.PerfumeResponse "Found perfumes"
// @success 204 {object} util.PerfumeResponse "No perfumes found"
// @failure 500 {object} util.PerfumeResponse "Something went wrong while processing request"
// @router /get [get]
func SelectHandler(w http.ResponseWriter, r *http.Request) {
	p := getSelectionParameters(r)
	perfumes, status := core.Select(p)
	response := util.PerfumeResponse{Perfumes: perfumes, State: status}
	if !status.Success {
		util.WriteResponse(w, http.StatusInternalServerError, response)
		return
	}
	log.Printf("Found perfumes: %d\n", len(perfumes))
	if len(perfumes) == 0 {
		util.WriteResponse(w, http.StatusNoContent, response)
		return
	}
	util.WriteResponse(w, http.StatusOK, response)
}

func getSelectionParameters(r *http.Request) *core.SelectParameters {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")

	return core.NewSelectParameters().WithBrand(brand).WithName(name)
}
