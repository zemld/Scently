package matching

import (
	"container/heap"

	"github.com/zemld/Scently/models"
)

type PerfumeHeap struct {
	perfumes []models.Ranked
	limit    int
}

func (h *PerfumeHeap) Push(x any) {
	if h.limit <= 0 {
		return
	}

	ranked, ok := x.(models.Ranked)
	if !ok {
		return
	}

	if h.Len() < h.limit {
		h.perfumes = append(h.perfumes, ranked)
		heap.Fix(h, h.Len()-1)
		return
	}

	if ranked.Score <= h.perfumes[0].Score {
		return
	}
	heap.Pop(h)
	h.perfumes = append(h.perfumes, ranked)
	heap.Fix(h, h.Len()-1)
}

func (h *PerfumeHeap) Pop() any {
	if h.Len() == 0 {
		return nil
	}
	x := h.perfumes[0]
	h.perfumes = h.perfumes[1:]
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
