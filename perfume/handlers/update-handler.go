package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/handlers/responses"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	p := getUpdateParametersFromRequest(r)
	perfumes, err := getPerfumesToUpdate(r)
	if err != nil {
		http.Error(w, "Failed to get perfumes", http.StatusBadRequest)
		return
	}
	core.Update(p, perfumes)
	w.WriteHeader(http.StatusOK)
}

func getPerfumesToUpdate(r *http.Request) ([]models.Perfume, error) {
	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		return nil, err
	}
	var perfumes responses.PerfumeCollection
	if err = json.Unmarshal(content, &perfumes); err != nil {
		log.Printf("Error unmarshalling request body: %v\n", err)
		return nil, err
	}
	return perfumes.Perfumes, nil
}
