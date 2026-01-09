package fetching

import (
	"context"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
)

type Fetcher interface {
	Fetch(ctx context.Context, parameter parameters.RequestPerfume) <-chan models.Perfume
	FetchMany(ctx context.Context, parameters []parameters.RequestPerfume) <-chan models.Perfume
}
