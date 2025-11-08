package app

import (
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/pkg/set"
)

const (
	familyWeight = 0.4
	notesWeight  = 0.55
	typeWeight   = 0.05
)

const (
	upperNotesWeight  = 0.15
	middleNotesWeight = 0.45
	baseNotesWeight   = 0.4
)

func FoundSimilarities(favourite models.Perfume, all []models.Perfume, suggestsCount int) []models.PerfumeWithScore {
	mostSimilar := make([]models.PerfumeWithScore, 0, suggestsCount)
	for _, perfume := range all {
		if favourite.Equal(perfume) {
			continue
		}
		similarityScore := GetPerfumeSimilarityScore(favourite.Properties, perfume.Properties)
		mostSimilar = updateMostSimilarIfNeeded(mostSimilar, perfume, similarityScore)
	}
	return mostSimilar
}

func updateMostSimilarIfNeeded(mostSimilar []models.PerfumeWithScore, perfume models.Perfume, similarityScore float64) []models.PerfumeWithScore {
	current := perfume
	for i := range mostSimilar {
		if similarityScore > mostSimilar[i].Score {
			tmp := mostSimilar[i]
			mostSimilar[i].Score = similarityScore
			mostSimilar[i].Perfume = current
			current = tmp.Perfume
			similarityScore = tmp.Score
		}
	}
	if len(mostSimilar) < cap(mostSimilar) {
		mostSimilar = append(mostSimilar, models.PerfumeWithScore{
			Perfume: current,
			Score:   similarityScore,
		})
	}
	return mostSimilar
}

func GetPerfumeSimilarityScore(first models.PerfumeProperties, second models.PerfumeProperties) float64 {
	familiesSimilarityScore := getListSimilarityScore(first.Family, second.Family)
	notesSimilarityScore := getNotesSimilarityScore(first, second)
	typeSimilarity := getTypeSimilarityScore(first.Type, second.Type)
	return familiesSimilarityScore*familyWeight + notesSimilarityScore*notesWeight + typeSimilarity*typeWeight
}

func getListSimilarityScore(first []string, second []string) float64 {
	firstSet := set.MakeSet(first)
	secondSet := set.MakeSet(second)
	intersection := set.Intersect(firstSet, secondSet)
	un := set.Union(firstSet, secondSet)

	if len(un) == 0 {
		return 0
	}
	return float64(len(intersection)) / float64(len(un))
}

func getNotesSimilarityScore(first models.PerfumeProperties, second models.PerfumeProperties) float64 {
	upperNotesSimilarityScore := getListSimilarityScore(first.UpperNotes, second.UpperNotes)
	middleNotesSimilarityScore := getListSimilarityScore(first.CoreNotes, second.CoreNotes)
	baseNotesSimilarityScore := getListSimilarityScore(first.BaseNotes, second.BaseNotes)

	return upperNotesSimilarityScore*upperNotesWeight + middleNotesSimilarityScore*middleNotesWeight + baseNotesSimilarityScore*baseNotesWeight
}

func getTypeSimilarityScore(first string, second string) float64 {
	if first == second {
		return 1
	}
	return 0
}
