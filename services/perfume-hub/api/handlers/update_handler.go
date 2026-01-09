package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/zemld/Scently/perfume-hub/internal/db/core"
	"github.com/zemld/Scently/perfume-hub/internal/errors"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

func Update(w http.ResponseWriter, r *http.Request) {
	params := models.NewUpdateParameters()

	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		validationErr := errors.NewValidationError("failed to read request body: " + err.Error())
		handleError(w, validationErr)
		return
	}

	if len(content) == 0 {
		validationErr := errors.NewValidationError("request body is empty")
		handleError(w, validationErr)
		return
	}

	if err := json.Unmarshal(content, params); err != nil {
		log.Printf("Error unmarshaling JSON: %v\n", err)
		validationErr := errors.NewValidationError("invalid JSON: " + err.Error())
		handleError(w, validationErr)
		return
	}

	if len(params.Perfumes) == 0 {
		validationErr := errors.NewValidationError("perfumes array is empty")
		handleError(w, validationErr)
		return
	}

	isHard, _ := strconv.ParseBool(r.URL.Query().Get("is_hard"))
	params.WithIsHard(isHard)

	updateStatus := core.Update(r.Context(), params)
	if updateStatus.Error != nil {
		handleError(w, updateStatus.Error)
		return
	}

	WriteResponse(w, http.StatusOK, updateStatus)
}
