package advising

import (
	"context"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/models/fetching"
	"github.com/zemld/Scently/perfumist/internal/models/matching"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
	"github.com/zemld/config-manager/pkg/cm"
)

type TagsBased struct {
	matcher matching.Matcher
	fetcher fetching.Fetcher
	cm      cm.ConfigManager
}

func NewTagsBased(matcher matching.Matcher, fetcher fetching.Fetcher, cm cm.ConfigManager) *TagsBased {
	return &TagsBased{matcher: matcher, fetcher: fetcher, cm: cm}
}

func (a *TagsBased) Advise(ctx context.Context, params parameters.RequestPerfume) ([]models.Ranked, error) {
	common := NewCommon(a.fetcher, a.matcher, a.cm)
	return common.Advise(ctx, params)
}
