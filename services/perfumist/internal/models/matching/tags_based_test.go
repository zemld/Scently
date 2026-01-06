package matching

import (
	"testing"

	"github.com/zemld/Scently/models"
)

func TestNewTagsBasedAdapter(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tagsBased := NewTagsBased()
	requestedTags := []string{"floral", "sweet", "floral", "romantic"}

	adapter := NewTagsBasedAdapter(*weights, requestedTags)

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
		t.Fatalf("expected floral count 2, got %d", adapter.RequestedTags["floral"])
	}
	if adapter.RequestedTags["sweet"] != 1 {
		t.Fatalf("expected sweet count 1, got %d", adapter.RequestedTags["sweet"])
	}
	if adapter.RequestedTags["romantic"] != 1 {
		t.Fatalf("expected romantic count 1, got %d", adapter.RequestedTags["romantic"])
	}
}

func TestNewTagsBasedAdapter_EmptyTags(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	requestedTags := []string{}

	adapter := NewTagsBasedAdapter(*weights, requestedTags)

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
	requestedTags := []string{"floral", "sweet", "romantic"}

	adapter := NewTagsBasedAdapter(*weights, requestedTags)

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

	// requestedTags: {"floral": 1, "sweet": 1, "romantic": 1}
	// perfumeTags: Upper: {"floral": 1*0.3=0.3≈0, "romantic": 1*0.3=0.3≈0}
	//              Core: {"sweet": 1*0.4=0.4≈0, "warm": 1*0.4=0.4≈0}
	//              Base: {"woody": 1*0.3=0.3≈0}
	// After rounding: all become 0, so perfumeTags = {}
	// intersection: {} -> 0 elements
	// union: {"floral": 1, "sweet": 1, "romantic": 1} -> 3 elements
	// score = 0/3 = 0.0
	expected := 0.0
	epsilon := 0.0001
	diff := score - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected score %f, got %f (diff: %f)", expected, score, diff)
	}
}

func TestTagsBasedAdapter_GetSimilarityScore_WithWeights(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.5, 0.3, 0.2)
	requestedTags := []string{"floral", "sweet"}

	adapter := NewTagsBasedAdapter(*weights, requestedTags)

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

	// requestedTags: {"floral": 1, "sweet": 1}
	// perfumeTags: Upper: {"floral": 1*0.5=0.5≈1}, Core: {"sweet": 1*0.3=0.3≈0}
	// After rounding: {"floral": 1}
	// intersection: {"floral": min(1,1)=1} -> 1 element
	// union: {"floral": max(1,1)=1, "sweet": 1} -> 2 elements
	// score = 1/2 = 0.5
	expected := 0.5
	epsilon := 0.0001
	diff := score - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected score %f, got %f (diff: %f)", expected, score, diff)
	}
}

func TestTagsBasedAdapter_GetSimilarityScore_EmptyNotes(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	requestedTags := []string{"floral", "sweet"}

	adapter := NewTagsBasedAdapter(*weights, requestedTags)

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

func TestTagsBased_GetSimilarityScore(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased()

	requestedTags := map[string]int{
		"floral":   2,
		"sweet":    1,
		"romantic": 1,
	}

	perfumeTags := map[string]int{
		"floral":   1,
		"sweet":    2,
		"romantic": 1,
		"warm":     1,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	// intersection: {"floral": min(2,1)=1, "sweet": min(1,2)=1, "romantic": min(1,1)=1} -> 3 elements
	// union: {"floral": max(2,1)=2, "sweet": max(1,2)=2, "romantic": max(1,1)=1, "warm": max(0,1)=1} -> 4 elements
	// score = 3/4 = 0.75
	expected := 0.75
	epsilon := 0.0001
	diff := score - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected score %f, got %f (diff: %f)", expected, score, diff)
	}
}

func TestTagsBased_GetSimilarityScore_Identical(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased()

	tags := map[string]int{
		"floral":   2,
		"sweet":    1,
		"romantic": 1,
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

	tagsBased := NewTagsBased()

	requestedTags := map[string]int{
		"floral": 1,
		"sweet":  1,
	}

	perfumeTags := map[string]int{
		"woody":  1,
		"earthy": 1,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	// No overlapping tags, score should be 0.0
	expected := 0.0
	epsilon := 0.0001
	diff := score - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected score %f for completely different tags, got %f (diff: %f)", expected, score, diff)
	}
}

func TestTagsBased_GetSimilarityScore_EmptyTags(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased()

	requestedTags := map[string]int{}
	perfumeTags := map[string]int{
		"floral": 1,
		"sweet":  1,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	expected := 0.0
	if score != expected {
		t.Fatalf("expected score %f for empty requested tags, got %f", expected, score)
	}
}

func TestTagsBased_GetSimilarityScore_PartialOverlap(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased()

	requestedTags := map[string]int{
		"floral":   2,
		"sweet":    1,
		"romantic": 1,
	}

	perfumeTags := map[string]int{
		"floral": 1,
		"warm":   1,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	// intersection: {"floral": min(2,1)=1} -> 1 element (only tags with positive values count)
	// union: {"floral": max(2,1)=2, "sweet": 1, "romantic": 1, "warm": 1} -> 4 elements
	// score = 1/4 = 0.25
	expected := 0.25
	epsilon := 0.0001
	diff := score - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected score %f, got %f (diff: %f)", expected, score, diff)
	}
}

func TestTagsBased_GetSimilarityScore_WeightedValues(t *testing.T) {
	t.Parallel()

	tagsBased := NewTagsBased()

	requestedTags := map[string]int{
		"floral": 3,
		"sweet":  2,
	}

	perfumeTags := map[string]int{
		"floral": 1,
		"sweet":  1,
	}

	score := tagsBased.GetSimilarityScore(requestedTags, perfumeTags)

	// intersection: {"floral": min(3,1)=1, "sweet": min(2,1)=1} -> 2 elements
	// union: {"floral": max(3,1)=3, "sweet": max(2,1)=2} -> 2 elements
	// score = 2/2 = 1.0
	expected := 1.0
	epsilon := 0.0001
	diff := score - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected score %f, got %f (diff: %f)", expected, score, diff)
	}
}

func TestTagsBasedAdapter_GetSimilarityScore_ComplexCase(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	requestedTags := []string{"floral", "sweet", "warm", "floral"}

	adapter := NewTagsBasedAdapter(*weights, requestedTags)

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

	// requestedTags: {"floral": 2, "sweet": 1, "warm": 1}
	// perfumeTags: Upper: {"floral": 2*0.3=0.6≈1, "romantic": 1*0.3=0.3≈0, "sweet": 1*0.3=0.3≈0}
	//              Core: {"sweet": 1*0.4=0.4≈0, "warm": 1*0.4=0.4≈0}
	//              Base: {"warm": 1*0.3=0.3≈0, "oriental": 1*0.3=0.3≈0}
	// After first rounding: {"floral": 1}
	// After final rounding: {"floral": 1}
	// intersection: {"floral": min(2,1)=1} -> 1 element
	// union: {"floral": max(2,1)=2, "sweet": 1, "warm": 1} -> 3 elements
	// score = 1/3 ≈ 0.3333
	expected := 1.0 / 3.0
	epsilon := 0.0001
	diff := score - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected score %f, got %f (diff: %f)", expected, score, diff)
	}
}
