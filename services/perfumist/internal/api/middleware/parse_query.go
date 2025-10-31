package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

func ParseAndValidateQuery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		brand := r.URL.Query().Get(parameters.BrandParamKey)
		name := r.URL.Query().Get(parameters.NameParamKey)
		sex := r.URL.Query().Get(parameters.SexParamKey)
		useAI := r.URL.Query().Get(parameters.UseAIParamKey)
		useAIBool, err := strconv.ParseBool(useAI)
		if err != nil {
			useAIBool = false
		}
		if brand == "" || name == "" {
			http.Error(w, "Brand and name are required", http.StatusBadRequest)
			return
		}
		if sex != parameters.SexMale && sex != parameters.SexFemale {
			sex = parameters.SexUnisex
		}

		params := parameters.RequestPerfume{
			Brand: brand,
			Name:  name,
			Sex:   sex,
			UseAI: useAIBool,
		}

		ctx := context.WithValue(r.Context(), parameters.ParamsKey, params)
		next(w, r.WithContext(ctx))
	}
}
