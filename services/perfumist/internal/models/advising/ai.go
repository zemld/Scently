package advising

import (
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/errors"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type AI struct {
	adviseFetcher fetching.Fetcher
	enrichFetcher fetching.Fetcher
}

func NewAI(adviseFetcher fetching.Fetcher, enrichFetcher fetching.Fetcher) *AI {
	return &AI{adviseFetcher: adviseFetcher, enrichFetcher: enrichFetcher}
}

func (a AI) Advise(params parameters.RequestPerfume) ([]perfume.Ranked, error) {
	adviseResults, ok := a.adviseFetcher.Fetch([]parameters.RequestPerfume{params})
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
	enrichmentResults, ok := a.enrichFetcher.Fetch(enrichmentParams)

	rankedMap := make(map[string]perfume.Perfume)
	if ok && enrichmentResults != nil && len(enrichmentResults) > 0 {
		for _, e := range enrichmentResults {
			rankedMap[getKey(e)] = e
		}
	}
	rankedResults := make([]perfume.Ranked, 0, len(adviseResults))
	for i, advise := range adviseResults {
		if enriched, ok := rankedMap[getKey(advise)]; ok {
			rankedResults = append(rankedResults, perfume.Ranked{
				Perfume: enriched,
				Rank:    i + 1,
			})
		} else {
			rankedResults = append(rankedResults, perfume.Ranked{
				Perfume: advise,
				Rank:    i + 1,
			})
		}
	}
	return rankedResults, nil
}

func getKey(p perfume.Perfume) string {
	return p.Brand + p.Name
}
