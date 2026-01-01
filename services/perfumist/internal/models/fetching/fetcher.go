package fetching

import (
	"context"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/Scently/models"
)

type Fetcher interface {
	Fetch(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool)
}
