package fetching

import (
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type Fetcher interface {
	Fetch(params []parameters.RequestPerfume) ([]perfume.Perfume, bool)
}
