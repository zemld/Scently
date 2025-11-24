package matching

import (
	"math"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

type TagsOverlay struct {
	Weights
}

func NewTagsMatcher(weights Weights) *TagsOverlay {
	return &TagsOverlay{Weights: weights}
}

func (m TagsOverlay) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	upperNotesScore := m.getNotesSimilarityScore(
		models.UniteTags(first.EnrichedUpperNotes),
		models.UniteTags(second.EnrichedUpperNotes),
	)
	coreNotesScore := m.getNotesSimilarityScore(
		models.UniteTags(first.EnrichedCoreNotes),
		models.UniteTags(second.EnrichedCoreNotes),
	)
	baseNotesScore := m.getNotesSimilarityScore(
		models.UniteTags(first.EnrichedBaseNotes),
		models.UniteTags(second.EnrichedBaseNotes),
	)

	return (upperNotesScore*m.UpperNotesWeight +
		coreNotesScore*m.CoreNotesWeight +
		baseNotesScore*m.BaseNotesWeight)
}

func (m TagsOverlay) getNotesSimilarityScore(first map[string]int, second map[string]int) float64 {
	matches := 0
	maxMatches := 0

	for tag, firstCount := range first {
		if secondCount, ok := second[tag]; !ok {
			maxMatches += firstCount
		} else {
			matches += int(math.Min(float64(firstCount), float64(secondCount)))
			maxMatches += firstCount + secondCount - matches
			delete(second, tag)
		}
	}

	for _, secondCount := range second {
		maxMatches += secondCount
	}

	if maxMatches == 0 {
		return 0
	}
	return float64(matches) / float64(maxMatches)
}
