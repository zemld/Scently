package matching

import (
	"testing"

	"github.com/zemld/Scently/models"
)

func TestNewCharacteristicsMatcher(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	matcher := NewCharacteristicsMatcher(*weights)

	if matcher == nil {
		t.Fatal("expected non-nil CharacteristicsMatcher")
	}
	if matcher.UpperNotesWeight != 0.3 {
		t.Fatalf("expected UpperNotesWeight 0.3, got %f", matcher.UpperNotesWeight)
	}
	if matcher.CoreNotesWeight != 0.4 {
		t.Fatalf("expected CoreNotesWeight 0.4, got %f", matcher.CoreNotesWeight)
	}
	if matcher.BaseNotesWeight != 0.3 {
		t.Fatalf("expected BaseNotesWeight 0.3, got %f", matcher.BaseNotesWeight)
	}
}

func TestCharacteristicsMatcher_GetPerfumeSimilarityScore(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	matcher := NewCharacteristicsMatcher(*weights)

	first := models.Properties{
		UpperCharacteristics: map[string]float64{
			"fresh": 0.5,
			"citrus": 0.3,
		},
		CoreCharacteristics: map[string]float64{
			"floral": 0.6,
			"spicy":  0.2,
		},
		BaseCharacteristics: map[string]float64{
			"woody": 0.4,
			"musk":  0.3,
		},
	}

	second := models.Properties{
		UpperCharacteristics: map[string]float64{
			"fresh": 0.4,
			"citrus": 0.4,
		},
		CoreCharacteristics: map[string]float64{
			"floral": 0.5,
			"spicy":  0.3,
		},
		BaseCharacteristics: map[string]float64{
			"woody": 0.5,
			"musk":  0.2,
		},
	}

	score := matcher.GetPerfumeSimilarityScore(first, second)

	if score <= 0 {
		t.Fatalf("expected positive score, got %f", score)
	}
	if score > 1.0 {
		t.Fatalf("expected score <= 1.0, got %f", score)
	}
}

func TestCharacteristicsMatcher_GetPerfumeSimilarityScore_Identical(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	matcher := NewCharacteristicsMatcher(*weights)

	props := models.Properties{
		UpperCharacteristics: map[string]float64{
			"fresh": 0.5,
			"citrus": 0.3,
		},
		CoreCharacteristics: map[string]float64{
			"floral": 0.6,
		},
		BaseCharacteristics: map[string]float64{
			"woody": 0.4,
		},
	}

	score := matcher.GetPerfumeSimilarityScore(props, props)

	// For identical characteristics, cosine similarity should be 1.0
	// Score should be: 1.0*0.3 + 1.0*0.4 + 1.0*0.3 = 1.0
	expected := 1.0
	if score != expected {
		t.Fatalf("expected score %f for identical properties, got %f", expected, score)
	}
}

func TestCharacteristicsMatcher_GetPerfumeSimilarityScore_CompletelyDifferent(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	matcher := NewCharacteristicsMatcher(*weights)

	first := models.Properties{
		UpperCharacteristics: map[string]float64{
			"fresh": 0.5,
		},
		CoreCharacteristics: map[string]float64{
			"floral": 0.6,
		},
		BaseCharacteristics: map[string]float64{
			"woody": 0.4,
		},
	}

	second := models.Properties{
		UpperCharacteristics: map[string]float64{
			"spicy": 0.5,
		},
		CoreCharacteristics: map[string]float64{
			"oriental": 0.6,
		},
		BaseCharacteristics: map[string]float64{
			"amber": 0.4,
		},
	}

	score := matcher.GetPerfumeSimilarityScore(first, second)

	// No overlapping characteristics, cosine similarity should be 0.0
	expected := 0.0
	if score != expected {
		t.Fatalf("expected score %f for completely different properties, got %f", expected, score)
	}
}

func TestCharacteristicsMatcher_GetPerfumeSimilarityScore_EmptyCharacteristics(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	matcher := NewCharacteristicsMatcher(*weights)

	empty := models.Properties{
		UpperCharacteristics: map[string]float64{},
		CoreCharacteristics:  map[string]float64{},
		BaseCharacteristics:  map[string]float64{},
	}

	filled := models.Properties{
		UpperCharacteristics: map[string]float64{"fresh": 0.5},
		CoreCharacteristics:  map[string]float64{"floral": 0.6},
		BaseCharacteristics:  map[string]float64{"woody": 0.4},
	}

	score := matcher.GetPerfumeSimilarityScore(empty, filled)

	// Empty characteristics should result in 0.0 cosine similarity
	expected := 0.0
	if score != expected {
		t.Fatalf("expected score %f for empty characteristics, got %f", expected, score)
	}
}

