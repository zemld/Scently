package matching

import (
	"container/heap"
	"sync"
	"testing"

	"github.com/zemld/Scently/models"
)

func TestPerfumeHeap_Len(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{}
	if h.Len() != 0 {
		t.Fatalf("expected length 0, got %d", h.Len())
	}

	h.perfumes = []models.Ranked{
		{Perfume: models.Perfume{Name: "Test1"}, Score: 0.5},
		{Perfume: models.Perfume{Name: "Test2"}, Score: 0.7},
	}
	if h.Len() != 2 {
		t.Fatalf("expected length 2, got %d", h.Len())
	}
}

func TestPerfumeHeap_Less(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{
			{Perfume: models.Perfume{Name: "Test1"}, Score: 0.5},
			{Perfume: models.Perfume{Name: "Test2"}, Score: 0.7},
		},
	}

	if !h.Less(0, 1) {
		t.Fatal("expected 0.5 < 0.7 to be true (min-heap)")
	}
	// Less(1, 0) should be false (0.7 < 0.5 is false)
	if h.Less(1, 0) {
		t.Fatal("expected 0.7 < 0.5 to be false (min-heap)")
	}
}

func TestPerfumeHeap_Swap(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{
			{Perfume: models.Perfume{Name: "Test1"}, Score: 0.5},
			{Perfume: models.Perfume{Name: "Test2"}, Score: 0.7},
		},
	}

	h.Swap(0, 1)

	if h.perfumes[0].Score != 0.7 {
		t.Fatalf("expected score 0.7 at index 0, got %f", h.perfumes[0].Score)
	}
	if h.perfumes[1].Score != 0.5 {
		t.Fatalf("expected score 0.5 at index 1, got %f", h.perfumes[1].Score)
	}
}

func TestPerfumeHeap_Push(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{}
	item := models.Ranked{
		Perfume: models.Perfume{Name: "Test1"},
		Score:   0.5,
	}

	h.Push(item)

	if h.Len() != 1 {
		t.Fatalf("expected length 1, got %d", h.Len())
	}
	if h.perfumes[0].Score != 0.5 {
		t.Fatalf("expected score 0.5, got %f", h.perfumes[0].Score)
	}
}

func TestPerfumeHeap_Pop(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{
			{Perfume: models.Perfume{Name: "Test1"}, Score: 0.5},
			{Perfume: models.Perfume{Name: "Test2"}, Score: 0.7},
		},
	}

	item := h.Pop().(models.Ranked)

	if h.Len() != 1 {
		t.Fatalf("expected length 1 after pop, got %d", h.Len())
	}
	if item.Score != 0.7 {
		t.Fatalf("expected popped score 0.7, got %f", item.Score)
	}
}

func TestPerfumeHeap_HeapOperations(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{}
	heap.Init(h)

	items := []models.Ranked{
		{Perfume: models.Perfume{Name: "Test1"}, Score: 0.9},
		{Perfume: models.Perfume{Name: "Test2"}, Score: 0.3},
		{Perfume: models.Perfume{Name: "Test3"}, Score: 0.7},
		{Perfume: models.Perfume{Name: "Test4"}, Score: 0.1},
		{Perfume: models.Perfume{Name: "Test5"}, Score: 0.5},
	}

	for _, item := range items {
		heap.Push(h, item)
	}

	if h.Len() != len(items) {
		t.Fatalf("expected length %d, got %d", len(items), h.Len())
	}

	expectedOrder := []float64{0.1, 0.3, 0.5, 0.7, 0.9}
	for i, expectedScore := range expectedOrder {
		if h.Len() == 0 {
			t.Fatalf("heap is empty at index %d", i)
		}
		item := heap.Pop(h).(models.Ranked)
		if item.Score != expectedScore {
			t.Fatalf("at index %d: expected score %f, got %f", i, expectedScore, item.Score)
		}
	}
}

func TestPerfumeHeap_PushSafe(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{}
	item := models.Ranked{
		Perfume: models.Perfume{Name: "Test1"},
		Score:   0.5,
	}

	h.PushSafeIfNeeded(item, 1)

	if h.Len() != 1 {
		t.Fatalf("expected length 1, got %d", h.Len())
	}
	if h.perfumes[0].Score != 0.5 {
		t.Fatalf("expected score 0.5, got %f", h.perfumes[0].Score)
	}
}

func TestPerfumeHeap_PopSafe(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{
		perfumes: []models.Ranked{
			{Perfume: models.Perfume{Name: "Test1"}, Score: 0.5},
			{Perfume: models.Perfume{Name: "Test2"}, Score: 0.7},
		},
	}

	initialLen := h.Len()
	h.PopSafe()

	if h.Len() != initialLen-1 {
		t.Fatalf("expected length %d after PopSafe, got %d", initialLen-1, h.Len())
	}
}

func TestPerfumeHeap_ConcurrentPushSafe(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{}
	const goroutines = 10
	const itemsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				item := models.Ranked{
					Perfume: models.Perfume{Name: "Test"},
					Score:   float64(id*itemsPerGoroutine + j),
				}
				h.PushSafeIfNeeded(item, 1)
			}
		}(i)
	}

	wg.Wait()

	expectedLen := 1
	if h.Len() != expectedLen {
		t.Fatalf("expected length %d after concurrent pushes, got %d", expectedLen, h.Len())
	}
}

func TestPerfumeHeap_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	h := &PerfumeHeap{}
	const goroutines = 5
	const itemsPerGoroutine = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Concurrent pushes
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				item := models.Ranked{
					Perfume: models.Perfume{Name: "Test"},
					Score:   float64(id*itemsPerGoroutine + j),
				}
				h.PushSafeIfNeeded(item, 1)
			}
		}(i)
	}

	wg.Wait()

	expectedLen := 1
	if h.Len() != expectedLen {
		t.Fatalf("expected length %d, got %d", expectedLen, h.Len())
	}

	// Concurrent pops
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				if h.Len() > 0 {
					h.PopSafe()
				}
			}
		}()
	}

	wg.Wait()

	if h.Len() != 0 {
		t.Fatalf("expected empty heap after concurrent pops, got length %d", h.Len())
	}
}
