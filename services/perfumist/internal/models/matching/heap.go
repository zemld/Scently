package matching

import (
	"container/heap"
	"sync"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type PerfumeHeap struct {
	mu       sync.RWMutex
	perfumes []perfume.WithScore
}

func (h *PerfumeHeap) Push(x any) {
	h.perfumes = append(h.perfumes, x.(perfume.WithScore))
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
	return h.perfumes[i].Score > h.perfumes[j].Score
}

func (h *PerfumeHeap) Swap(i, j int) {
	h.perfumes[i], h.perfumes[j] = h.perfumes[j], h.perfumes[i]
}

func (h *PerfumeHeap) PushSafe(x perfume.WithScore) {
	h.mu.Lock()
	defer h.mu.Unlock()

	heap.Push(h, x)
}

func (h *PerfumeHeap) PopSafe() perfume.WithScore {
	h.mu.Lock()
	defer h.mu.Unlock()

	return heap.Pop(h).(perfume.WithScore)
}