func TestCharacteristicsMatcher_GetPerfumeSimilarityScore_PartialOverlap(t *testing.T) {
	t.Parallel()

	weights := NewBaseWeights(0.3, 0.4, 0.3)
	matcher := NewCharacteristicsMatcher(*weights)

	first := models.Properties{
		UpperCharacteristics: map[string]float64{
			"fresh": 0.5,
			"citrus": 0.3,
		},
		CoreCharacteristics: map[string]float64{
			"floral": 0.6,
		},
		BaseCharacteristics: map[string]float64{
			"woody": 0.4,
		},
	}

	second := models.Properties{
		UpperCharacteristics: map[string]float64{
			"fresh": 0.4,
			"spicy": 0.3,
		},
		CoreCharacteristics: map[string]float64{
			"floral": 0.5,
		},
		BaseCharacteristics: map[string]float64{
			"woody": 0.5,
		},
	}

	score := matcher.GetPerfumeSimilarityScore(first, second)

	// Should have partial overlap (fresh, floral, woody)
	if score <= 0 {
		t.Fatalf("expected positive score for partial overlap, got %f", score)
	}
	if score >= 1.0 {
		t.Fatalf("expected score < 1.0 for partial overlap, got %f", score)
	}
}

func TestCharacteristicsMatcher_GetPerfumeSimilarityScore_Weights(t *testing.T) {
	t.Parallel()

	// Test with different weights
	weights := NewBaseWeights(0.5, 0.3, 0.2)
	matcher := NewCharacteristicsMatcher(*weights)

	props := models.Properties{
		UpperCharacteristics: map[string]float64{"fresh": 0.5},
		CoreCharacteristics:  map[string]float64{"floral": 0.6},
		BaseCharacteristics:  map[string]float64{"woody": 0.4},
	}

	score := matcher.GetPerfumeSimilarityScore(props, props)

	// For identical properties, should be 1.0 regardless of weights
	expected := 1.0
	if score != expected {
		t.Fatalf("expected score %f for identical properties, got %f", expected, score)
	}
}

func TestUniteCharacteristics(t *testing.T) {
	t.Parallel()

	notes := []models.EnrichedNote{
		{
			Name: "Rose",
			Characteristics: []models.NoteCharacteristic{
				{Name: "floral", Value: 0.6},
				{Name: "sweet", Value: 0.3},
			},
		},
		{
			Name: "Jasmine",
			Characteristics: []models.NoteCharacteristic{
				{Name: "floral", Value: 0.5},
				{Name: "sweet", Value: 0.4},
			},
		},
	}

	united := uniteCharacteristics(notes)

	// floral: (0.6 + 0.5) / 2 = 0.55
	// sweet: (0.3 + 0.4) / 2 = 0.35
	expectedFloral := 0.55
	expectedSweet := 0.35

	if united["floral"] != expectedFloral {
		t.Fatalf("expected floral %f, got %f", expectedFloral, united["floral"])
	}
	if united["sweet"] != expectedSweet {
		t.Fatalf("expected sweet %f, got %f", expectedSweet, united["sweet"])
	}
}

func TestUniteCharacteristics_EmptyNotes(t *testing.T) {
	t.Parallel()

	notes := []models.EnrichedNote{}
	united := uniteCharacteristics(notes)

	if len(united) != 0 {
		t.Fatalf("expected empty map for empty notes, got %d items", len(united))
	}
}

func TestUniteCharacteristics_SingleNote(t *testing.T) {
	t.Parallel()

	notes := []models.EnrichedNote{
		{
			Name: "Rose",
			Characteristics: []models.NoteCharacteristic{
				{Name: "floral", Value: 0.6},
			},
		},
	}

	united := uniteCharacteristics(notes)

	// Single note, value should be divided by 1 (no change)
	expectedFloral := 0.6
	if united["floral"] != expectedFloral {
		t.Fatalf("expected floral %f, got %f", expectedFloral, united["floral"])
	}
}

func TestPreparePerfumeCharacteristics(t *testing.T) {
	t.Parallel()

	perfume := &models.Perfume{
		Properties: models.Properties{
			EnrichedUpperNotes: []models.EnrichedNote{
				{
					Name: "Rose",
					Characteristics: []models.NoteCharacteristic{
						{Name: "floral", Value: 0.6},
					},
				},
			},
			EnrichedCoreNotes: []models.EnrichedNote{
				{
					Name: "Vanilla",
					Characteristics: []models.NoteCharacteristic{
						{Name: "sweet", Value: 0.5},
					},
				},
			},
			EnrichedBaseNotes: []models.EnrichedNote{
				{
					Name: "Musk",
					Characteristics: []models.NoteCharacteristic{
						{Name: "woody", Value: 0.4},
					},
				},
			},
		},
	}

	preparePerfumeCharacteristics(perfume)

	if perfume.Properties.UpperCharacteristics == nil {
		t.Fatal("expected UpperCharacteristics to be set")
	}
	if perfume.Properties.CoreCharacteristics == nil {
		t.Fatal("expected CoreCharacteristics to be set")
	}
	if perfume.Properties.BaseCharacteristics == nil {
		t.Fatal("expected BaseCharacteristics to be set")
	}

	if perfume.Properties.UpperCharacteristics["floral"] != 0.6 {
		t.Fatalf("expected UpperCharacteristics[floral] 0.6, got %f", perfume.Properties.UpperCharacteristics["floral"])
	}
	if perfume.Properties.CoreCharacteristics["sweet"] != 0.5 {
		t.Fatalf("expected CoreCharacteristics[sweet] 0.5, got %f", perfume.Properties.CoreCharacteristics["sweet"])
	}
	if perfume.Properties.BaseCharacteristics["woody"] != 0.4 {
		t.Fatalf("expected BaseCharacteristics[woody] 0.4, got %f", perfume.Properties.BaseCharacteristics["woody"])
	}
}

