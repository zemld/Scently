package middleware

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

func ParseQuery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brand := r.URL.Query().Get("brand")
		name := r.URL.Query().Get("name")
		sex := r.URL.Query().Get("sex")
		sp := core.NewSelectParameters().WithBrand(brand).WithName(name).WithSex(sex)

		up := models.NewUpdateParameters()
		if err := getPerfumesToUpdate(*r, up); err != nil {
			http.Error(w, "Failed to get perfumes", http.StatusBadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), core.SelectParametersContextKey, sp))
		r = r.WithContext(context.WithValue(r.Context(), models.UpdateParametersContextKey, up))
		next(w, r)
	}
}

func getPerfumesToUpdate(r http.Request, up *models.UpdateParameters) error {
	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		return err
	}

	json.Unmarshal(content, up)
	return nil
}
