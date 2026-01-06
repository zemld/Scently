package advising

import (
	"context"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/errors"
	"github.com/zemld/Scently/perfumist/internal/models/fetching"
	"github.com/zemld/Scently/perfumist/internal/models/matching"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
	"github.com/zemld/config-manager/pkg/cm"
)

type Base struct {
	fetcher fetching.Fetcher
	matcher matching.Matcher
	cm      cm.ConfigManager
}

func NewBase(fetcher fetching.Fetcher, matcher matching.Matcher, cm cm.ConfigManager) *Base {
	return &Base{fetcher: fetcher, matcher: matcher, cm: cm}
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
	matches := matching.Find(
		matching.NewMatchData(
			a.matcher,
			favouritePerfumes[0],
			allPerfumes,
			a.cm.GetIntWithDefault("suggest_count", 4),
			a.cm.GetIntWithDefault("threads_count", 8),
		),
	)
	for i := range matches {
		matching.CalculatePerfumeTags(&matches[i].Perfume, a.cm.GetIntWithDefault("minimal_tag_count", 3))
	}
	return matches, nil
}
