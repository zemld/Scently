package matching

import (
	"maps"
	"math"

	"github.com/zemld/Scently/models"
)

type TagsBasedAdapter struct {
	Weights       Weights
	TagsBased     *TagsBased
	RequestedTags map[string]int
}

type TagsBased struct{}

func NewTagsBasedAdapter(weights Weights, requestedTags []string) *TagsBasedAdapter {
	requestedTagsMap := make(map[string]int)
	for _, tag := range requestedTags {
		requestedTagsMap[tag]++
	}
	return &TagsBasedAdapter{Weights: weights, TagsBased: NewTagsBased(), RequestedTags: requestedTagsMap}
}

func NewTagsBased() *TagsBased {
	return &TagsBased{}
}

func (a *TagsBasedAdapter) GetSimilarityScore(first models.Properties, second models.Properties) float64 {
	return a.TagsBased.GetSimilarityScore(a.RequestedTags, CalculatePerfumeTags(&second, a.Weights))
}

func (m *TagsBased) GetSimilarityScore(requestedTags map[string]int, perfumeTags map[string]int) float64 {
	intersection := getTagMapsIntersection(requestedTags, perfumeTags)
	union := getTagMapsUnion(requestedTags, perfumeTags)
	if len(union) == 0 {
		return 0.0
	}
	return float64(len(intersection)) / float64(len(union))
}

func CalculatePerfumeTags(p *models.Properties, weights Weights) map[string]int {
	perfumeUpperNotesTags := uniteTags(p.EnrichedUpperNotes)
	perfumeCoreNotesTags := uniteTags(p.EnrichedCoreNotes)
	perfumeBaseNotesTags := uniteTags(p.EnrichedBaseNotes)

	perfumeTags := make(map[string]int)
	for tag, count := range perfumeUpperNotesTags {
		perfumeTags[tag] += int(math.Round(float64(count) * weights.UpperNotesWeight))
	}
	for tag, count := range perfumeCoreNotesTags {
		perfumeTags[tag] += int(math.Round(float64(count) * weights.CoreNotesWeight))
	}
	for tag, count := range perfumeBaseNotesTags {
		perfumeTags[tag] += int(math.Round(float64(count) * weights.BaseNotesWeight))
	}

	roundedPerfumeTags := make(map[string]int)
	for tag, count := range perfumeTags {
		roundedCount := math.Round(float64(count))
		if roundedCount > 0 {
			roundedPerfumeTags[tag] = int(roundedCount)
		}
	}
	return roundedPerfumeTags
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

func getTagMapsIntersection(requestedTags map[string]int, perfumeTags map[string]int) map[string]int {
	if len(requestedTags) > len(perfumeTags) {
		return getTagMapsIntersection(perfumeTags, requestedTags)
	}
	mapsIntersect := make(map[string]int, len(requestedTags))
	for tag, count := range requestedTags {
		if otherCount, ok := perfumeTags[tag]; ok {
			mapsIntersect[tag] = int(math.Min(float64(count), float64(otherCount)))
		}
	}
	return mapsIntersect
}

func getTagMapsUnion(requestedTags map[string]int, perfumeTags map[string]int) map[string]int {
	mapsUnion := make(map[string]int, len(requestedTags))
	maps.Copy(mapsUnion, requestedTags)
	for tag, count := range perfumeTags {
		mapsUnion[tag] = int(math.Max(float64(mapsUnion[tag]), float64(count)))
	}
	return mapsUnion
}
