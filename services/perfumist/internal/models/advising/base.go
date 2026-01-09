package advising

import (
	"context"
	"log"

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
	favouritePerfume, err := a.fetchFavouritePerfume(ctx, params)
	if err != nil {
		return nil, err
	}
	log.Printf("favouritePerfume: %+v\n", favouritePerfume)
	common := NewCommon(a.fetcher, a.matcher, a.cm).WithFavouritePerfume(favouritePerfume)
	return common.Advise(ctx, params)
}

func (a *Base) fetchFavouritePerfume(ctx context.Context, params parameters.RequestPerfume) (models.Perfume, error) {
	select {
	case <-ctx.Done():
		return models.Perfume{}, errors.NewServiceError("context cancelled", nil)
	case favouritePerfume, ok := <-a.fetcher.Fetch(ctx, params):
		if !ok {
			return models.Perfume{}, errors.NewServiceError("failed to interact with perfume service", nil)
		}
		return favouritePerfume, nil
	}
}
