package matching

import "github.com/zemld/Scently/models"

type CombinedMatcher struct {
	Weights
}

func NewCombinedMatcher(weights Weights) *CombinedMatcher {
	return &CombinedMatcher{Weights: weights}
}

func (m CombinedMatcher) GetSimilarityScore(first models.Properties, second models.Properties) float64 {
	cm := NewCharacteristicsMatcher(m.Weights)
	tm := NewTags(m.Weights)
	om := NewOverlay(m.Weights)

	return (m.CharacteristicsWeight*cm.GetSimilarityScore(first, second) +
		m.TagsWeight*tm.GetSimilarityScore(first, second) +
		m.OverlayWeight*om.GetSimilarityScore(first, second))
}
