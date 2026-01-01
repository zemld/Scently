package advising

import (
	"context"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/matching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/Scently/models"
)

type Base struct {
	fetcher     fetching.Fetcher
	matcher     matching.Matcher
	adviseCount int
}

func NewBase(fetcher fetching.Fetcher, matcher matching.Matcher, adviseCount int) *Base {
	return &Base{fetcher: fetcher, matcher: matcher, adviseCount: adviseCount}
}

func (a *Base) Advise(ctx context.Context, params parameters.RequestPerfume) ([]models.Ranked, error) {
	favouritePerfumes, ok := a.fetcher.Fetch(ctx, []parameters.RequestPerfume{params})
	if !ok {
		return nil, errors.NewServiceError("failed to interact with perfume service", nil)
	}
	if len(favouritePerfumes) == 0 {
		return nil, errors.NewNotFoundError("perfume not found")
	}
	allPerfumes, ok := a.fetcher.Fetch(ctx, []parameters.RequestPerfume{*parameters.NewGet().WithSex(params.Sex)})
	if !ok {
		return nil, errors.NewServiceError("failed to interact with perfume service", nil)
	}
	if len(allPerfumes) == 0 {
		return nil, errors.NewServiceError("no perfumes available in database", nil)
	}
	similarities := a.matcher.Find(favouritePerfumes[0], allPerfumes, a.adviseCount)
	return similarities, nil
}
