package matching

import (
	"container/heap"
	"sync"

	"github.com/zemld/Scently/models"
)

type Matcher interface {
	GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64
}

type MatchData struct {
	Matcher
	favourite    models.Perfume
	all          []models.Perfume
	matchesCount int
	threadsCount int
}

func NewMatchData(
	matcher Matcher,
	favourite models.Perfume,
	all []models.Perfume,
	mathesCount int,
	threadsCount int,
) *MatchData {
	return &MatchData{
		Matcher:      matcher,
		favourite:    favourite,
		all:          all,
		matchesCount: mathesCount,
		threadsCount: threadsCount,
	}
}

func Find(md *MatchData) []models.Ranked {
	wg := sync.WaitGroup{}
	actualThreads := min(md.threadsCount, len(md.all))
	if actualThreads == 0 {
		return []models.Ranked{}
	}
	wg.Add(actualThreads)
	results := make(chan PerfumeHeap, actualThreads)

	chunkSize := len(md.all) / actualThreads
	for i := range actualThreads {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == actualThreads-1 {
			end = len(md.all)
		}
		go findChunk(md, md.all[start:end], results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return getMatchingResults(md, mergeHeaps(results, md.matchesCount))
}

func findChunk(md *MatchData, perfumes []models.Perfume, results chan<- PerfumeHeap, wg *sync.WaitGroup) {
	defer wg.Done()
	h := PerfumeHeap{
		limit: md.matchesCount,
	}
	heap.Init(&h)

	for _, p := range perfumes {
		if md.favourite.Equal(p) {
			continue
		}
		similarityScore := md.Matcher.GetPerfumeSimilarityScore(
			md.favourite.Properties,
			p.Properties,
		)

		h.Push(models.Ranked{
			Perfume: p,
			Score:   similarityScore,
		})
	}

	results <- h
}

func mergeHeaps(heaps <-chan PerfumeHeap, limit int) *PerfumeHeap {
	finalHeap := PerfumeHeap{
		limit: limit,
	}
	heap.Init(&finalHeap)

	for h := range heaps {
		for _, p := range h.perfumes {
			finalHeap.Push(p)
		}
	}
	return &finalHeap
}

func getMatchingResults(md *MatchData, h *PerfumeHeap) []models.Ranked {
	md.matchesCount = min(md.matchesCount, h.Len())
	ranked := make([]models.Ranked, md.matchesCount)
	for i := md.matchesCount - 1; i >= 0; i-- {
		ranked[i] = heap.Pop(h).(models.Ranked)
		ranked[i].Rank = i + 1
	}
	return ranked
}
