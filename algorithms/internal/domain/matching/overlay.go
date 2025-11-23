package matching

import (
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/pkg/set"
)

type Overlay struct {
	Weights
}

func NewOverlay(weights Weights) *Overlay {
	return &Overlay{Weights: weights}
}

func (m Overlay) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	familiesSimilarityScore := m.getListSimilarityScore(first.Family, second.Family)
	notesSimilarityScore := m.getNotesSimilarityScore(first, second)
	typeSimilarity := m.getTypeSimilarityScore(first.Type, second.Type)
	return familiesSimilarityScore*m.familyWeight + notesSimilarityScore*m.notesWeight + typeSimilarity*m.typeWeight
}

func (m Overlay) getListSimilarityScore(first []string, second []string) float64 {
	firstSet := set.MakeSet(first)
	secondSet := set.MakeSet(second)
	intersection := set.Intersect(firstSet, secondSet)
	un := set.Union(firstSet, secondSet)

	if len(un) == 0 {
		return 0
	}
	return float64(len(intersection)) / float64(len(un))
}

func (m Overlay) getNotesSimilarityScore(first models.Properties, second models.Properties) float64 {
	upperNotesSimilarityScore := m.getListSimilarityScore(first.UpperNotes, second.UpperNotes)
	middleNotesSimilarityScore := m.getListSimilarityScore(first.CoreNotes, second.CoreNotes)
	baseNotesSimilarityScore := m.getListSimilarityScore(first.BaseNotes, second.BaseNotes)

	return upperNotesSimilarityScore*m.upperNotesWeight + middleNotesSimilarityScore*m.coreNotesWeight + baseNotesSimilarityScore*m.baseNotesWeight
}

func (m Overlay) getTypeSimilarityScore(first string, second string) float64 {
	if first == second {
		return 1
	}
	return 0
}
