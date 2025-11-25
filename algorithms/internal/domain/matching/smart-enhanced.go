package matching

import "github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"

type SmartEnhanced struct {
	Weights
}

func NewSmartEnhanced(weights Weights) *SmartEnhanced {
	return &SmartEnhanced{Weights: weights}
}

func (m SmartEnhanced) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	cm := NewCharacteristicsMatcher(m.Weights)
	tm := NewTags(m.Weights)
	om := NewOverlay(m.Weights)

	return (m.CharacteristicsWeight*cm.GetPerfumeSimilarityScore(first, second) +
		m.TagsWeight*tm.GetPerfumeSimilarityScore(first, second) +
		m.OverlayWeight*om.GetPerfumeSimilarityScore(first, second))
}
