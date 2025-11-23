package matching

import (
	"container/heap"
	"sync"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

type PerfumeHeap struct {
	mu       sync.RWMutex
	perfumes []models.WithScore
}

func (h *PerfumeHeap) Push(x any) {
	h.perfumes = append(h.perfumes, x.(models.WithScore))
}

func (h *PerfumeHeap) Pop() any {
	n := len(h.perfumes)
	x := h.perfumes[n-1]
	h.perfumes = h.perfumes[0 : n-1]
	return x
}

func (h *PerfumeHeap) Len() int {
	return len(h.perfumes)
}

func (h *PerfumeHeap) Less(i, j int) bool {
	return h.perfumes[i].Score < h.perfumes[j].Score
}

func (h *PerfumeHeap) Swap(i, j int) {
	h.perfumes[i], h.perfumes[j] = h.perfumes[j], h.perfumes[i]
}

func (h *PerfumeHeap) PopSafe() models.WithScore {
	h.mu.Lock()
	defer h.mu.Unlock()

	return heap.Pop(h).(models.WithScore)
}

func (h *PerfumeHeap) PushSafeIfNeeded(p models.WithScore, limit int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if limit <= 0 {
		return
	}

	if h.Len() < limit {
		heap.Push(h, p)
		return
	}

	if len(h.perfumes) == 0 || p.Score <= h.perfumes[0].Score {
		return
	}
	heap.Pop(h)
	heap.Push(h, p)
}
