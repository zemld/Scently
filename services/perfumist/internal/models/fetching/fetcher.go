package fetching

import (
	"context"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
)

type Fetcher interface {
	Fetch(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool)
}
