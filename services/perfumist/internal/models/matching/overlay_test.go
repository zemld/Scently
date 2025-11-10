package matching

import (
	"math"
	"sync"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

func TestNewSimpleMatcher(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 4)

	if matcher == nil {
		t.Fatal("expected non-nil matcher")
	}
	if matcher.familyWeight != 0.3 {
		t.Fatalf("expected familyWeight 0.3, got %f", matcher.familyWeight)
	}
	if matcher.notesWeight != 0.4 {
		t.Fatalf("expected notesWeight 0.4, got %f", matcher.notesWeight)
	}
	if matcher.typeWeight != 0.3 {
		t.Fatalf("expected typeWeight 0.3, got %f", matcher.typeWeight)
	}
	if matcher.upperNotesWeight != 0.2 {
		t.Fatalf("expected upperNotesWeight 0.2, got %f", matcher.upperNotesWeight)
	}
	if matcher.middleNotesWeight != 0.3 {
		t.Fatalf("expected middleNotesWeight 0.3, got %f", matcher.middleNotesWeight)
	}
	if matcher.baseNotesWeight != 0.5 {
		t.Fatalf("expected baseNotesWeight 0.5, got %f", matcher.baseNotesWeight)
	}
	if matcher.threadsCount != 4 {
		t.Fatalf("expected threadsCount 4, got %d", matcher.threadsCount)
	}
}

func TestOverlay_getTypeSimilarityScore_Identical(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)
	score := matcher.getTypeSimilarityScore("EDT", "EDT")

	if score != 1.0 {
		t.Fatalf("expected score 1.0 for identical types, got %f", score)
	}
}

func TestOverlay_getTypeSimilarityScore_Different(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)
	score := matcher.getTypeSimilarityScore("EDT", "EDP")

	if score != 0.0 {
		t.Fatalf("expected score 0.0 for different types, got %f", score)
	}
}

func TestOverlay_getListSimilarityScore_Identical(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)
	first := []string{"floral", "woody", "spicy"}
	second := []string{"floral", "woody", "spicy"}

	score := matcher.getListSimilarityScore(first, second)

	if score != 1.0 {
		t.Fatalf("expected score 1.0 for identical lists, got %f", score)
	}
}

func TestOverlay_getListSimilarityScore_PartialOverlap(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)
	first := []string{"floral", "woody"}
	second := []string{"floral", "woody", "spicy"}

	score := matcher.getListSimilarityScore(first, second)

	// Intersection: {"floral", "woody"} = 2
	// Union: {"floral", "woody", "spicy"} = 3
	// Score: 2/3 = 0.666...
	expected := 2.0 / 3.0
	if score != expected {
		t.Fatalf("expected score %f, got %f", expected, score)
	}
}

func TestOverlay_getListSimilarityScore_NoOverlap(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)
	first := []string{"floral", "woody"}
	second := []string{"citrus", "spicy"}

	score := matcher.getListSimilarityScore(first, second)

	if score != 0.0 {
		t.Fatalf("expected score 0.0 for no overlap, got %f", score)
	}
}

func TestOverlay_getListSimilarityScore_EmptyLists(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)
	first := []string{}
	second := []string{}

	score := matcher.getListSimilarityScore(first, second)

	if score != 0.0 {
		t.Fatalf("expected score 0.0 for empty lists, got %f", score)
	}
}

func TestOverlay_getListSimilarityScore_OneEmpty(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)
	first := []string{"floral", "woody"}
	second := []string{}

	score := matcher.getListSimilarityScore(first, second)

	if score != 0.0 {
		t.Fatalf("expected score 0.0 when one list is empty, got %f", score)
	}
}

func TestOverlay_getListSimilarityScore_Duplicates(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)
	first := []string{"floral", "floral", "woody"}
	second := []string{"floral", "woody", "woody"}

	score := matcher.getListSimilarityScore(first, second)

	// Sets: first = {"floral", "woody"}, second = {"floral", "woody"}
	// Intersection: {"floral", "woody"} = 2
	// Union: {"floral", "woody"} = 2
	// Score: 2/2 = 1.0
	if score != 1.0 {
		t.Fatalf("expected score 1.0 (duplicates should be ignored), got %f", score)
	}
}

func TestOverlay_getNotesSimilarityScore(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)

	first := perfume.Properties{
		UpperNotes: []string{"bergamot", "lemon"},
		CoreNotes:  []string{"lavender", "rose"},
		BaseNotes:  []string{"musk", "vanilla"},
	}

	second := perfume.Properties{
		UpperNotes: []string{"bergamot"},
		CoreNotes:  []string{"lavender", "rose"},
		BaseNotes:  []string{"musk"},
	}

	score := matcher.getNotesSimilarityScore(first, second)

	// Upper: {"bergamot"} / {"bergamot", "lemon"} = 1/2 = 0.5
	// Core: {"lavender", "rose"} / {"lavender", "rose"} = 2/2 = 1.0
	// Base: {"musk"} / {"musk", "vanilla"} = 1/2 = 0.5
	// Weighted: 0.5*0.2 + 1.0*0.3 + 0.5*0.5 = 0.1 + 0.3 + 0.25 = 0.65
	expected := 0.5*0.2 + 1.0*0.3 + 0.5*0.5
	if math.Abs(score-expected) > 0.0001 {
		t.Fatalf("expected score %f, got %f", expected, score)
	}
}

