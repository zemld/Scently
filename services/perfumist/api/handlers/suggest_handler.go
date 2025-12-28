package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/advising"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/matching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	ph "github.com/zemld/Scently/generated/proto/perfume-hub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := parseQueryParams(r)

	if err := params.Validate(); err != nil {
		handleError(w, err)
		return
	}

	advisor := createAdvisor(params)

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

func createAdvisor(params parameters.RequestPerfume) advising.Advisor {
	getPerfumesURL := os.Getenv(config.GetPerfumesURLEnv)
	if getPerfumesURL == "" {
		getPerfumesURL = config.DefaultGetPerfumesURL
	}

	aiSuggestURL := os.Getenv(config.AISuggestURLEnv)
	if aiSuggestURL == "" {
		aiSuggestURL = config.DefaultAISuggestURL
	}

	var fetcher fetching.Fetcher

	conn, err := grpc.NewClient(
		os.Getenv(config.PerfumeHubURLEnv),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)),
	)
	if err != nil {
		log.Printf("failed to connect to perfume hub: %v", err)
		fetcher = fetching.NewDB(getPerfumesURL, os.Getenv(config.PerfumeInternalTokenEnv))
	} else {
		fetcher = fetching.NewPerfumeHub(ph.NewPerfumeStorageClient(conn))
	}

	if params.UseAI {
		return advising.NewAI(fetching.NewAI(aiSuggestURL), fetcher)
	}

	return advising.NewBase(
		fetcher,
		matching.NewOverlay(
			config.FamilyWeight,
			config.NotesWeight,
			config.TypeWeight,
			config.UpperNotesWeight,
			config.MiddleNotesWeight,
			config.BaseNotesWeight,
			config.ThreadsCount,
		),
		config.SuggestCount,
	)
}
