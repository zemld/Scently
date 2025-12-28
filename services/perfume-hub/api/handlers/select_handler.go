package handlers

import (
	"log"
	"net/http"

	"github.com/zemld/Scently/perfume-hub/internal/db/core"
	"github.com/zemld/Scently/perfume-hub/internal/errors"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

func Select(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")
	sex := r.URL.Query().Get("sex")

	params := models.NewSelectParameters().WithBrand(brand).WithName(name).WithSex(sex)

	perfumes, status := core.Select(r.Context(), params)
	if status.Error != nil {
		handleError(w, status.Error)
		return
	}

	response := PerfumeResponse{Perfumes: perfumes, State: status}
	log.Printf("Found perfumes: %d\n", len(perfumes))
	if len(perfumes) == 0 {
		WriteResponse(w, http.StatusNotFound, response)
		return
	}
	WriteResponse(w, http.StatusOK, response)
}

func handleError(w http.ResponseWriter, err error) {
	serviceErr, ok := err.(errors.ServiceError)
	if !ok {
		WriteResponse(w, http.StatusInternalServerError, models.NewProcessedState())
		return
	}

	status := serviceErr.HTTPStatus()
	WriteResponse(w, status, models.NewProcessedState())
}
