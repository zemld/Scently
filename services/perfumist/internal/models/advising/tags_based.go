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

type TagsBased struct {
	matcher matching.TagsBasedAdapter
	fetcher fetching.Fetcher
	cm      cm.ConfigManager
}

func NewTagsBased(matcher matching.TagsBasedAdapter, fetcher fetching.Fetcher, cm cm.ConfigManager) *TagsBased {
	return &TagsBased{matcher: matcher, fetcher: fetcher, cm: cm}
}

func (a *TagsBased) Advise(ctx context.Context, params parameters.RequestPerfume) ([]models.Ranked, error) {
	perfumes, ok := a.fetcher.Fetch(ctx, []parameters.RequestPerfume{*parameters.NewGet().WithSex(params.Sex)})
	if !ok {
		return nil, errors.NewServiceError("failed to interact with perfume service", nil)
	}
	if len(perfumes) == 0 {
		return nil, errors.NewServiceError("no perfumes available in database", nil)
	}
	matches := matching.Find(
		matching.NewMatchData(
			&a.matcher,
			models.Perfume{},
			perfumes,
			a.cm.GetIntWithDefault("suggest_count", 4),
			a.cm.GetIntWithDefault("threads_count", 8),
		))
	for i := range matches {
		matching.CalculatePerfumeTags(&matches[i].Perfume, a.cm.GetIntWithDefault("minimal_tag_count", 3))
	}
	return matches, nil
}
