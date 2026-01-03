package matching

import "github.com/zemld/Scently/models"

type CombinedMatcher struct {
	Weights
}

func NewCombinedMatcher(weights Weights) *CombinedMatcher {
	return &CombinedMatcher{Weights: weights}
}

func (m CombinedMatcher) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	cm := NewCharacteristicsMatcher(m.Weights)
	tm := NewTags(m.Weights)
	om := NewOverlay(m.Weights)

	return (m.CharacteristicsWeight*cm.GetPerfumeSimilarityScore(first, second) +
		m.TagsWeight*tm.GetPerfumeSimilarityScore(first, second) +
		m.OverlayWeight*om.GetPerfumeSimilarityScore(first, second))
}
