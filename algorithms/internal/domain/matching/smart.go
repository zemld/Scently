package matching

import "github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"

type Smart struct {
	Weights
}

func NewSmart(weights Weights) *Smart {
	return &Smart{Weights: weights}
}

func (m Smart) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	cm := NewCharacteristicsMatcher(m.Weights)
	tm := NewTags(m.Weights)

	return (m.characteristicsWeight*cm.GetPerfumeSimilarityScore(first, second) +
		m.tagsWeight*tm.GetPerfumeSimilarityScore(first, second))
}
