package app

import (
	"math"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
)

func TestGetPerfumeSimilarityScore_EmptyProperties(t *testing.T) {
	t.Parallel()

	a := models.PerfumeProperties{}
	b := models.PerfumeProperties{}
	// Empty properties: families=0, notes=0, type=1 (empty strings are equal)
	// Score = 0*0.4 + 0*0.55 + 1*0.05 = 0.05
	got := GetPerfumeSimilarityScore(a, b)
	if diff := math.Abs(got - 0.05); diff > 1e-9 {
		t.Fatalf("empty properties expected 0.05, got %v", got)
	}
}

func TestGetPerfumeSimilarityScore_FullMatch(t *testing.T) {
	t.Parallel()

	a := models.PerfumeProperties{
		Type:       "edt",
		Family:     []string{"woody", "spicy"},
		UpperNotes: []string{"bergamot"},
		CoreNotes:  []string{"lavender"},
		BaseNotes:  []string{"cedar"},
	}
	b := a
	got := GetPerfumeSimilarityScore(a, b)
	if diff := math.Abs(got - 1.0); diff > 1e-9 {
		t.Fatalf("full match should be 1.0, got %v", got)
	}
}

func TestGetPerfumeSimilarityScore_EmptyLists(t *testing.T) {
	t.Parallel()

	a := models.PerfumeProperties{Type: "edt"}
	b := models.PerfumeProperties{Type: "edt"}
	// No notes/families -> list similarities 0; type equal -> 1
	// Score = 0*0.4 + 0*0.55 + 1*0.05 = 0.05
	got := GetPerfumeSimilarityScore(a, b)
	if diff := math.Abs(got - 0.05); diff > 1e-9 {
		t.Fatalf("empty lists score should be 0.05, got %v", got)
	}
}

func TestGetPerfumeSimilarityScore_PartialNotes(t *testing.T) {
	t.Parallel()

	a := models.PerfumeProperties{
		Type:       "edt",
		Family:     []string{"woody", "spicy"},
		UpperNotes: []string{"bergamot", "lemon"},
		CoreNotes:  []string{"lavender"},
		BaseNotes:  []string{"cedar", "musk"},
	}
	b := models.PerfumeProperties{
		Type:       "edp",
		Family:     []string{"woody", "amber"},
		UpperNotes: []string{"bergamot"},
		CoreNotes:  []string{"lavender", "rose"},
		BaseNotes:  []string{"musk"},
	}

	// families: {woody} / {woody, spicy, amber} -> Jaccard 1/3
	// upper: {bergamot} / {bergamot, lemon} -> 1/2
	// core: {lavender} / {lavender, rose} -> 1/2
	// base: {musk} / {cedar, musk} -> 1/2
	// notes weighted = 0.15*0.5 + 0.45*0.5 + 0.4*0.5 = 0.5
	// type: diff -> 0
	// total = 0.4*(1/3) + 0.55*0.5 + 0.05*0 = 0.133333.. + 0.275 = 0.408333..
	got := GetPerfumeSimilarityScore(a, b)
	if diff := math.Abs(got - 0.4083333333); diff > 1e-9 {
		t.Fatalf("partial notes score unexpected, got %v", got)
	}
}

func TestGetTypeSimilarityScore(t *testing.T) {
	t.Parallel()

	if s := getTypeSimilarityScore("a", "a"); s != 1 {
		t.Fatalf("same type expected 1, got %v", s)
	}
	if s := getTypeSimilarityScore("a", "b"); s != 0 {
		t.Fatalf("different type expected 0, got %v", s)
	}
}

func TestGetListSimilarityScore(t *testing.T) {
	t.Parallel()

	if s := getListSimilarityScore(nil, nil); s != 0 {
		t.Fatalf("empty lists expected 0, got %v", s)
	}
	if s := getListSimilarityScore([]string{"a"}, []string{"a", "b"}); math.Abs(s-0.5) > 1e-9 {
		t.Fatalf("expected 0.5, got %v", s)
	}
	if s := getListSimilarityScore([]string{"a", "b"}, []string{"c"}); s != 0 {
		t.Fatalf("disjoint expected 0, got %v", s)
	}
}

func TestUpdateMostSimilarIfNeeded(t *testing.T) {
	t.Parallel()

	arr := make([]models.PerfumeWithScore, 0, 4)
	for _, s := range []float64{0.1, 0.2, 0.3, 0.4} {
		arr = updateMostSimilarIfNeeded(arr, models.Perfume{Name: "n"}, s)
	}
	if len(arr) != 4 {
		t.Fatalf("expected 4 items, got %d", len(arr))
	}
	for i, expectedScore := range []float64{0.4, 0.3, 0.2, 0.1} {
		if arr[i].Score != expectedScore {
			t.Fatalf("position %d expected score %v, got %v", i, expectedScore, arr[i].Score)
		}
	}
	arr = updateMostSimilarIfNeeded(arr, models.Perfume{Name: "m"}, 0.25)
	if !(arr[1].Score >= 0.25 && arr[2].Score >= 0.25) {
		t.Fatalf("middle insertion not reflected: %+v", arr)
	}
}
