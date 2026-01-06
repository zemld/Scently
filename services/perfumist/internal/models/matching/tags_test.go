package matching

import (
	"testing"

	"github.com/zemld/Scently/models"
)

func TestNewTags(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tags := NewTags(*weights)

	if tags == nil {
		t.Fatal("expected non-nil Tags")
	}
	if tags.UpperNotesWeight != 0.3 {
		t.Fatalf("expected UpperNotesWeight 0.3, got %f", tags.UpperNotesWeight)
	}
	if tags.CoreNotesWeight != 0.4 {
		t.Fatalf("expected CoreNotesWeight 0.4, got %f", tags.CoreNotesWeight)
	}
	if tags.BaseNotesWeight != 0.3 {
		t.Fatalf("expected BaseNotesWeight 0.3, got %f", tags.BaseNotesWeight)
	}
}

func TestTags_GetSimilarityScore(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tags := NewTags(*weights)

	first := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral", "romantic"}},
			{Name: "Jasmine", Tags: []string{"floral", "sweet"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet", "warm"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody", "animalic"}},
		},
	}

	second := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral", "romantic"}},
			{Name: "Lily", Tags: []string{"floral", "fresh"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet", "warm"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody", "animalic"}},
		},
	}

	score := tags.GetSimilarityScore(first, second)

	if score <= 0 {
		t.Fatalf("expected positive score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

func TestTags_GetSimilarityScore_Identical(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tags := NewTags(*weights)

	props := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral", "romantic"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody"}},
		},
	}

	score := tags.GetSimilarityScore(props, props)

	// For identical tags, cosine similarity should be 1.0
	// Score should be: 1.0*0.3 + 1.0*0.4 + 1.0*0.3 = 1.0
	expected := 1.0
	if score != expected {
		t.Fatalf("expected score %f for identical properties, got %f", expected, score)
	}
}

func TestTags_GetSimilarityScore_CompletelyDifferent(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tags := NewTags(*weights)

	first := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody"}},
		},
	}

	second := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Citrus", Tags: []string{"fresh", "spicy"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Amber", Tags: []string{"oriental", "warm"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Patchouli", Tags: []string{"earthy", "green"}},
		},
	}

	score := tags.GetSimilarityScore(first, second)

	// No overlapping tags, cosine similarity should be 0.0
	expected := 0.0
	if score != expected {
		t.Fatalf("expected score %f for completely different properties, got %f", expected, score)
	}
}

func TestTags_GetSimilarityScore_EmptyNotes(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tags := NewTags(*weights)

	empty := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{},
		EnrichedCoreNotes:  []models.EnrichedNote{},
		EnrichedBaseNotes:  []models.EnrichedNote{},
	}

	filled := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody"}},
		},
	}

	score := tags.GetSimilarityScore(empty, filled)

	// Empty notes should result in 0.0 cosine similarity
	expected := 0.0
	if score != expected {
		t.Fatalf("expected score %f for empty notes, got %f", expected, score)
	}
}

func TestTags_GetSimilarityScore_PartialOverlap(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	tags := NewTags(*weights)

	first := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral", "romantic"}},
			{Name: "Jasmine", Tags: []string{"floral", "sweet"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet", "warm"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody"}},
		},
	}

	second := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral", "romantic"}},
			{Name: "Lily", Tags: []string{"floral", "fresh"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody", "animalic"}},
		},
	}

	score := tags.GetSimilarityScore(first, second)

	// Should have partial overlap
	if score <= 0 {
		t.Fatalf("expected positive score for partial overlap, got %f", score)
	}
	if score >= 1.0 {
		t.Fatalf("expected score < 1.0 for partial overlap, got %f", score)
	}
}

func TestTags_GetSimilarityScore_Weights(t *testing.T) {
	t.Parallel()

	// Test with different weights
	weights := NewBaseWeights(0.5, 0.3, 0.2)
	tags := NewTags(*weights)

	props := models.Properties{
		EnrichedUpperNotes: []models.EnrichedNote{
			{Name: "Rose", Tags: []string{"floral"}},
		},
		EnrichedCoreNotes: []models.EnrichedNote{
			{Name: "Vanilla", Tags: []string{"sweet"}},
		},
		EnrichedBaseNotes: []models.EnrichedNote{
			{Name: "Musk", Tags: []string{"woody"}},
		},
	}

	score := tags.GetSimilarityScore(props, props)

	// For identical properties, should be 1.0 regardless of weights
	expected := 1.0
	if score != expected {
		t.Fatalf("expected score %f for identical properties, got %f", expected, score)
	}
}

func TestUniteTags(t *testing.T) {
	t.Parallel()

	notes := []models.EnrichedNote{
		{Name: "Rose", Tags: []string{"floral", "romantic"}},
		{Name: "Jasmine", Tags: []string{"floral", "sweet"}},
	}

	united := uniteTags(notes)

	// floral: 2 occurrences
	// romantic: 1 occurrence
	// sweet: 1 occurrence
	if united["floral"] != 2 {
		t.Fatalf("expected floral count 2, got %d", united["floral"])
	}
	if united["romantic"] != 1 {
		t.Fatalf("expected romantic count 1, got %d", united["romantic"])
	}
	if united["sweet"] != 1 {
		t.Fatalf("expected sweet count 1, got %d", united["sweet"])
	}
}

func TestUniteTags_EmptyNotes(t *testing.T) {
	t.Parallel()

	notes := []models.EnrichedNote{}
	united := uniteTags(notes)

	if len(united) != 0 {
		t.Fatalf("expected empty map for empty notes, got %d items", len(united))
	}
}

func TestUniteTags_NoTags(t *testing.T) {
	t.Parallel()

	notes := []models.EnrichedNote{
		{Name: "Rose", Tags: []string{}},
		{Name: "Jasmine", Tags: []string{}},
	}

	united := uniteTags(notes)

	if len(united) != 0 {
		t.Fatalf("expected empty map for notes without tags, got %d items", len(united))
	}
}
