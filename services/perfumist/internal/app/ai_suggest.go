package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

const (
	aiSuggestUrl = "http://ai_advisor:8000/v1/advise"
)

type aISuggestion struct {
	Perfumes []models.Perfume `json:"perfumes"`
}

func GetAIEnrichedSuggestions(ctx context.Context, params parameters.RequestPerfume) []models.PerfumeWithScore {
	var mostSimilar []models.PerfumeWithScore
	aiSuggests, err := AISuggest(ctx, params)
	if err == nil && aiSuggests != nil {
		mostSimilar = aiSuggests
	}

	if mostSimilar != nil {
		enrichmentParams := make([]parameters.RequestPerfume, len(mostSimilar))
		for i, suggestion := range mostSimilar {
			enrichmentParams[i] = *parameters.NewGet().WithBrand(suggestion.Perfume.Brand).WithName(suggestion.Perfume.Name).WithSex(params.Sex)
		}
		enrichedSuggests, ok := FetchPerfumes(ctx, enrichmentParams)
		if ok && enrichedSuggests != nil {
			for i, suggestion := range enrichedSuggests {
				mostSimilar[i].Perfume = suggestion
			}
		}
	}
	return mostSimilar
}

func AISuggest(ctx context.Context, params parameters.RequestPerfume) ([]models.PerfumeWithScore, error) {
	r, _ := http.NewRequestWithContext(ctx, "GET", aiSuggestUrl, nil)
	updateQuery(r, params)

	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get AI suggestions: %v", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil || len(body) == 0 {
		return nil, err
	}
	var suggestions aISuggestion
	err = json.Unmarshal(body, &suggestions)
	if err != nil {
		return nil, err
	}
	var suggestionsWithScore []models.PerfumeWithScore
	for _, suggestion := range suggestions.Perfumes {
		suggestionsWithScore = append(suggestionsWithScore, models.PerfumeWithScore{
			Perfume: suggestion,
		})
	}
	return suggestionsWithScore, nil
}
