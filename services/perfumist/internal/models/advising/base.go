package advising

import (
	"errors"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/fetching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/matching"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type Base struct {
	fetcher     fetching.Fetcher
	matcher     matching.Matcher
	adviseCount int
}

func NewBase(fetcher fetching.Fetcher, matcher matching.Matcher, adviseCount int) *Base {
	return &Base{fetcher: fetcher, matcher: matcher, adviseCount: adviseCount}
}

func (a *Base) Advise(params parameters.RequestPerfume) ([]perfume.Ranked, error) {
	favouritePerfumes, ok := a.fetcher.Fetch([]parameters.RequestPerfume{params})
	if !ok || favouritePerfumes == nil || len(favouritePerfumes) == 0 {
		return nil, errors.New("failed to get favourite perfumes")
	}
	allPerfumes, ok := a.fetcher.Fetch([]parameters.RequestPerfume{*parameters.NewGet().WithSex(params.Sex)})
	if !ok || allPerfumes == nil || len(allPerfumes) == 0 {
		return nil, errors.New("failed to get all perfumes")
	}
	similarities := a.matcher.Find(favouritePerfumes[0], allPerfumes, a.adviseCount)
	return similarities, nil
}
