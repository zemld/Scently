package advising

import (
	"context"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type Advisor interface {
	Advise(ctx context.Context, params parameters.RequestPerfume) ([]perfume.Ranked, error)
}
