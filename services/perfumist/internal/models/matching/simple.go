package matching

import (
	"container/heap"
	"sync"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/pkg/set"
)

type Overlay struct {
	familyWeight float64
	notesWeight  float64
	typeWeight   float64

	upperNotesWeight  float64
	middleNotesWeight float64
	baseNotesWeight   float64

	threadsCount int
}

func NewSimpleMatcher(familyWeight float64,
	notesWeight float64,
	typeWeight float64,
	upperNotesWeight float64,
	middleNotesWeight float64,
	baseNotesWeight float64,
	threadsCount int,
) *Overlay {
	return &Overlay{
		familyWeight:      familyWeight,
		notesWeight:       notesWeight,
		typeWeight:        typeWeight,
		upperNotesWeight:  upperNotesWeight,
		middleNotesWeight: middleNotesWeight,
		baseNotesWeight:   baseNotesWeight,
		threadsCount:      threadsCount,
	}
}

func (m Overlay) Find(favourite perfume.Perfume, all []perfume.Perfume, matchesCount int) []perfume.Ranked {
	matchingHeap := PerfumeHeap{}
	heap.Init(&matchingHeap)

	wg := sync.WaitGroup{}
	wg.Add(m.threadsCount)

	for range m.threadsCount {
		go m.buildHeapAsync(&matchingHeap, &wg, favourite, all)
	}

	wg.Wait()

	availableCount := matchingHeap.Len()
	if matchesCount > availableCount {
		matchesCount = availableCount
	}

	if matchesCount == 0 {
		return []perfume.Ranked{}
	}

	mostSimilar := matchingHeap.perfumes[:matchesCount]
	ranked := make([]perfume.Ranked, 0, len(mostSimilar))
	for i := range mostSimilar {
		ranked = append(ranked, perfume.Ranked{
			Perfume: mostSimilar[i].Perfume,
			Rank:    i + 1,
			Score:   mostSimilar[i].Score,
		})
	}
	return ranked
}

func (m Overlay) buildHeapAsync(h *PerfumeHeap, wg *sync.WaitGroup, favourite perfume.Perfume, all []perfume.Perfume) {
	defer wg.Done()

	for _, p := range all {
		if favourite.Equal(p) {
			continue
		}
		similarityScore := m.GetPerfumeSimilarityScore(favourite.Properties, p.Properties)

		h.PushSafe(perfume.WithScore{
			Perfume: p,
			Score:   similarityScore,
		})
	}
}

func (m Overlay) GetPerfumeSimilarityScore(first perfume.Properties, second perfume.Properties) float64 {
	familiesSimilarityScore := m.getListSimilarityScore(first.Family, second.Family)
	notesSimilarityScore := m.getNotesSimilarityScore(first, second)
	typeSimilarity := m.getTypeSimilarityScore(first.Type, second.Type)
	return familiesSimilarityScore*m.familyWeight + notesSimilarityScore*m.notesWeight + typeSimilarity*m.typeWeight
}

func (m Overlay) getListSimilarityScore(first []string, second []string) float64 {
	firstSet := set.MakeSet(first)
	secondSet := set.MakeSet(second)
	intersection := set.Intersect(firstSet, secondSet)
	un := set.Union(firstSet, secondSet)

	if len(un) == 0 {
		return 0
	}
	return float64(len(intersection)) / float64(len(un))
}

func (m Overlay) getNotesSimilarityScore(first perfume.Properties, second perfume.Properties) float64 {
	upperNotesSimilarityScore := m.getListSimilarityScore(first.UpperNotes, second.UpperNotes)
	middleNotesSimilarityScore := m.getListSimilarityScore(first.CoreNotes, second.CoreNotes)
	baseNotesSimilarityScore := m.getListSimilarityScore(first.BaseNotes, second.BaseNotes)

	return upperNotesSimilarityScore*m.upperNotesWeight + middleNotesSimilarityScore*m.middleNotesWeight + baseNotesSimilarityScore*m.baseNotesWeight
}

func (m Overlay) getTypeSimilarityScore(first string, second string) float64 {
	if first == second {
		return 1
	}
	return 0
}
