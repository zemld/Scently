package matching

import "testing"

func TestNewBaseWeights(t *testing.T) {
	t.Parallel()

	upperNotesWeight := 0.3
	coreNotesWeight := 0.4
	baseNotesWeight := 0.3

	weights := NewBaseWeights(upperNotesWeight, coreNotesWeight, baseNotesWeight)

	if weights == nil {
		t.Fatal("expected non-nil weights")
	}
	if weights.UpperNotesWeight != upperNotesWeight {
		t.Fatalf("expected UpperNotesWeight %f, got %f", upperNotesWeight, weights.UpperNotesWeight)
	}
	if weights.CoreNotesWeight != coreNotesWeight {
		t.Fatalf("expected CoreNotesWeight %f, got %f", coreNotesWeight, weights.CoreNotesWeight)
	}
	if weights.BaseNotesWeight != baseNotesWeight {
		t.Fatalf("expected BaseNotesWeight %f, got %f", baseNotesWeight, weights.BaseNotesWeight)
	}
	// Other weights should be zero
	if weights.FamilyWeight != 0 {
		t.Fatalf("expected FamilyWeight 0, got %f", weights.FamilyWeight)
	}
	if weights.NotesWeight != 0 {
		t.Fatalf("expected NotesWeight 0, got %f", weights.NotesWeight)
	}
	if weights.TypeWeight != 0 {
		t.Fatalf("expected TypeWeight 0, got %f", weights.TypeWeight)
	}
	if weights.CharacteristicsWeight != 0 {
		t.Fatalf("expected CharacteristicsWeight 0, got %f", weights.CharacteristicsWeight)
	}
	if weights.TagsWeight != 0 {
		t.Fatalf("expected TagsWeight 0, got %f", weights.TagsWeight)
	}
	if weights.OverlayWeight != 0 {
		t.Fatalf("expected OverlayWeight 0, got %f", weights.OverlayWeight)
	}
}

func TestNewOverlayWeights(t *testing.T) {
	t.Parallel()

	upperNotesWeight := 0.2
	coreNotesWeight := 0.3
	baseNotesWeight := 0.2
	familyWeight := 0.1
	notesWeight := 0.15
	typeWeight := 0.05

	weights := NewOverlayWeights(
		upperNotesWeight,
		coreNotesWeight,
		baseNotesWeight,
		familyWeight,
		notesWeight,
		typeWeight,
	)

	if weights == nil {
		t.Fatal("expected non-nil weights")
	}
	// Base weights should be set
	if weights.UpperNotesWeight != upperNotesWeight {
		t.Fatalf("expected UpperNotesWeight %f, got %f", upperNotesWeight, weights.UpperNotesWeight)
	}
	if weights.CoreNotesWeight != coreNotesWeight {
		t.Fatalf("expected CoreNotesWeight %f, got %f", coreNotesWeight, weights.CoreNotesWeight)
	}
	if weights.BaseNotesWeight != baseNotesWeight {
		t.Fatalf("expected BaseNotesWeight %f, got %f", baseNotesWeight, weights.BaseNotesWeight)
	}
	// Overlay weights should be set
	if weights.FamilyWeight != familyWeight {
		t.Fatalf("expected FamilyWeight %f, got %f", familyWeight, weights.FamilyWeight)
	}
	if weights.NotesWeight != notesWeight {
		t.Fatalf("expected NotesWeight %f, got %f", notesWeight, weights.NotesWeight)
	}
	if weights.TypeWeight != typeWeight {
		t.Fatalf("expected TypeWeight %f, got %f", typeWeight, weights.TypeWeight)
	}
	// Other weights should be zero
	if weights.CharacteristicsWeight != 0 {
		t.Fatalf("expected CharacteristicsWeight 0, got %f", weights.CharacteristicsWeight)
	}
	if weights.TagsWeight != 0 {
		t.Fatalf("expected TagsWeight 0, got %f", weights.TagsWeight)
	}
	if weights.OverlayWeight != 0 {
		t.Fatalf("expected OverlayWeight 0, got %f", weights.OverlayWeight)
	}
}

func TestNewSmartWeights(t *testing.T) {
	t.Parallel()

	upperNotesWeight := 0.2
	coreNotesWeight := 0.3
	baseNotesWeight := 0.2
	characteristicsWeight := 0.15
	tagsWeight := 0.15

	weights := NewSmartWeights(
		upperNotesWeight,
		coreNotesWeight,
		baseNotesWeight,
		characteristicsWeight,
		tagsWeight,
	)

	if weights == nil {
		t.Fatal("expected non-nil weights")
	}
	// Base weights should be set
	if weights.UpperNotesWeight != upperNotesWeight {
		t.Fatalf("expected UpperNotesWeight %f, got %f", upperNotesWeight, weights.UpperNotesWeight)
	}
	if weights.CoreNotesWeight != coreNotesWeight {
		t.Fatalf("expected CoreNotesWeight %f, got %f", coreNotesWeight, weights.CoreNotesWeight)
	}
	if weights.BaseNotesWeight != baseNotesWeight {
		t.Fatalf("expected BaseNotesWeight %f, got %f", baseNotesWeight, weights.BaseNotesWeight)
	}
	// Smart weights should be set
	if weights.CharacteristicsWeight != characteristicsWeight {
		t.Fatalf("expected CharacteristicsWeight %f, got %f", characteristicsWeight, weights.CharacteristicsWeight)
	}
	if weights.TagsWeight != tagsWeight {
		t.Fatalf("expected TagsWeight %f, got %f", tagsWeight, weights.TagsWeight)
	}
	// Other weights should be zero
	if weights.FamilyWeight != 0 {
		t.Fatalf("expected FamilyWeight 0, got %f", weights.FamilyWeight)
	}
	if weights.NotesWeight != 0 {
		t.Fatalf("expected NotesWeight 0, got %f", weights.NotesWeight)
	}
	if weights.TypeWeight != 0 {
		t.Fatalf("expected TypeWeight 0, got %f", weights.TypeWeight)
	}
	if weights.OverlayWeight != 0 {
		t.Fatalf("expected OverlayWeight 0, got %f", weights.OverlayWeight)
	}
}

