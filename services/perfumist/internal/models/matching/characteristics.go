package matching

import (
	"github.com/zemld/Scently/models"
)

type CharacteristicsMatcher struct {
	Weights
}

func NewCharacteristicsMatcher(weights Weights) *CharacteristicsMatcher {
	return &CharacteristicsMatcher{Weights: weights}
}

func (m CharacteristicsMatcher) GetSimilarityScore(first models.Properties, second models.Properties) float64 {
	upperNotesScore := cosineSimilarity(first.UpperCharacteristics, second.UpperCharacteristics)
	coreNotesScore := cosineSimilarity(first.CoreCharacteristics, second.CoreCharacteristics)
	baseNotesScore := cosineSimilarity(first.BaseCharacteristics, second.BaseCharacteristics)

	return (upperNotesScore*m.UpperNotesWeight +
		coreNotesScore*m.CoreNotesWeight +
		baseNotesScore*m.BaseNotesWeight)
}

func uniteCharacteristics(notes []models.EnrichedNote) map[string]float64 {
	united := make(map[string]float64)

	for _, note := range notes {
		for _, characteristic := range note.Characteristics {
			united[characteristic.Name] += characteristic.Value
		}
	}

	for characteristic, value := range united {
		united[characteristic] = value / float64(len(notes))
	}

	return united
}

func PreparePerfumeCharacteristics(p *models.Perfume) {
	p.Properties.UpperCharacteristics = uniteCharacteristics(p.Properties.EnrichedUpperNotes)
	p.Properties.CoreCharacteristics = uniteCharacteristics(p.Properties.EnrichedCoreNotes)
	p.Properties.BaseCharacteristics = uniteCharacteristics(p.Properties.EnrichedBaseNotes)
}
