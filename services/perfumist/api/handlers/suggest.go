package handlers

import (
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/advising"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/matching"
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	log.Println("Suggest request received")
	ctx := r.Context()

	params, err := generalParseSimilarParameters(r)
	if err != nil {
		handleError(w, err)
		return
	}

	perfumeHubFetcher, err := createPerfumeHubFetcher(config.Manager())
	if err != nil {
		handleError(w, err)
		return
	}
	_ = perfumeHubFetcher

	advisor := advising.NewBase(
		perfumeHubFetcher,
		matching.NewCombinedMatcher(
			*matching.NewWeights(
				config.Manager().GetFloatWithDefault("family_weight", 0.4),
				config.Manager().GetFloatWithDefault("notes_weight", 0.55),
				config.Manager().GetFloatWithDefault("type_weight", 0.05),
				config.Manager().GetFloatWithDefault("upper_notes_weight", 0.2),
				config.Manager().GetFloatWithDefault("core_notes_weight", 0.35),
				config.Manager().GetFloatWithDefault("base_notes_weight", 0.45),
				config.Manager().GetFloatWithDefault("characteristics_weight", 0.3),
				config.Manager().GetFloatWithDefault("tags_weight", 0.5),
				config.Manager().GetFloatWithDefault("overlay_weight", 0.2),
			),
		),
		config.Manager(),
	)

	suggested, err := advisor.Advise(ctx, params)
	if err != nil {
		handleError(w, err)
		return
	}

	WriteResponse(w, SuggestResponse{Suggested: suggested}, http.StatusOK)
}
