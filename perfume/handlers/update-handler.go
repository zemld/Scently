package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/handlers/util"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

// @description Update database state about perfumes. Can truncate table
// @tags Perfumes
// @summary Update perfumes
// @accept json
// @param hard query boolean false "Is hard update"
// @param password query string false "Password for hard update"
// @param perfumes body util.PerfumeCollection true "List of perfumes to update"
// @success 200 {object} util.UpdateResponse "Update successful"
// @failure 400 {object} util.UpdateResponse "Wrong perfumes format"
// @failure 500 {object} util.UpdateResponse "Something went wrong while update perfumes state"
// @router /update [post]
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	p := getUpdateParametersFromRequest(r)
	perfumes, err := getPerfumesToUpdate(r)
	if err != nil {
		http.Error(w, "Failed to get perfumes", http.StatusBadRequest)
		return
	}
	updateStatus := core.Update(p, perfumes)
	if !updateStatus.State.Success {
		util.WriteResponse(w, http.StatusInternalServerError, updateStatus)
		return
	}
	util.WriteResponse(w, http.StatusOK, updateStatus)
}

func getPerfumesToUpdate(r *http.Request) ([]models.Perfume, error) {
	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		return nil, err
	}
	var perfumes util.PerfumeCollection
	if err = json.Unmarshal(content, &perfumes); err != nil {
		log.Printf("Error unmarshalling request body: %v\n", err)
		return nil, err
	}
	return perfumes.Perfumes, nil
}
