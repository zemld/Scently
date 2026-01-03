package matching

import "github.com/zemld/Scently/models"

type Tags struct {
	Weights
}

func NewTags(weights Weights) *Tags {
	return &Tags{Weights: weights}
}

func (m Tags) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	firstUpperNotesTags := m.normalizeTags(uniteTags(first.EnrichedUpperNotes))
	firstCoreNotesTags := m.normalizeTags(uniteTags(first.EnrichedCoreNotes))
	firstBaseNotesTags := m.normalizeTags(uniteTags(first.EnrichedBaseNotes))

	secondUpperNotesTags := m.normalizeTags(uniteTags(second.EnrichedUpperNotes))
	secondCoreNotesTags := m.normalizeTags(uniteTags(second.EnrichedCoreNotes))
	secondBaseNotesTags := m.normalizeTags(uniteTags(second.EnrichedBaseNotes))

	upperNotesScore := cosineSimilarity(firstUpperNotesTags, secondUpperNotesTags)
	coreNotesScore := cosineSimilarity(firstCoreNotesTags, secondCoreNotesTags)
	baseNotesScore := cosineSimilarity(firstBaseNotesTags, secondBaseNotesTags)

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

func uniteTags(notes []models.EnrichedNote) map[string]int {
	united := make(map[string]int)

	for _, note := range notes {
		for _, tag := range note.Tags {
			united[tag]++
		}
	}

	return united
}

func calculatePerfumeTags(p *models.Perfume) {
	tags := uniteTags(
		append(
			append(p.Properties.EnrichedUpperNotes, p.Properties.EnrichedCoreNotes...),
			p.Properties.EnrichedBaseNotes...,
		),
	)
	chosenTags := make([]string, 0, len(tags))
	for tag, count := range tags {
		if count > 1 {
			chosenTags = append(chosenTags, tag)
		}
	}
	p.Properties.Tags = chosenTags
}
