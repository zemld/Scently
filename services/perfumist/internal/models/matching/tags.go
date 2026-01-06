package matching

import (
	"github.com/zemld/Scently/models"
)

type Tags struct {
	Weights
}

func NewTags(weights Weights) *Tags {
	return &Tags{Weights: weights}
}

func (m Tags) GetSimilarityScore(first models.Properties, second models.Properties) float64 {
	firstUpperNotesTags := uniteTags(first.EnrichedUpperNotes)
	firstCoreNotesTags := uniteTags(first.EnrichedCoreNotes)
	firstBaseNotesTags := uniteTags(first.EnrichedBaseNotes)

	secondUpperNotesTags := uniteTags(second.EnrichedUpperNotes)
	secondCoreNotesTags := uniteTags(second.EnrichedCoreNotes)
	secondBaseNotesTags := uniteTags(second.EnrichedBaseNotes)

	upperNotesScore := cosineSimilarity(firstUpperNotesTags, secondUpperNotesTags)
	coreNotesScore := cosineSimilarity(firstCoreNotesTags, secondCoreNotesTags)
	baseNotesScore := cosineSimilarity(firstBaseNotesTags, secondBaseNotesTags)

	return (upperNotesScore*m.UpperNotesWeight +
		coreNotesScore*m.CoreNotesWeight +
		baseNotesScore*m.BaseNotesWeight)
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

func CalculatePerfumeTags(p *models.Perfume, minimalTagCount int) {
	tags := uniteTags(
		append(
			append(p.Properties.EnrichedUpperNotes, p.Properties.EnrichedCoreNotes...),
			p.Properties.EnrichedBaseNotes...,
		),
	)
	chosenTags := make([]string, 0, len(tags))
	for tag, count := range tags {
		if count > minimalTagCount {
			chosenTags = append(chosenTags, tag)
		}
	}
	p.Properties.Tags = chosenTags
}
