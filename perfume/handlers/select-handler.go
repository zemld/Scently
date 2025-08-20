package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/handlers/responses"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

// @description Get info about perfumes. Can accept brand and name parameters
// @tags Perfumes
// @summary Get list of perfumes
// @produce json
// @param brand query string false "Brand of the perfume"
// @param name query string false "Name of the perfume"
// @success 200 {object} responses.PerfumeCollection "Found perfumes"
// @success 204 "No perfumes found"
// @router /get [get]
func SelectHandler(w http.ResponseWriter, r *http.Request) {
	p := getSelectionParameters(r)
	perfumes := core.Select(p)
	log.Printf("Found perfumes: %d\n", len(perfumes))
	if len(perfumes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeResponseWithFoundPerfumes(w, perfumes)
}

func getSelectionParameters(r *http.Request) *core.SelectParameters {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")

	return core.NewSelectParameters().WithBrand(brand).WithName(name)
}

func writeResponseWithFoundPerfumes(w http.ResponseWriter, perfumes []models.Perfume) {
	w.Header().Set("Content-Type", "application/json")
	perfumeResponse := responses.PerfumeCollection{Perfumes: perfumes}
	encodedPerfumes, _ := json.Marshal(perfumeResponse)
	w.Write(encodedPerfumes)
}