func TestNewSmartEnhancedWeights(t *testing.T) {
	t.Parallel()

	upperNotesWeight := 0.15
	coreNotesWeight := 0.2
	baseNotesWeight := 0.15
	characteristicsWeight := 0.2
	tagsWeight := 0.2
	overlayWeight := 0.1

	weights := NewSmartEnhancedWeights(
		upperNotesWeight,
		coreNotesWeight,
		baseNotesWeight,
		characteristicsWeight,
		tagsWeight,
		overlayWeight,
	)

	if weights == nil {
		t.Fatal("expected non-nil weights")
	}
	// Base weights should be set
	if weights.UpperNotesWeight != upperNotesWeight {
		t.Fatalf("expected UpperNotesWeight %f, got %f", upperNotesWeight, weights.UpperNotesWeight)
	}
	if weights.CoreNotesWeight != coreNotesWeight {
		t.Fatalf("expected CoreNotesWeight %f, got %f", coreNotesWeight, weights.CoreNotesWeight)
	}
	if weights.BaseNotesWeight != baseNotesWeight {
		t.Fatalf("expected BaseNotesWeight %f, got %f", baseNotesWeight, weights.BaseNotesWeight)
	}
	// Smart weights should be set
	if weights.CharacteristicsWeight != characteristicsWeight {
		t.Fatalf("expected CharacteristicsWeight %f, got %f", characteristicsWeight, weights.CharacteristicsWeight)
	}
	if weights.TagsWeight != tagsWeight {
		t.Fatalf("expected TagsWeight %f, got %f", tagsWeight, weights.TagsWeight)
	}
	// Enhanced weight should be set
	if weights.OverlayWeight != overlayWeight {
		t.Fatalf("expected OverlayWeight %f, got %f", overlayWeight, weights.OverlayWeight)
	}
	// Other weights should be zero
	if weights.FamilyWeight != 0 {
		t.Fatalf("expected FamilyWeight 0, got %f", weights.FamilyWeight)
	}
	if weights.NotesWeight != 0 {
		t.Fatalf("expected NotesWeight 0, got %f", weights.NotesWeight)
	}
	if weights.TypeWeight != 0 {
		t.Fatalf("expected TypeWeight 0, got %f", weights.TypeWeight)
	}
}

func TestNewOverlayWeights_UsesNewBaseWeights(t *testing.T) {
	t.Parallel()

	// Verify that NewOverlayWeights calls NewBaseWeights internally
	weights := NewOverlayWeights(0.1, 0.2, 0.3, 0.15, 0.15, 0.1)

	// Should have base weights set
	if weights.UpperNotesWeight != 0.1 {
		t.Fatalf("expected UpperNotesWeight to be set by NewBaseWeights")
	}
	if weights.CoreNotesWeight != 0.2 {
		t.Fatalf("expected CoreNotesWeight to be set by NewBaseWeights")
	}
	if weights.BaseNotesWeight != 0.3 {
		t.Fatalf("expected BaseNotesWeight to be set by NewBaseWeights")
	}
}

func TestNewSmartWeights_UsesNewBaseWeights(t *testing.T) {
	t.Parallel()

	// Verify that NewSmartWeights calls NewBaseWeights internally
	weights := NewSmartWeights(0.1, 0.2, 0.3, 0.2, 0.2)

	// Should have base weights set
	if weights.UpperNotesWeight != 0.1 {
		t.Fatalf("expected UpperNotesWeight to be set by NewBaseWeights")
	}
	if weights.CoreNotesWeight != 0.2 {
		t.Fatalf("expected CoreNotesWeight to be set by NewBaseWeights")
	}
	if weights.BaseNotesWeight != 0.3 {
		t.Fatalf("expected BaseNotesWeight to be set by NewBaseWeights")
	}
}

func TestNewSmartEnhancedWeights_UsesNewSmartWeights(t *testing.T) {
	t.Parallel()

	// Verify that NewSmartEnhancedWeights calls NewSmartWeights internally
	weights := NewSmartEnhancedWeights(0.1, 0.2, 0.3, 0.15, 0.15, 0.1)

	// Should have base weights set
	if weights.UpperNotesWeight != 0.1 {
		t.Fatalf("expected UpperNotesWeight to be set by NewSmartWeights")
	}
	if weights.CoreNotesWeight != 0.2 {
		t.Fatalf("expected CoreNotesWeight to be set by NewSmartWeights")
	}
	if weights.BaseNotesWeight != 0.3 {
		t.Fatalf("expected BaseNotesWeight to be set by NewSmartWeights")
	}
	// Should have smart weights set
	if weights.CharacteristicsWeight != 0.15 {
		t.Fatalf("expected CharacteristicsWeight to be set by NewSmartWeights")
	}
	if weights.TagsWeight != 0.15 {
		t.Fatalf("expected TagsWeight to be set by NewSmartWeights")
	}
}

