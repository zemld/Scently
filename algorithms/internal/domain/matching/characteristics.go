package matching

import "github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"

type CharacteristicsMatcher struct {
	Weights
}

func NewCharacteristicsMatcher(weights Weights) *CharacteristicsMatcher {
	return &CharacteristicsMatcher{Weights: weights}
}

func (m CharacteristicsMatcher) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	return multiplyMaps(
		models.UniteCharacteristics(first.EnrichedUpperNotes),
		models.UniteCharacteristics(second.EnrichedUpperNotes),
	)*m.UpperNotesWeight +
		multiplyMaps(
			models.UniteCharacteristics(first.EnrichedCoreNotes),
			models.UniteCharacteristics(second.EnrichedCoreNotes),
		)*m.CoreNotesWeight +
		multiplyMaps(
			models.UniteCharacteristics(first.EnrichedBaseNotes),
			models.UniteCharacteristics(second.EnrichedBaseNotes),
		)*m.BaseNotesWeight
}