func TestOverlay_GetPerfumeSimilarityScore(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)

	first := perfume.Properties{
		Type:       "EDT",
		Family:     []string{"floral", "woody"},
		UpperNotes: []string{"bergamot"},
		CoreNotes:  []string{"rose"},
		BaseNotes:  []string{"musk"},
	}

	second := perfume.Properties{
		Type:       "EDT",
		Family:     []string{"floral"},
		UpperNotes: []string{"bergamot"},
		CoreNotes:  []string{"rose"},
		BaseNotes:  []string{"musk"},
	}

	score := matcher.GetPerfumeSimilarityScore(first, second)

	// Type: 1.0 (identical) * 0.3 = 0.3
	// Family: 1/2 = 0.5 * 0.3 = 0.15
	// Notes: 1.0 (all identical) * 0.4 = 0.4
	// Total: 0.3 + 0.15 + 0.4 = 0.85
	expected := 0.3 + 0.15 + 0.4
	if math.Abs(score-expected) > 0.0001 {
		t.Fatalf("expected score %f, got %f", expected, score)
	}
}

func TestOverlay_GetPerfumeSimilarityScore_DifferentType(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)

	first := perfume.Properties{
		Type:       "EDT",
		Family:     []string{"floral"},
		UpperNotes: []string{"bergamot"},
		CoreNotes:  []string{"rose"},
		BaseNotes:  []string{"musk"},
	}

	second := perfume.Properties{
		Type:       "EDP",
		Family:     []string{"floral"},
		UpperNotes: []string{"bergamot"},
		CoreNotes:  []string{"rose"},
		BaseNotes:  []string{"musk"},
	}

	score := matcher.GetPerfumeSimilarityScore(first, second)

	// Type: 0.0 (different) * 0.3 = 0.0
	// Family: 1.0 (identical) * 0.3 = 0.3
	// Notes: 1.0 (all identical) * 0.4 = 0.4
	// Total: 0.0 + 0.3 + 0.4 = 0.7
	expected := 0.0 + 0.3 + 0.4
	if math.Abs(score-expected) > 0.0001 {
		t.Fatalf("expected score %f, got %f", expected, score)
	}
}

func TestOverlay_buildHeapAsync_SkipsEqual(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)

	favourite := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	all := []perfume.Perfume{
		favourite, // Should be skipped
		{Brand: "Dior", Name: "Sauvage", Sex: "male"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}

	h := &PerfumeHeap{}
	wg := &sync.WaitGroup{}
	wg.Add(1)

	matcher.processPerfumes(h, wg, favourite, all, 2)
	wg.Wait()

	// Should have 2 items (favourite is skipped)
	if h.Len() != 2 {
		t.Fatalf("expected 2 items in heap (favourite skipped), got %d", h.Len())
	}
}

func TestOverlay_Find(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 2)

	favourite := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
		Properties: perfume.Properties{
			Type:       "EDT",
			Family:     []string{"floral"},
			UpperNotes: []string{"bergamot"},
			CoreNotes:  []string{"rose"},
			BaseNotes:  []string{"musk"},
		},
	}

	// Create perfumes with different similarity scores
	// Perfume1: same type, same family, same notes - highest similarity
	perfume1 := perfume.Perfume{
		Brand: "Dior",
		Name:  "J'adore",
		Sex:   "female",
		Properties: perfume.Properties{
			Type:       "EDT",
			Family:     []string{"floral"},
			UpperNotes: []string{"bergamot"},
			CoreNotes:  []string{"rose"},
			BaseNotes:  []string{"musk"},
		},
	}

	// Perfume2: same type, different family - medium similarity
	perfume2 := perfume.Perfume{
		Brand: "Tom Ford",
		Name:  "Black Orchid",
		Sex:   "unisex",
		Properties: perfume.Properties{
			Type:       "EDT",
			Family:     []string{"woody"},
			UpperNotes: []string{"citrus"},
			CoreNotes:  []string{"spicy"},
			BaseNotes:  []string{"amber"},
		},
	}

	// Perfume3: different type, different everything - lowest similarity
	perfume3 := perfume.Perfume{
		Brand: "Versace",
		Name:  "Eros",
		Sex:   "male",
		Properties: perfume.Properties{
			Type:       "EDP",
			Family:     []string{"fresh"},
			UpperNotes: []string{"mint"},
			CoreNotes:  []string{"tonka"},
			BaseNotes:  []string{"vanilla"},
		},
	}

	all := []perfume.Perfume{perfume1, perfume2, perfume3}

	result := matcher.Find(favourite, all, 2)

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}

	// Results should be sorted by score descending (highest first)
	// Since heap is max-heap, first element should have highest score
	if result[0].Rank != 1 {
		t.Fatalf("expected first rank to be 1, got %d", result[0].Rank)
	}
	if result[1].Rank != 2 {
		t.Fatalf("expected second rank to be 2, got %d", result[1].Rank)
	}

	// Verify scores are in descending order
	if result[0].Score < result[1].Score {
		t.Fatalf("expected first score (%f) to be >= second (%f)", result[0].Score, result[1].Score)
	}

	// Verify that perfume1 (identical properties) has highest or equal score
	// Since it has identical properties, it should have score close to 1.0
	if result[0].Score < 0.9 {
		t.Fatalf("expected highest score to be >= 0.9, got %f", result[0].Score)
	}
}

