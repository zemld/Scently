package advising

import (
	"errors"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type AiAdvisor struct {
	adviseFetcher fetching.Fetcher
	enrichFetcher fetching.Fetcher
}

func NewAiAdvisor(adviseFetcher fetching.Fetcher, enrichFetcher fetching.Fetcher) *AiAdvisor {
	return &AiAdvisor{adviseFetcher: adviseFetcher, enrichFetcher: enrichFetcher}
}

func (a AiAdvisor) Advise(params parameters.RequestPerfume) ([]perfume.Ranked, error) {
	adviseResults, ok := a.adviseFetcher.Fetch([]parameters.RequestPerfume{params})
	if !ok || adviseResults == nil || len(adviseResults) == 0 {
		return nil, errors.New("failed to get AI suggestions")
	}

	enrichmentParams := make([]parameters.RequestPerfume, len(adviseResults))
	for i, suggestion := range adviseResults {
		enrichmentParams[i] = *parameters.NewGet().WithBrand(suggestion.Brand).WithName(suggestion.Name).WithSex(params.Sex)
	}
	enrichmentResults, ok := a.enrichFetcher.Fetch(enrichmentParams)
	if !ok || enrichmentResults == nil || len(enrichmentResults) == 0 {
		return nil, errors.New("failed to get enrichment results")
	}

	rankedResults := make([]perfume.Ranked, len(enrichmentResults))
	for i, suggestion := range enrichmentResults {
		rankedResults[i] = perfume.Ranked{
			Perfume: suggestion,
			Rank:    i + 1,
		}
	}
	return rankedResults, nil
}
