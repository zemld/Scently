package matching

import (
	"container/heap"
	"sync"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
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

func Find(md MatchData) []models.Ranked {
	h, wg := initThreadingTools(md.threadsCount)
	findThreaded(md, h, wg)
	wg.Wait()
	return getMatchingResults(md, h)
}

func initThreadingTools(threadsCount int) (*PerfumeHeap, *sync.WaitGroup) {
	matchingHeap := PerfumeHeap{}
	heap.Init(&matchingHeap)

	wg := sync.WaitGroup{}
	wg.Add(threadsCount)

	return &matchingHeap, &wg
}

func findThreaded(md MatchData, h *PerfumeHeap, wg *sync.WaitGroup) {
	chunkSize := len(md.all) / md.threadsCount
	for i := range md.threadsCount {
		start := i * chunkSize
		end := start + chunkSize
		if i == md.threadsCount-1 {
			end = len(md.all)
		}
		go processPerfumes(
			h,
			wg,
			md,
			md.all[start:end],
		)
	}
}

func processPerfumes(
	h *PerfumeHeap,
	wg *sync.WaitGroup,
	md MatchData,
	all []models.Perfume,
) {
	defer wg.Done()

	for _, p := range all {
		if md.favourite.Equal(p) {
			continue
		}
		similarityScore := md.Matcher.GetPerfumeSimilarityScore(
			md.favourite.Properties,
			p.Properties,
		)

		h.PushSafeIfNeeded(models.WithScore{
			Perfume: p,
			Score:   similarityScore,
		}, md.matchesCount)
	}
}

func getMatchingResults(md MatchData, h *PerfumeHeap) []models.Ranked {
	md.matchesCount = min(md.matchesCount, h.Len())

	ranked := make([]models.Ranked, md.matchesCount)
	for i := md.matchesCount - 1; i >= 0; i-- {
		mostSimilar := h.PopSafe()
		ranked[i] = models.Ranked{
			Perfume: mostSimilar.Perfume,
			Rank:    i + 1,
			Score:   mostSimilar.Score,
		}
	}
	return ranked
}

func multiplyMaps(first map[string]float64, second map[string]float64) float64 {
	if len(second) < len(first) {
		return multiplyMaps(second, first)
	}

	score := 0.0
	for tag, value := range first {
		if secondValue, ok := second[tag]; ok {
			score += value * secondValue
		}
	}
	return score
}
