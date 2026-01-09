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

type AI struct {
	adviseFetcher fetching.Fetcher
	enrichFetcher fetching.Fetcher
	cm            cm.ConfigManager
}

func NewAI(adviseFetcher fetching.Fetcher, enrichFetcher fetching.Fetcher, configManager cm.ConfigManager) *AI {
	return &AI{adviseFetcher: adviseFetcher, enrichFetcher: enrichFetcher, cm: configManager}
}

func (a *AI) Advise(ctx context.Context, params parameters.RequestPerfume) ([]models.Ranked, error) {
	adviseResults, err := a.tryFetchRawAdvise(ctx, params)
	if err != nil {
		return nil, err
	}
	rankedMap := a.tryFetchEnrichments(ctx, adviseResults, params)

	return a.prepareSuggestionsWithEnrichments(adviseResults, rankedMap), nil
}

func (a *AI) tryFetchRawAdvise(ctx context.Context, params parameters.RequestPerfume) ([]models.Perfume, error) {
	adviseChan := a.adviseFetcher.Fetch(ctx, params)
	adviseResults := make([]models.Perfume, 0, a.cm.GetIntWithDefault("suggest_count", 4))
	for perfume, ok := <-adviseChan; ok; perfume, ok = <-adviseChan {
		adviseResults = append(adviseResults, perfume)
	}
	if len(adviseResults) == 0 {
		return nil, errors.NewNotFoundError("perfume not found")
	}
	return adviseResults, nil
}

func (a *AI) tryFetchEnrichments(ctx context.Context, adviseResults []models.Perfume, params parameters.RequestPerfume) map[string]models.Perfume {
	enrichmentParams := make([]parameters.RequestPerfume, len(adviseResults))
	for i, suggestion := range adviseResults {
		enrichmentParams[i] = *parameters.NewGet().WithBrand(suggestion.Brand).WithName(suggestion.Name).WithSex(params.Sex)
	}
	enrichmentChan := a.enrichFetcher.FetchMany(ctx, enrichmentParams)

	rankedMap := make(map[string]models.Perfume)
	for enrichment, ok := <-enrichmentChan; ok; enrichment, ok = <-enrichmentChan {
		rankedMap[getKey(enrichment)] = enrichment
	}
	return rankedMap
}

func (a *AI) prepareSuggestionsWithEnrichments(adviseResults []models.Perfume, rankedMap map[string]models.Perfume) []models.Ranked {
	enrichedResults := make([]models.Ranked, 0, len(adviseResults))
	for i, advise := range adviseResults {
		if enriched, ok := rankedMap[getKey(advise)]; ok {
			matching.PreparePerfumeCharacteristics(&enriched)
			enriched.Properties.Tags = matching.CalculatePerfumeTags(
				&enriched.Properties,
				*matching.NewBaseWeights(
					a.cm.GetFloatWithDefault("upper_notes_weight", 0.2),
					a.cm.GetFloatWithDefault("core_notes_weight", 0.35),
					a.cm.GetFloatWithDefault("base_notes_weight", 0.45),
				),
			)
			enrichedResults = append(enrichedResults, models.Ranked{
				Perfume: enriched,
				Rank:    i + 1,
			})
		} else {
			enrichedResults = append(enrichedResults, models.Ranked{
				Perfume: advise,
				Rank:    i + 1,
			})
		}
	}
	return enrichedResults
}

func getKey(p models.Perfume) string {
	return p.Brand + p.Name
}
