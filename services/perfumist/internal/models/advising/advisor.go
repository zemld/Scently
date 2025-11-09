package advising

import (
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type Advisor interface {
	Advise(params parameters.RequestPerfume) ([]perfume.Ranked, error)
}
