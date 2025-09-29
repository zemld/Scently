package similarity

import (
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/models"
)

const (
	male   = "male"
	female = "female"
)

func GetPerfumeSimilarityScore(first models.Perfume, second models.Perfume) float64 {
	if first.Sex == male && second.Sex == female || first.Sex == female && second.Sex == male {
		return 0
	}
	familiesSimilarityScore := getListSimilarityScore(first.Family, second.Family)
	notesSimilarityScore := getNotesSimilarityScore(first, second)
	typeSimilarity := getTypeSimilarityScore(first.Type, second.Type)
	return familiesSimilarityScore*familyWeight + notesSimilarityScore*notesWeight + typeSimilarity*typeWeight
}
