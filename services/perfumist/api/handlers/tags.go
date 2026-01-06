package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/zemld/Scently/perfumist/internal/config"
	"github.com/zemld/Scently/perfumist/internal/errors"
	"github.com/zemld/Scently/perfumist/internal/models/advising"
	"github.com/zemld/Scently/perfumist/internal/models/matching"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
)

func SuggestByTags(w http.ResponseWriter, r *http.Request) {
	log.Println("SuggestByTags request received")
	ctx := r.Context()

	sex := parseSexParameter(r)
	rawTags := r.URL.Query().Get("tags")
	tags := strings.Split(rawTags, ",")
	log.Printf("Tags: %+v, length: %d", tags, len(tags))
	if len(tags) == 0 || rawTags == "" {
		handleError(w, errors.NewValidationError("tags", "are required"))
		return
	}

	perfumeHubFetcher, err := createPerfumeHubFetcher(config.Manager())
	if err != nil {
		handleError(w, err)
		return
	}

	advisor := advising.NewTagsBased(
		*matching.NewTagsBasedAdapter(
			matching.Weights{
				UpperNotesWeight: config.Manager().GetFloatWithDefault("upper_notes_weight", 0.2),
				CoreNotesWeight:  config.Manager().GetFloatWithDefault("core_notes_weight", 0.35),
				BaseNotesWeight:  config.Manager().GetFloatWithDefault("base_notes_weight", 0.45),
			},
			matching.NewTagsBased(config.Manager().GetIntWithDefault("suggest_count", 4)),
			tags,
		),
		perfumeHubFetcher,
		config.Manager(),
	)

	suggested, err := advisor.Advise(ctx, *parameters.NewGet().WithSex(sex))
	if err != nil {
		handleError(w, err)
		return
	}
	WriteResponse(w, SuggestResponse{Suggested: suggested}, http.StatusOK)
}
