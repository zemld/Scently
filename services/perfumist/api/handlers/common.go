package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/Scently/models"
	"github.com/zemld/config-manager/pkg/cm"
)

type SuggestResponse struct {
	Suggested []models.Ranked `json:"suggested"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func generalParseSimilarParameters(r *http.Request) (parameters.RequestPerfume, error) {
	query := r.URL.Query()

	brand := query.Get(parameters.BrandParamKey)
	name := query.Get(parameters.NameParamKey)
	sex := query.Get(parameters.SexParamKey)

	if sex != parameters.SexMale && sex != parameters.SexFemale {
		sex = parameters.SexUnisex
	}
	params := parameters.NewGet().WithBrand(brand).WithName(name).WithSex(sex)
	return *params, params.Validate()
}

func createPerfumeHubFetcher(cm cm.ConfigManager) (fetching.Fetcher, error) {
	getPerfumesUrl, err := cm.GetString("get_perfumes_url")
	if err != nil {
		return nil, errors.NewServiceError("failed to get get_perfumes_url", err)
	}
	perfumeHubInternalTokenEnv, err := cm.GetString("perfume_hub_internal_token_env_name")
	if err != nil {
		return nil, errors.NewServiceError("failed to get perfume_hub_internal_token_env_name", err)
	}
	return fetching.NewPerfumeHub(getPerfumesUrl, os.Getenv(perfumeHubInternalTokenEnv), cm), nil
}

func createAIFetcher(cm cm.ConfigManager) fetching.Fetcher {
	return fetching.NewAI(
		os.Getenv("BASE_URL"),
		os.Getenv("FOLDER_ID"),
		os.Getenv("MODEL_NAME"),
		os.Getenv("API_KEY"),
		cm,
	)
}

func handleError(w http.ResponseWriter, err error) {
	var status int
	var errorMsg string

	switch e := err.(type) {
	case *errors.ValidationError:
		status = http.StatusBadRequest
		errorMsg = e.Error()
	case *errors.NotFoundError:
		status = http.StatusNotFound
		errorMsg = e.Error()
	case *errors.ServiceError:
		status = http.StatusInternalServerError
		errorMsg = e.Error()
	default:
		status = http.StatusInternalServerError
		errorMsg = "internal server error"
	}

	WriteResponse(w, ErrorResponse{Error: errorMsg}, status)
}

func WriteResponse(w http.ResponseWriter, response any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
