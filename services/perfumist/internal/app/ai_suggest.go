package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

const (
	aiSuggestUrl = "http://ai_advisor:8000/v1/advise"
)

type aISuggestion struct {
	Perfumes []models.GluedPerfume `json:"perfumes"`
}

func AISuggest(ctx context.Context, params parameters.RequestPerfume) ([]models.GluedPerfumeWithScore, error) {
	r, _ := http.NewRequestWithContext(ctx, "GET", aiSuggestUrl, nil)
	updateQuery(r, params)

	response, err := http.DefaultClient.Do(r)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil || len(body) == 0 {
		return nil, err
	}
	var suggestions aISuggestion
	err = json.Unmarshal(body, &suggestions)
	if err != nil {
		return nil, err
	}
	var suggestionsWithScore []models.GluedPerfumeWithScore
	for _, suggestion := range suggestions.Perfumes {
		suggestionsWithScore = append(suggestionsWithScore, models.GluedPerfumeWithScore{
			GluedPerfume: suggestion,
		})
	}
	return suggestionsWithScore, nil
}
