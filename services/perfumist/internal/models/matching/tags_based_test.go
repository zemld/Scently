package matching

import (
	"testing"

	"github.com/zemld/Scently/models"
)

func TestNewTagsBasedAdapter(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tagsBased := NewTagsBased(10)
	requestedTags := []string{"floral", "sweet", "floral", "romantic"}

	adapter := NewTagsBasedAdapter(*weights, tagsBased, requestedTags)

	if adapter == nil {
		t.Fatal("expected non-nil TagsBasedAdapter")
	}
	if adapter.Weights.UpperNotesWeight != 0.3 {
		t.Fatalf("expected UpperNotesWeight 0.3, got %f", adapter.Weights.UpperNotesWeight)
	}
	if adapter.TagsBased != tagsBased {
		t.Fatal("expected TagsBased to be set")
	}
	// floral appears 2 times, sweet 1 time, romantic 1 time
	if adapter.RequestedTags["floral"] != 2 {
		t.Fatalf("expected floral count 2, got %f", adapter.RequestedTags["floral"])
	}
	if adapter.RequestedTags["sweet"] != 1 {
		t.Fatalf("expected sweet count 1, got %f", adapter.RequestedTags["sweet"])
	}
	if adapter.RequestedTags["romantic"] != 1 {
		t.Fatalf("expected romantic count 1, got %f", adapter.RequestedTags["romantic"])
	}
}

func TestNewTagsBasedAdapter_EmptyTags(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tagsBased := NewTagsBased(10)
	requestedTags := []string{}

	adapter := NewTagsBasedAdapter(*weights, tagsBased, requestedTags)

	if adapter == nil {
		t.Fatal("expected non-nil TagsBasedAdapter")
	}
	if len(adapter.RequestedTags) != 0 {
		t.Fatalf("expected empty RequestedTags, got %d items", len(adapter.RequestedTags))
	}
}

func TestTagsBasedAdapter_GetSimilarityScore(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tagsBased := NewTagsBased(10)
	requestedTags := []string{"floral", "sweet", "romantic"}

	adapter := NewTagsBasedAdapter(*weights, tagsBased, requestedTags)

	first := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{},
		EnrichedCoreNotes:  []models.EnrichedNote{},
		EnrichedBaseNotes:  []models.EnrichedNote{},
	}

	second := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral", "romantic"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet", "warm"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody"}},
		},
	}

	score := adapter.GetSimilarityScore(first, second)

	if score < 0 {
		t.Fatalf("expected non-negative score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

func TestTagsBasedAdapter_GetSimilarityScore_WithWeights(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.5, 0.3, 0.2)
	tagsBased := NewTagsBased(10)
	requestedTags := []string{"floral", "sweet"}

	adapter := NewTagsBasedAdapter(*weights, tagsBased, requestedTags)

	first := models.Properties{}

	second := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{},
	}

	score := adapter.GetSimilarityScore(first, second)

	if score < 0 {
		t.Fatalf("expected non-negative score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

func TestTagsBasedAdapter_GetSimilarityScore_EmptyNotes(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tagsBased := NewTagsBased(10)
	requestedTags := []string{"floral", "sweet"}

	adapter := NewTagsBasedAdapter(*weights, tagsBased, requestedTags)

	first := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{},
		EnrichedCoreNotes:  []models.EnrichedNote{},
		EnrichedBaseNotes:  []models.EnrichedNote{},
	}

	second := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{},
		EnrichedCoreNotes:  []models.EnrichedNote{},
		EnrichedBaseNotes:  []models.EnrichedNote{},
	}

	score := adapter.GetSimilarityScore(first, second)

	expected := 0.0
	if score != expected {
		t.Fatalf("expected score %f for empty notes, got %f", expected, score)
	}
}

func TestNewTagsBased(t *testing.T) {
	t.Parallel()

	limit := 5
	tagsBased := NewTagsBased(limit)

	if tagsBased == nil {
		t.Fatal("expected non-nil TagsBased")
	}
	if tagsBased.limit != limit {
		t.Fatalf("expected limit %d, got %d", limit, tagsBased.limit)
	}
}

func TestTagsBased_GetSimilarityScore(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased(10)

	requestedTags := map[string]float64{
		"floral":   2.0,
		"sweet":    1.0,
		"romantic": 1.0,
	}

	perfumeTags := map[string]float64{
		"floral":   1.0,
		"sweet":    2.0,
		"romantic": 1.0,
		"warm":     1.0,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	if score < 0 {
		t.Fatalf("expected non-negative score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

func TestTagsBased_GetSimilarityScore_Identical(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased(10)

	tags := map[string]float64{
		"floral":   2.0,
		"sweet":    1.0,
		"romantic": 1.0,
	}

	score := tagsBased.GetSimilarityScore(tags, tags)

	// For identical tags, cosine similarity should be 1.0
	expected := 1.0
	epsilon := 0.0001
	diff := score - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected score %f for identical tags, got %f (diff: %f)", expected, score, diff)
	}
}

func TestTagsBased_GetSimilarityScore_CompletelyDifferent(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased(10)

	requestedTags := map[string]float64{
		"floral": 1.0,
		"sweet":  1.0,
	}

	perfumeTags := map[string]float64{
		"woody":  1.0,
		"earthy": 1.0,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	// No overlapping tags, cosine similarity should be 0.0
	expected := 0.0
	if score != expected {
		t.Fatalf("expected score %f for completely different tags, got %f", expected, score)
	}
}

func TestTagsBased_GetSimilarityScore_EmptyTags(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased(10)

	requestedTags := map[string]float64{}
	perfumeTags := map[string]float64{
		"floral": 1.0,
		"sweet":  1.0,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	expected := 0.0
	if score != expected {
		t.Fatalf("expected score %f for empty requested tags, got %f", expected, score)
	}
}

func TestTagsBased_GetSimilarityScore_PartialOverlap(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased(10)

	requestedTags := map[string]float64{
		"floral":   2.0,
		"sweet":    1.0,
		"romantic": 1.0,
	}

	perfumeTags := map[string]float64{
		"floral": 1.0,
		"warm":   1.0,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	// Should have partial overlap (only "floral")
	if score <= 0 {
		t.Fatalf("expected positive score for partial overlap, got %f", score)
	}
	if score >= 1.0 {
		t.Fatalf("expected score < 1.0 for partial overlap, got %f", score)
	}
}

func TestTagsBased_GetSimilarityScore_WeightedValues(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased(10)

	requestedTags := map[string]float64{
		"floral": 3.0,
		"sweet":  2.0,
	}

	perfumeTags := map[string]float64{
		"floral": 1.0,
		"sweet":  1.0,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	// Should have overlap, score should be between 0 and 1
	if score <= 0 {
		t.Fatalf("expected positive score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

func TestTagsBasedAdapter_GetSimilarityScore_ComplexCase(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tagsBased := NewTagsBased(10)
	requestedTags := []string{"floral", "sweet", "warm", "floral"}

	adapter := NewTagsBasedAdapter(*weights, tagsBased, requestedTags)

	first := models.Properties{}

	second := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral", "romantic"}},
			{Name: "Jasmine", Tags: []string{"floral", "sweet"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet", "warm"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Amber", Tags: []string{"warm", "oriental"}},
		},
	}

	score := adapter.GetSimilarityScore(first, second)

	// Should have good overlap: floral (2 in requested, 2 in upper), sweet (1 in requested, 1 in core), warm (1 in requested, 1 in core+base)
	if score <= 0 {
		t.Fatalf("expected positive score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

