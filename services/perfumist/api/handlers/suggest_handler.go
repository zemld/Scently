package handlers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/advising"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/matching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/config-manager/pkg/cm"
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := parseQueryParams(r)
	if err := params.Validate(); err != nil {
		handleError(w, err)
		return
	}

	advisor := createAdvisor(params, config.Manager())

	if advisor == nil {
		handleError(w, errors.NewServiceError("failed to create advisor", nil))
		return
	}

	suggested, err := advisor.Advise(ctx, params)
	if err != nil {
		handleError(w, err)
		return
	}

	WriteResponse(w, SuggestResponse{Suggested: suggested}, http.StatusOK)
}

func parseQueryParams(r *http.Request) parameters.RequestPerfume {
	query := r.URL.Query()

	brand := query.Get(parameters.BrandParamKey)
	name := query.Get(parameters.NameParamKey)
	sex := query.Get(parameters.SexParamKey)
	useAIStr := query.Get(parameters.UseAIParamKey)

	useAI, err := strconv.ParseBool(useAIStr)
	if err != nil {
		useAI = false
	}

	if sex != parameters.SexMale && sex != parameters.SexFemale {
		sex = parameters.SexUnisex
	}

	return *parameters.NewGet().WithBrand(brand).WithName(name).WithSex(sex).WithUseAI(useAI)
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

func createAdvisor(params parameters.RequestPerfume, cm cm.ConfigManager) advising.Advisor {
	getPerfumesUrl, err := cm.GetString("get_perfumes_url")
	if err != nil {
		return nil
	}
	perfumeHubInternalTokenEnv, err := cm.GetString("perfume_hub_internal_token_env_name")
	if err != nil {
		return nil
	}

	fetcher := fetching.NewPerfumeHub(getPerfumesUrl, os.Getenv(perfumeHubInternalTokenEnv), cm)

	if params.UseAI {
		return advising.NewAI(
			fetching.NewAI(
				os.Getenv("BASE_URL"),
				os.Getenv("FOLDER_ID"),
				os.Getenv("MODEL_NAME"),
				os.Getenv("API_KEY"),
				cm,
			),
			fetcher,
			cm,
		)
	}

	return advising.NewBase(
		fetcher,
		matching.NewCombinedMatcher(
			*matching.NewWeights(
				cm.GetFloatWithDefault("family_weight", 0.4),
				cm.GetFloatWithDefault("notes_weight", 0.55),
				cm.GetFloatWithDefault("type_weight", 0.05),
				cm.GetFloatWithDefault("upper_notes_weight", 0.2),
				cm.GetFloatWithDefault("core_notes_weight", 0.35),
				cm.GetFloatWithDefault("base_notes_weight", 0.45),
				cm.GetFloatWithDefault("characteristics_weight", 0.3),
				cm.GetFloatWithDefault("tags_weight", 0.5),
				cm.GetFloatWithDefault("overlay_weight", 0.2),
			),
		),
		cm,
	)
}
