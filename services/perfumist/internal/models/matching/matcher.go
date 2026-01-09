package matching

import (
	"math"

	"github.com/zemld/Scently/models"
)

type Matcher interface {
	GetSimilarityScore(first models.Properties, second models.Properties) float64
}

func cosineSimilarity[Number ~int | ~float64](first map[string]Number, second map[string]Number) float64 {
	dotProduct := multiplyMaps(first, second)
	firstNorm := math.Sqrt(multiplyMaps(first, first))
	secondNorm := math.Sqrt(multiplyMaps(second, second))

	if firstNorm == 0 || secondNorm == 0 {
		return 0.0
	}

	return dotProduct / (firstNorm * secondNorm)
}

func multiplyMaps[Number ~int | ~float64](first map[string]Number, second map[string]Number) float64 {
	if len(second) < len(first) {
		return multiplyMaps(second, first)
	}

	score := 0.0
	for tag, value := range first {
		if secondValue, ok := second[tag]; ok {
			score += float64(value) * float64(secondValue)
		}
	}
	return score
}