func TestOverlay_Find_EmptyList(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)

	favourite := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	all := []perfume.Perfume{}

	result := matcher.Find(favourite, all, 5)

	if result == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d items", len(result))
	}
}

func TestOverlay_Find_OnlyFavourite(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)

	favourite := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	all := []perfume.Perfume{favourite}

	result := matcher.Find(favourite, all, 5)

	if result == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result (favourite is skipped), got %d items", len(result))
	}
}

func TestOverlay_Find_MoreThanAvailable(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)

	favourite := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
		Properties: perfume.Properties{
			Type:       "EDT",
			Family:     []string{"floral"},
			UpperNotes: []string{"bergamot"},
			CoreNotes:  []string{"rose"},
			BaseNotes:  []string{"musk"},
		},
	}

	all := []perfume.Perfume{
		{Brand: "Dior", Name: "Sauvage", Sex: "male"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}

	// Request more than available
	result := matcher.Find(favourite, all, 10)

	// Should return only available items (2)
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}

	// Verify all results have valid ranks and scores
	for i, r := range result {
		if r.Rank != i+1 {
			t.Fatalf("expected rank %d at index %d, got %d", i+1, i, r.Rank)
		}
		if r.Score < 0 || r.Score > 1 {
			t.Fatalf("expected score between 0 and 1, got %f", r.Score)
		}
	}
}

func TestOverlay_Find_OrderedByScore(t *testing.T) {
	t.Parallel()

	matcher := NewOverlay(0.3, 0.4, 0.3, 0.2, 0.3, 0.5, 1)

	favourite := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
		Properties: perfume.Properties{
			Type:       "EDT",
			Family:     []string{"floral", "woody"},
			UpperNotes: []string{"bergamot", "lemon"},
			CoreNotes:  []string{"rose", "jasmine"},
			BaseNotes:  []string{"musk", "vanilla"},
		},
	}

	// Create perfumes with clearly different similarity scores
	// High similarity - same type, overlapping family and notes
	highSimilar := perfume.Perfume{
		Brand: "Dior",
		Name:  "J'adore",
		Sex:   "female",
		Properties: perfume.Properties{
			Type:       "EDT",
			Family:     []string{"floral", "woody"},
			UpperNotes: []string{"bergamot"},
			CoreNotes:  []string{"rose"},
			BaseNotes:  []string{"musk"},
		},
	}

	// Medium similarity - same type, some overlap
	mediumSimilar := perfume.Perfume{
		Brand: "Tom Ford",
		Name:  "Black Orchid",
		Sex:   "unisex",
		Properties: perfume.Properties{
			Type:       "EDT",
			Family:     []string{"woody"},
			UpperNotes: []string{"citrus"},
			CoreNotes:  []string{"spicy"},
			BaseNotes:  []string{"amber"},
		},
	}

	// Low similarity - different type, no overlap
	lowSimilar := perfume.Perfume{
		Brand: "Versace",
		Name:  "Eros",
		Sex:   "male",
		Properties: perfume.Properties{
			Type:       "EDP",
			Family:     []string{"fresh"},
			UpperNotes: []string{"mint"},
			CoreNotes:  []string{"tonka"},
			BaseNotes:  []string{"vanilla"},
		},
	}

	all := []perfume.Perfume{highSimilar, mediumSimilar, lowSimilar}

	result := matcher.Find(favourite, all, 3)

	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}

	// Verify scores are in descending order
	for i := 0; i < len(result)-1; i++ {
		if result[i].Score < result[i+1].Score {
			t.Fatalf("expected result[%d].Score (%f) >= result[%d].Score (%f)", i, result[i].Score, i+1, result[i+1].Score)
		}
	}

	// Verify highest score is from highSimilar perfume
	if result[0].Score <= result[1].Score || result[0].Score <= result[2].Score {
		t.Fatalf("expected first result to have highest score, got scores: %f, %f, %f", result[0].Score, result[1].Score, result[2].Score)
	}
}
