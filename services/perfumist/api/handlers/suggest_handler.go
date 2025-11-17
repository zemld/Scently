package handlers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/advising"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/errors"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/matching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

const aiSuggestUrl = "http://ai_advisor:8000/v1/advise"

const (
	getPerfumesUrl          = "http://perfume:8000/v1/perfumes/get"
	perfumeInternalTokenEnv = "PERFUME_INTERNAL_TOKEN"
)

const (
	familyWeight = 0.4
	notesWeight  = 0.55
	typeWeight   = 0.05
)

const (
	upperNotesWeight  = 0.15
	middleNotesWeight = 0.45
	baseNotesWeight   = 0.4
)

const (
	threadsCount = 5
)

const (
	suggestCount = 4
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	params := parseQueryParams(r)

	if err := params.Validate(); err != nil {
		handleError(w, err)
		return
	}

	advisor := createAdvisor(params)

	suggested, err := advisor.Advise(params)
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

	useAI, _ := strconv.ParseBool(useAIStr)

	if sex != parameters.SexMale && sex != parameters.SexFemale {
		sex = parameters.SexUnisex
	}

	return parameters.RequestPerfume{
		Brand: brand,
		Name:  name,
		Sex:   sex,
		UseAI: useAI,
	}
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
func createAdvisor(params parameters.RequestPerfume) advising.Advisor {
	dbFetcher := fetching.NewDB(getPerfumesUrl, os.Getenv(perfumeInternalTokenEnv))

	if params.UseAI {
		return advising.NewAI(fetching.NewAI(aiSuggestUrl), dbFetcher)
	}

	return advising.NewBase(
		dbFetcher,
		matching.NewOverlay(
			familyWeight,
			notesWeight,
			typeWeight,
			upperNotesWeight,
			middleNotesWeight,
			baseNotesWeight,
			threadsCount,
		),
		suggestCount,
	)
}
