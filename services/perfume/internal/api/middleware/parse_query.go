package middleware

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
)

func ParseQuery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brand := r.URL.Query().Get("brand")
		name := r.URL.Query().Get("name")
		sex := r.URL.Query().Get("sex")
		sp := core.NewSelectParameters().WithBrand(brand).WithName(name).WithSex(sex)

		up := core.NewUpdateParameters()
		if err := getPerfumesToUpdate(*r, up); err != nil {
			http.Error(w, "Failed to get perfumes", http.StatusBadRequest)
			return
		}

		hardBool, err := strconv.ParseBool(r.URL.Query().Get("hard"))
		if err != nil {
			hardBool = false
		}
		if hardBool {
			up.WithTruncate()
		}

		r = r.WithContext(context.WithValue(r.Context(), core.SelectParametersContextKey, sp))
		r = r.WithContext(context.WithValue(r.Context(), core.UpdateParametersContextKey, up))
		next(w, r)
	}
}

func getPerfumesToUpdate(r http.Request, up *core.UpdateParameters) error {
	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		return err
	}

	if err = json.Unmarshal(content, up); err != nil {
		log.Printf("Error unmarshalling request body: %v\n", err)
		return err
	}
	return nil
}
