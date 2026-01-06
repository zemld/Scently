package matching

import (
	"github.com/zemld/Scently/models"
)

type Tags struct {
	Weights
	TagsBased *TagsBased
}

func NewTags(weights Weights) *Tags {
	return &Tags{Weights: weights, TagsBased: NewTagsBased()}
}

func (m Tags) GetSimilarityScore(first models.Properties, second models.Properties) float64 {
	firstTags := CalculatePerfumeTags(&first, m.Weights)
	secondTags := CalculatePerfumeTags(&second, m.Weights)
	return m.TagsBased.GetSimilarityScore(firstTags, secondTags)
}
