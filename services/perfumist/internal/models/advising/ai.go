package advising

import (
	"context"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/matching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/Scently/models"
	"github.com/zemld/config-manager/pkg/cm"
)

type AI struct {
	adviseFetcher fetching.Fetcher
	enrichFetcher fetching.Fetcher
	cm            cm.ConfigManager
}

func NewAI(adviseFetcher fetching.Fetcher, enrichFetcher fetching.Fetcher, configManager cm.ConfigManager) *AI {
	return &AI{adviseFetcher: adviseFetcher, enrichFetcher: enrichFetcher, cm: configManager}
}

func (a *AI) Advise(ctx context.Context, params parameters.RequestPerfume) ([]models.Ranked, error) {
	adviseResults, ok := a.adviseFetcher.Fetch(ctx, []parameters.RequestPerfume{params})
	if !ok {
		return nil, errors.NewServiceError("failed to interact with AI advisor service", nil)
	}
	if len(adviseResults) == 0 {
		return nil, errors.NewNotFoundError("perfume not found")
	}

	enrichmentParams := make([]parameters.RequestPerfume, len(adviseResults))
	for i, suggestion := range adviseResults {
		enrichmentParams[i] = *parameters.NewGet().WithBrand(suggestion.Brand).WithName(suggestion.Name).WithSex(params.Sex)
	}
	enrichmentResults, ok := a.enrichFetcher.Fetch(ctx, enrichmentParams)

	rankedMap := make(map[string]models.Perfume)
	if ok && enrichmentResults != nil && len(enrichmentResults) > 0 {
		for _, e := range enrichmentResults {
			rankedMap[getKey(e)] = e
		}
	}
	rankedResults := make([]models.Ranked, 0, len(adviseResults))
	for i, advise := range adviseResults {
		if enriched, ok := rankedMap[getKey(advise)]; ok {
			matching.PreparePerfumeCharacteristics(&enriched)
			matching.CalculatePerfumeTags(&enriched, a.cm.GetIntWithDefault("minimal_tag_count", 3))
			rankedResults = append(rankedResults, models.Ranked{
				Perfume: enriched,
				Rank:    i + 1,
			})
		} else {
			rankedResults = append(rankedResults, models.Ranked{
				Perfume: advise,
				Rank:    i + 1,
			})
		}
	}
	return rankedResults, nil
}

func getKey(p models.Perfume) string {
	return p.Brand + p.Name
}
