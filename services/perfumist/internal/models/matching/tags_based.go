package matching

import "github.com/zemld/Scently/models"

type TagsBasedAdapter struct {
	Weights       Weights
	TagsBased     *TagsBased
	RequestedTags map[string]float64
}

type TagsBased struct {
	limit int
}

func NewTagsBasedAdapter(weights Weights, tagsBased *TagsBased, requestedTags []string) *TagsBasedAdapter {
	requestedTagsMap := make(map[string]float64)
	for _, tag := range requestedTags {
		requestedTagsMap[tag]++
	}
	return &TagsBasedAdapter{Weights: weights, TagsBased: tagsBased, RequestedTags: requestedTagsMap}
}

func (a *TagsBasedAdapter) GetSimilarityScore(first models.Properties, second models.Properties) float64 {
	perfumeUpperNotesTags := uniteTags(second.EnrichedUpperNotes)
	perfumeCoreNotesTags := uniteTags(second.EnrichedCoreNotes)
	perfumeBaseNotesTags := uniteTags(second.EnrichedBaseNotes)

	perfumeTags := make(map[string]float64)
	for tag, count := range perfumeUpperNotesTags {
		perfumeTags[tag] += float64(count) * a.Weights.UpperNotesWeight
	}
	for tag, count := range perfumeCoreNotesTags {
		perfumeTags[tag] += float64(count) * a.Weights.CoreNotesWeight
	}
	for tag, count := range perfumeBaseNotesTags {
		perfumeTags[tag] += float64(count) * a.Weights.BaseNotesWeight
	}

	return a.TagsBased.GetSimilarityScore(a.RequestedTags, perfumeTags)
}

func NewTagsBased(limit int) *TagsBased {
	return &TagsBased{limit: limit}
}

func (m *TagsBased) GetSimilarityScore(requestedTags map[string]float64, perfumeTags map[string]float64) float64 {
	return cosineSimilarity(requestedTags, perfumeTags)
}
