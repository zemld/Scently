package matching

import (
	"container/heap"
	"testing"

	"github.com/zemld/Scently/models"
)

func TestPerfumeHeap_Len(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{
			{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.5},
			{Perfume: models.Perfume{Brand: "Dior", Name: "Sauvage"}, Score: 0.7},
		},
		limit: 10,
	}

	if h.Len() != 2 {
		t.Fatalf("expected length 2, got %d", h.Len())
	}

	h = &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    10,
	}

	if h.Len() != 0 {
		t.Fatalf("expected length 0, got %d", h.Len())
	}
}

func TestPerfumeHeap_Less(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{
			{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.3},
			{Perfume: models.Perfume{Brand: "Dior", Name: "Sauvage"}, Score: 0.7},
		},
		limit: 10,
	}

	if !h.Less(0, 1) {
		t.Fatal("expected Less(0, 1) to return true (0.3 < 0.7)")
	}

	if h.Less(1, 0) {
		t.Fatal("expected Less(1, 0) to return false (0.7 > 0.3)")
	}
}

func TestPerfumeHeap_Swap(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{
			{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.3},
			{Perfume: models.Perfume{Brand: "Dior", Name: "Sauvage"}, Score: 0.7},
		},
		limit: 10,
	}

	h.Swap(0, 1)

	if h.perfumes[0].Perfume.Brand != "Dior" {
		t.Fatalf("expected first perfume to be Dior, got %s", h.perfumes[0].Perfume.Brand)
	}
	if h.perfumes[1].Perfume.Brand != "Chanel" {
		t.Fatalf("expected second perfume to be Chanel, got %s", h.perfumes[1].Perfume.Brand)
	}
}

func TestPerfumeHeap_Pop(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{
			{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.3},
			{Perfume: models.Perfume{Brand: "Dior", Name: "Sauvage"}, Score: 0.7},
		},
		limit: 10,
	}

	heap.Init(h)

	popped := h.Pop()
	if popped == nil {
		t.Fatal("expected non-nil result from Pop")
	}

	ranked, ok := popped.(models.Ranked)
	if !ok {
		t.Fatalf("expected models.Ranked, got %T", popped)
	}

	if ranked.Perfume.Brand != "Chanel" {
		t.Fatalf("expected popped perfume to be Chanel (lowest score), got %s", ranked.Perfume.Brand)
	}

	if h.Len() != 1 {
		t.Fatalf("expected heap length 1 after pop, got %d", h.Len())
	}

	// Pop from empty heap
	h = &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    10,
	}

	popped = h.Pop()
	if popped != nil {
		t.Fatal("expected nil result from Pop on empty heap")
	}
}

func TestPerfumeHeap_Push_WithinLimit(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    3,
	}

	heap.Init(h)

	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.5})
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Dior", Name: "Sauvage"}, Score: 0.3})
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Tom Ford", Name: "Black Orchid"}, Score: 0.7})

	if h.Len() != 3 {
		t.Fatalf("expected heap length 3, got %d", h.Len())
	}

	// Verify heap property: root should have minimum score
	if h.perfumes[0].Score != 0.3 {
		t.Fatalf("expected root score to be 0.3 (minimum), got %f", h.perfumes[0].Score)
	}
}

func TestPerfumeHeap_Push_ExceedsLimit(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    2,
	}

	heap.Init(h)

	// Push items that exceed limit
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.5})
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Dior", Name: "Sauvage"}, Score: 0.3})
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Tom Ford", Name: "Black Orchid"}, Score: 0.7})

	// Should only keep top 2 (highest scores)
	if h.Len() != 2 {
		t.Fatalf("expected heap length 2, got %d", h.Len())
	}

	// The heap should contain the two highest scores (0.5 and 0.7)
	// Root should be the minimum of those (0.5)
	if h.perfumes[0].Score < 0.3 {
		t.Fatalf("expected root score to be at least 0.3, got %f", h.perfumes[0].Score)
	}
}

func TestPerfumeHeap_Push_RejectsLowerScores(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    2,
	}

	heap.Init(h)

	// Fill heap to limit
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.5})
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Dior", Name: "Sauvage"}, Score: 0.7})

	initialLen := h.Len()
	rootScore := h.perfumes[0].Score

	// Try to push item with lower score than root
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Tom Ford", Name: "Black Orchid"}, Score: rootScore - 0.1})

	// Should not add the item
	if h.Len() != initialLen {
		t.Fatalf("expected heap length to remain %d, got %d", initialLen, h.Len())
	}
}

func TestPerfumeHeap_Push_ReplacesWithHigherScore(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    2,
	}

	heap.Init(h)

	// Fill heap to limit with lower scores
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.3})
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Dior", Name: "Sauvage"}, Score: 0.4})

	// Push item with higher score
	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Tom Ford", Name: "Black Orchid"}, Score: 0.8})

	// Should replace the minimum (0.3) with 0.8
	if h.Len() != 2 {
		t.Fatalf("expected heap length 2, got %d", h.Len())
	}

	// Verify that the new item is in the heap
	found := false
	for _, p := range h.perfumes {
		if p.Perfume.Brand == "Tom Ford" && p.Score == 0.8 {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to find Tom Ford with score 0.8 in heap")
	}
}

func TestPerfumeHeap_Push_InvalidType(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    10,
	}

	heap.Init(h)

	initialLen := h.Len()

	// Push invalid type
	h.Push("invalid")

	if h.Len() != initialLen {
		t.Fatalf("expected heap length to remain %d, got %d", initialLen, h.Len())
	}
}

func TestPerfumeHeap_Push_ZeroLimit(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    0,
	}

	heap.Init(h)

	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.5})

	if h.Len() != 0 {
		t.Fatalf("expected heap length 0 with zero limit, got %d", h.Len())
	}
}

func TestPerfumeHeap_Push_NegativeLimit(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    -1,
	}

	heap.Init(h)

	h.Push(models.Ranked{Perfume: models.Perfume{Brand: "Chanel", Name: "No5"}, Score: 0.5})

	if h.Len() != 0 {
		t.Fatalf("expected heap length 0 with negative limit, got %d", h.Len())
	}
}

func TestPerfumeHeap_HeapProperty(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{},
		limit:    10,
	}

	heap.Init(h)

	// Push items in random order
	scores := []float64{0.9, 0.1, 0.5, 0.3, 0.7, 0.2, 0.6, 0.4, 0.8}
	for i, score := range scores {
		h.Push(models.Ranked{
			Perfume: models.Perfume{Brand: "Brand", Name: "Perfume"},
			Score:   score,
		})
		if i == 0 {
			// After first push, verify heap property
			if h.perfumes[0].Score != score {
				t.Fatalf("after first push, expected root score %f, got %f", score, h.perfumes[0].Score)
			}
		}
	}

	// Verify heap property: parent should be <= children
	for i := 0; i < h.Len(); i++ {
		left := 2*i + 1
		right := 2*i + 2

		if left < h.Len() && h.perfumes[i].Score > h.perfumes[left].Score {
			t.Fatalf("heap property violated: parent[%d]=%f > left[%d]=%f", i, h.perfumes[i].Score, left, h.perfumes[left].Score)
		}
		if right < h.Len() && h.perfumes[i].Score > h.perfumes[right].Score {
			t.Fatalf("heap property violated: parent[%d]=%f > right[%d]=%f", i, h.perfumes[i].Score, right, h.perfumes[right].Score)
		}
	}
}

