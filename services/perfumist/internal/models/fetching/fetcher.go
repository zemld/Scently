package fetching

import (
	"context"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type Fetcher interface {
	Fetch(ctx context.Context, params []parameters.RequestPerfume) ([]perfume.Perfume, bool)
}
