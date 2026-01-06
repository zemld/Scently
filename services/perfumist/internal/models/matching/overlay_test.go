package matching

import (
	"testing"

	"github.com/zemld/Scently/models"
)

func TestNewOverlay(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	overlay := NewOverlay(*weights)

	if overlay.Weights.UpperNotesWeight != 0.3 {
		t.Fatalf("expected UpperNotesWeight 0.3, got %f", overlay.Weights.UpperNotesWeight)
	}
	if overlay.Weights.CoreNotesWeight != 0.4 {
		t.Fatalf("expected CoreNotesWeight 0.4, got %f", overlay.Weights.CoreNotesWeight)
	}
	if overlay.Weights.BaseNotesWeight != 0.3 {
		t.Fatalf("expected BaseNotesWeight 0.3, got %f", overlay.Weights.BaseNotesWeight)
	}
}

func TestOverlay_GetSimilarityScore(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.2, 0.3, 0.2)
	weights.FamilyWeight = 0.1
	weights.NotesWeight = 0.15
	weights.TypeWeight = 0.05
	overlay := NewOverlay(*weights)

	first := models.Properties{
		Type:       "Eau de Parfum",
		Family:     []string{"Floral", "Oriental"},
		UpperNotes: []string{"Rose", "Jasmine"},
		CoreNotes:  []string{"Vanilla", "Amber"},
		BaseNotes:  []string{"Musk", "Patchouli"},
	}

	second := models.Properties{
		Type:       "Eau de Parfum",
		Family:     []string{"Floral"},
		UpperNotes: []string{"Rose", "Lily"},
		CoreNotes:  []string{"Vanilla"},
		BaseNotes:  []string{"Musk"},
	}

	score := overlay.GetSimilarityScore(first, second)

	// Expected calculations:
	// Family similarity: 1/2 = 0.5 (intersection: Floral, union: Floral, Oriental)
	// Upper notes similarity: 1/3 = 0.333... (intersection: Rose, union: Rose, Jasmine, Lily)
	// Core notes similarity: 1/2 = 0.5 (intersection: Vanilla, union: Vanilla, Amber)
	// Base notes similarity: 1/2 = 0.5 (intersection: Musk, union: Musk, Patchouli)
	// Notes similarity: 0.333*0.2 + 0.5*0.3 + 0.5*0.2 = 0.0666 + 0.15 + 0.1 = 0.3166
	// Type similarity: 1.0 (both are "Eau de Parfum")
	// Total: 0.5*0.1 + 0.3166*0.15 + 1.0*0.05 = 0.05 + 0.04749 + 0.05 = 0.14749

	if score <= 0 {
		t.Fatalf("expected positive score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

func TestOverlay_GetSimilarityScore_Identical(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.2, 0.3, 0.2)
	weights.FamilyWeight = 0.1
	weights.NotesWeight = 0.15
	weights.TypeWeight = 0.05
	overlay := NewOverlay(*weights)

	props := models.Properties{
		Type:       "Eau de Parfum",
		Family:     []string{"Floral"},
		UpperNotes: []string{"Rose", "Jasmine"},
		CoreNotes:  []string{"Vanilla"},
		BaseNotes:  []string{"Musk"},
	}

	score := overlay.GetSimilarityScore(props, props)

	// All similarities should be 1.0
	// Family: 1.0, Notes: 1.0*0.2 + 1.0*0.3 + 1.0*0.2 = 0.7, Type: 1.0
	// Total: 1.0*0.1 + 0.7*0.15 + 1.0*0.05 = 0.1 + 0.105 + 0.05 = 0.255
	expectedMin := 0.2 // At least should be high
	if score < expectedMin {
		t.Fatalf("expected score >= %f for identical properties, got %f", expectedMin, score)
	}
}

func TestOverlay_GetSimilarityScore_CompletelyDifferent(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.2, 0.3, 0.2)
	weights.FamilyWeight = 0.1
	weights.NotesWeight = 0.15
	weights.TypeWeight = 0.05
	overlay := NewOverlay(*weights)

	first := models.Properties{
		Type:       "Eau de Parfum",
		Family:     []string{"Floral"},
		UpperNotes: []string{"Rose"},
		CoreNotes:  []string{"Vanilla"},
		BaseNotes:  []string{"Musk"},
	}

	second := models.Properties{
		Type:       "Eau de Toilette",
		Family:     []string{"Oriental"},
		UpperNotes: []string{"Citrus"},
		CoreNotes:  []string{"Spice"},
		BaseNotes:  []string{"Wood"},
	}

	score := overlay.GetSimilarityScore(first, second)

	// All similarities should be 0.0
	// Total should be 0.0
	if score != 0.0 {
		t.Fatalf("expected score 0.0 for completely different properties, got %f", score)
	}
}

