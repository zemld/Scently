package matching

import "github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"

type Tags struct {
	Weights
}

func NewTags(weights Weights) *Tags {
	return &Tags{Weights: weights}
}

func (m Tags) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	firstUpperNotesTags := m.normalizeTags(models.UniteTags(first.EnrichedUpperNotes))
	firstCoreNotesTags := m.normalizeTags(models.UniteTags(first.EnrichedCoreNotes))
	firstBaseNotesTags := m.normalizeTags(models.UniteTags(first.EnrichedBaseNotes))

	secondUpperNotesTags := m.normalizeTags(models.UniteTags(second.EnrichedUpperNotes))
	secondCoreNotesTags := m.normalizeTags(models.UniteTags(second.EnrichedCoreNotes))
	secondBaseNotesTags := m.normalizeTags(models.UniteTags(second.EnrichedBaseNotes))

	upperNotesScore := multiplyMaps(firstUpperNotesTags, secondUpperNotesTags)
	coreNotesScore := multiplyMaps(firstCoreNotesTags, secondCoreNotesTags)
	baseNotesScore := multiplyMaps(firstBaseNotesTags, secondBaseNotesTags)

	return (upperNotesScore*m.UpperNotesWeight +
		coreNotesScore*m.CoreNotesWeight +
		baseNotesScore*m.BaseNotesWeight)
}

func (m Tags) normalizeTags(tags map[string]int) map[string]float64 {
	tagsSum := 0

	normalized := make(map[string]float64, len(tags))
	for tag, count := range tags {
		normalized[tag] += float64(count)
		tagsSum += count
	}

	for tag, raw := range normalized {
		if tagsSum != 0 {
			normalized[tag] = raw / float64(tagsSum)
		}
	}
	return normalized
}
