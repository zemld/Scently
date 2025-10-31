package app

import (
	"context"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

const suggestsCount = 4

func GetComparisionSuggestions(ctx context.Context, params parameters.RequestPerfume) []models.GluedPerfumeWithScore {
	favouritePerfumes, ok := FetchPerfumes(ctx, []parameters.RequestPerfume{params})
	if !ok || favouritePerfumes == nil || len(favouritePerfumes) == 0 {
		return nil
	}
	allPerfumes, ok := FetchPerfumes(ctx, []parameters.RequestPerfume{*parameters.NewGet().WithSex(params.Sex)})
	if !ok {
		return nil
	}

	return FoundSimilarities(favouritePerfumes[0], allPerfumes, suggestsCount)
}
