package matching

import "github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"

type Matcher interface {
	Find(favourite perfume.Perfume, all []perfume.Perfume, matchesCount int) []perfume.Ranked
}
