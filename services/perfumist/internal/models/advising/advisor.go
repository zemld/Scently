package advising

import (
	"context"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/Scently/models"
)

type Advisor interface {
	Advise(ctx context.Context, params parameters.RequestPerfume) ([]models.Ranked, error)
}
