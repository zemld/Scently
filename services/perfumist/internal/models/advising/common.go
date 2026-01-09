package advising

import (
	"container/heap"
	"context"
	"log"
	"sync"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/models/fetching"
	"github.com/zemld/Scently/perfumist/internal/models/matching"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
	"github.com/zemld/config-manager/pkg/cm"
)

type Common struct {
	favouritePerfume models.Perfume

	matcher matching.Matcher
	fetcher fetching.Fetcher

	cm cm.ConfigManager

	matchesCount int
	workersCount int
}

func NewCommon(fetcher fetching.Fetcher, matcher matching.Matcher, cm cm.ConfigManager) *Common {
	return &Common{
		fetcher:      fetcher,
		matcher:      matcher,
		cm:           cm,
		matchesCount: cm.GetIntWithDefault("suggest_count", 4),
		workersCount: cm.GetIntWithDefault("threads_count", 8),
	}
}

func (a *Common) WithFavouritePerfume(favouritePerfume models.Perfume) *Common {
	a.favouritePerfume = favouritePerfume
	return a
}

func (a *Common) Advise(ctx context.Context, parameter parameters.RequestPerfume) ([]models.Ranked, error) {
	log.Printf("parameter: %+v\n", parameter)
	allPerfumesChan := a.fetcher.Fetch(ctx, *parameters.NewGet().WithSex(parameter.Sex))

	resultsChan := make(chan *matching.PerfumeHeap)

	wg := a.initWaitGroup()

	a.runAllAdvisingWorkers(
		ctx,
		wg,
		allPerfumesChan,
		resultsChan,
	)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	results := getMatchingResults(mergeHeaps(resultsChan, a.matchesCount), a.matchesCount)
	for i := range results {
		results[i].Perfume.Properties.Tags = matching.CalculatePerfumeTags(
			&results[i].Perfume.Properties,
			*matching.NewBaseWeights(
				a.cm.GetFloatWithDefault("upper_notes_weight", 0.2),
				a.cm.GetFloatWithDefault("core_notes_weight", 0.35),
				a.cm.GetFloatWithDefault("base_notes_weight", 0.45),
			),
		)
	}
	return results, nil
}

func (a *Common) initWaitGroup() *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(a.workersCount)
	return &wg
}

func (a *Common) runAllAdvisingWorkers(
	ctx context.Context,
	wg *sync.WaitGroup,
	perfumesChan <-chan models.Perfume,
	resultsChan chan<- *matching.PerfumeHeap,
) {
	for range a.workersCount {
		go func() {
			defer wg.Done()
			a.runAdvisingWorker(
				ctx,
				perfumesChan,
				resultsChan,
			)
		}()
	}
}

func (a *Common) runAdvisingWorker(
	ctx context.Context,
	jobs <-chan models.Perfume,
	results chan<- *matching.PerfumeHeap,
) {
	h := matching.NewPerfumeHeap(a.matchesCount)
	heap.Init(h)

	for {
		select {
		case <-ctx.Done():
			results <- h
			return
		case perfume, ok := <-jobs:
			if !ok {
				results <- h
				return
			}
			a.processPerfume(ctx, perfume, h)
		}
	}
}

func (a *Common) processPerfume(ctx context.Context, perfume models.Perfume, h *matching.PerfumeHeap) {
	if a.favouritePerfume.Equal(perfume) {
		return
	}
	matching.PreparePerfumeCharacteristics(&perfume)
	similarityScore := a.matcher.GetSimilarityScore(
		a.favouritePerfume.Properties,
		perfume.Properties,
	)

	h.Push(models.Ranked{
		Perfume: perfume,
		Score:   similarityScore,
	})
}

func mergeHeaps(heaps <-chan *matching.PerfumeHeap, matchesCount int) *matching.PerfumeHeap {
	suggestionsHeap := matching.NewPerfumeHeap(matchesCount)
	heap.Init(suggestionsHeap)
	for h := range heaps {
		for _, p := range h.GetPerfumes() {
			suggestionsHeap.Push(p)
		}
	}
	return suggestionsHeap
}

func getMatchingResults(h *matching.PerfumeHeap, matchesCount int) []models.Ranked {
	matchesCount = min(matchesCount, h.Len())
	ranked := make([]models.Ranked, matchesCount)
	for i := matchesCount - 1; i >= 0; i-- {
		ranked[i] = heap.Pop(h).(models.Ranked)
		ranked[i].Rank = i + 1
	}
	return ranked
}