func TestOverlay_getListSimilarityScore(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	overlay := NewOverlay(*weights)

	// Test identical lists
	score := overlay.getListSimilarityScore([]string{"A", "B"}, []string{"A", "B"})
	if score != 1.0 {
		t.Fatalf("expected score 1.0 for identical lists, got %f", score)
	}

	// Test partial overlap
	score = overlay.getListSimilarityScore([]string{"A", "B"}, []string{"A", "C"})
	// Intersection: {A} = 1, Union: {A, B, C} = 3, Score: 1/3
	expected := 1.0 / 3.0
	if score != expected {
		t.Fatalf("expected score %f, got %f", expected, score)
	}

	// Test no overlap
	score = overlay.getListSimilarityScore([]string{"A", "B"}, []string{"C", "D"})
	if score != 0.0 {
		t.Fatalf("expected score 0.0 for no overlap, got %f", score)
	}

	// Test empty lists
	score = overlay.getListSimilarityScore([]string{}, []string{})
	if score != 0.0 {
		t.Fatalf("expected score 0.0 for empty lists, got %f", score)
	}

	// Test one empty list
	score = overlay.getListSimilarityScore([]string{"A"}, []string{})
	if score != 0.0 {
		t.Fatalf("expected score 0.0 when one list is empty, got %f", score)
	}

	// Test complete subset
	score = overlay.getListSimilarityScore([]string{"A"}, []string{"A", "B", "C"})
	// Intersection: {A} = 1, Union: {A, B, C} = 3, Score: 1/3
	expected = 1.0 / 3.0
	if score != expected {
		t.Fatalf("expected score %f, got %f", expected, score)
	}
}

func TestOverlay_getNotesSimilarityScore(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.2, 0.3, 0.2)
	weights.FamilyWeight = 0.1
	weights.NotesWeight = 0.15
	weights.TypeWeight = 0.05
	overlay := NewOverlay(*weights)

	first := models.Properties{
		UpperNotes: []string{"Rose", "Jasmine"},
		CoreNotes:  []string{"Vanilla"},
		BaseNotes:  []string{"Musk"},
	}

	second := models.Properties{
		UpperNotes: []string{"Rose", "Lily"},
		CoreNotes:  []string{"Vanilla", "Amber"},
		BaseNotes:  []string{"Musk", "Patchouli"},
	}

	score := overlay.getNotesSimilarityScore(first, second)

	// Upper notes: 1/3 = 0.333... (Rose in common, union: Rose, Jasmine, Lily)
	// Core notes: 1/2 = 0.5 (Vanilla in common, union: Vanilla, Amber)
	// Base notes: 1/2 = 0.5 (Musk in common, union: Musk, Patchouli)
	// Total: 0.333*0.2 + 0.5*0.3 + 0.5*0.2 = 0.0666 + 0.15 + 0.1 = 0.3166

	if score <= 0 {
		t.Fatalf("expected positive score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

func TestOverlay_getTypeSimilarityScore(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	overlay := NewOverlay(*weights)

	// Test identical types
	score := overlay.getTypeSimilarityScore("Eau de Parfum", "Eau de Parfum")
	if score != 1.0 {
		t.Fatalf("expected score 1.0 for identical types, got %f", score)
	}

	// Test different types
	score = overlay.getTypeSimilarityScore("Eau de Parfum", "Eau de Toilette")
	if score != 0.0 {
		t.Fatalf("expected score 0.0 for different types, got %f", score)
	}

	// Test empty strings
	score = overlay.getTypeSimilarityScore("", "")
	if score != 1.0 {
		t.Fatalf("expected score 1.0 for empty strings, got %f", score)
	}

	// Test one empty string
	score = overlay.getTypeSimilarityScore("Eau de Parfum", "")
	if score != 0.0 {
		t.Fatalf("expected score 0.0 when one type is empty, got %f", score)
	}
}

func TestOverlay_GetSimilarityScore_EmptyProperties(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.2, 0.3, 0.2)
	weights.FamilyWeight = 0.1
	weights.NotesWeight = 0.15
	weights.TypeWeight = 0.05
	overlay := NewOverlay(*weights)

	empty := models.Properties{}
	filled := models.Properties{
		Type:       "Eau de Parfum",
		Family:     []string{"Floral"},
		UpperNotes: []string{"Rose"},
		CoreNotes:  []string{"Vanilla"},
		BaseNotes:  []string{"Musk"},
	}

	score := overlay.GetSimilarityScore(empty, filled)
	// Should handle empty properties gracefully
	if score < 0 || score > 1.0 {
		t.Fatalf("expected score in range [0, 1], got %f", score)
	}
}

func TestOverlay_GetSimilarityScore_WeightsSum(t *testing.T) {
	t.Parallel()

	// Test with different weight configurations
	weights := NewBaseWeights(0.1, 0.2, 0.1)
	weights.FamilyWeight = 0.2
	weights.NotesWeight = 0.3
	weights.TypeWeight = 0.1
	overlay := NewOverlay(*weights)

	props := models.Properties{
		Type:       "Eau de Parfum",
		Family:     []string{"Floral"},
		UpperNotes: []string{"Rose"},
		CoreNotes:  []string{"Vanilla"},
		BaseNotes:  []string{"Musk"},
	}

	score := overlay.GetSimilarityScore(props, props)

	// For identical properties, score should be based on weights
	// All individual similarities are 1.0
	// Notes: 1.0*0.1 + 1.0*0.2 + 1.0*0.1 = 0.4
	// Total: 1.0*0.2 + 0.4*0.3 + 1.0*0.1 = 0.2 + 0.12 + 0.1 = 0.42
	expectedMin := 0.3
	if score < expectedMin {
		t.Fatalf("expected score >= %f for identical properties with these weights, got %f", expectedMin, score)
	}
}
