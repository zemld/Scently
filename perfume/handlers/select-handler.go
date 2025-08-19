package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/handlers/responses"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

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
