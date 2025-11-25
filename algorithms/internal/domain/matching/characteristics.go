package matching

import (
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

type CharacteristicsMatcher struct {
	Weights
}

func NewCharacteristicsMatcher(weights Weights) *CharacteristicsMatcher {
	return &CharacteristicsMatcher{Weights: weights}
}

func (m CharacteristicsMatcher) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	firstUpperNotesCharacteristics := models.UniteCharacteristics(first.EnrichedUpperNotes)
	secondUpperNotesCharacteristics := models.UniteCharacteristics(second.EnrichedUpperNotes)
	firstCoreNotesCharacteristics := models.UniteCharacteristics(first.EnrichedCoreNotes)
	secondCoreNotesCharacteristics := models.UniteCharacteristics(second.EnrichedCoreNotes)
	firstBaseNotesCharacteristics := models.UniteCharacteristics(first.EnrichedBaseNotes)
	secondBaseNotesCharacteristics := models.UniteCharacteristics(second.EnrichedBaseNotes)

	upperNotesScore := cosineSimilarity(firstUpperNotesCharacteristics, secondUpperNotesCharacteristics)
	coreNotesScore := cosineSimilarity(firstCoreNotesCharacteristics, secondCoreNotesCharacteristics)
	baseNotesScore := cosineSimilarity(firstBaseNotesCharacteristics, secondBaseNotesCharacteristics)

	return (upperNotesScore*m.UpperNotesWeight +
		coreNotesScore*m.CoreNotesWeight +
		baseNotesScore*m.BaseNotesWeight)
}
