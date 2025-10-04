package similarity

import "github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"

func getNotesSimilarityScore(first models.PerfumeProperties, second models.PerfumeProperties) float64 {
	upperNotesSimilarityScore := getListSimilarityScore(first.UpperNotes, second.UpperNotes)
	middleNotesSimilarityScore := getListSimilarityScore(first.MiddleNotes, second.MiddleNotes)
	baseNotesSimilarityScore := getListSimilarityScore(first.BaseNotes, second.BaseNotes)

	return upperNotesSimilarityScore*upperNotesWeight + middleNotesSimilarityScore*middleNotesWeight + baseNotesSimilarityScore*baseNotesWeight
}
