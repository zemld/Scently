package similarity

import "github.com/zemld/PerfumeRecommendationSystem/perfumist/models"

func getNotesSimilarityScore(first models.Perfume, second models.Perfume) float64 {
	upperNotesSimilarityScore := getListSimilarityScore(first.UpperNotes, second.UpperNotes)
	middleNotesSimilarityScore := getListSimilarityScore(first.MiddleNotes, second.MiddleNotes)
	baseNotesSimilarityScore := getListSimilarityScore(first.BaseNotes, second.BaseNotes)

	return upperNotesSimilarityScore*upperNotesWeight + middleNotesSimilarityScore*middleNotesWeight + baseNotesSimilarityScore*baseNotesWeight
}
