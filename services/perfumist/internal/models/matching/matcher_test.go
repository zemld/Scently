package matching

import (
	"testing"

	"github.com/zemld/Scently/models"
)

type MockMatcher struct {
	GetSimilarityScoreFunc func(first models.Properties, second models.Properties) float64
}

func (m *MockMatcher) GetSimilarityScore(first models.Properties, second models.Properties) float64 {
	if m.GetSimilarityScoreFunc != nil {
		return m.GetSimilarityScoreFunc(first, second)
	}
	return 0.0
}

func TestNewMatchData(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{}
	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}
	all := []models.Perfume{
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}
	matchesCount := 5
	threadsCount := 2

	md := NewMatchData(matcher, favourite, all, matchesCount, threadsCount)

	if md == nil {
		t.Fatal("expected non-nil MatchData")
	}
	if md.Matcher != matcher {
		t.Fatal("expected matcher to be set")
	}
	if !md.favourite.Equal(favourite) {
		t.Fatal("expected favourite to be set")
	}
	if len(md.all) != len(all) {
		t.Fatalf("expected all length %d, got %d", len(all), len(md.all))
	}
	if md.matchesCount != matchesCount {
		t.Fatalf("expected matchesCount %d, got %d", matchesCount, md.matchesCount)
	}
	if md.threadsCount != threadsCount {
		t.Fatalf("expected threadsCount %d, got %d", threadsCount, md.threadsCount)
	}
}

func TestFind_EmptyAll(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{}
	favourite := models.Perfume{Brand: "Chanel", Name: "No5", Sex: "female"}
	all := []models.Perfume{}

	md := NewMatchData(matcher, favourite, all, 5, 2)
	result := Find(md)

	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d items", len(result))
	}
}

func TestFind_ZeroThreads(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{}
	favourite := models.Perfume{Brand: "Chanel", Name: "No5", Sex: "female"}
	all := []models.Perfume{
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
	}

	md := NewMatchData(matcher, favourite, all, 5, 0)
	result := Find(md)

	if len(result) != 0 {
		t.Fatalf("expected empty result with 0 threads, got %d items", len(result))
	}
}

func TestFind_ExcludesFavourite(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	all := []models.Perfume{
		favourite, // Same as favourite
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}

	md := NewMatchData(matcher, favourite, all, 5, 2)
	result := Find(md)

	// Should exclude favourite perfume
	for _, r := range result {
		if r.Perfume.Equal(favourite) {
			t.Fatal("favourite perfume should not be in results")
		}
	}

	// Should have 2 results (excluding favourite)
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestFind_RespectsMatchesCount(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	// Create more perfumes than matchesCount
	all := make([]models.Perfume, 10)
	for i := range all {
		all[i] = models.Perfume{
			Brand: "Brand",
			Name:  "Perfume",
			Sex:   "female",
		}
	}

	matchesCount := 3
	md := NewMatchData(matcher, favourite, all, matchesCount, 2)
	result := Find(md)

	if len(result) > matchesCount {
		t.Fatalf("expected at most %d results, got %d", matchesCount, len(result))
	}
}

func TestFind_AssignsRanks(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	all := []models.Perfume{
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
		{Brand: "Yves Saint Laurent", Name: "Opium", Sex: "unisex"},
	}

	md := NewMatchData(matcher, favourite, all, 5, 2)
	result := Find(md)

	expectedRanks := make(map[int]bool)
	for i := 0; i < len(result); i++ {
		expectedRanks[i+1] = true
	}

	for i, r := range result {
		if r.Rank < 1 || r.Rank > len(result) {
			t.Fatalf("result[%d]: expected rank in range [1, %d], got %d", i, len(result), r.Rank)
		}
		if !expectedRanks[r.Rank] {
			t.Fatalf("result[%d]: got duplicate or invalid rank %d", i, r.Rank)
		}
		delete(expectedRanks, r.Rank)
	}

	if len(expectedRanks) > 0 {
		t.Fatalf("not all ranks were assigned: missing %v", expectedRanks)
	}
}

func TestFind_SortsByScore(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	all := []models.Perfume{
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
		{Brand: "Yves Saint Laurent", Name: "Opium", Sex: "unisex"},
	}

	md := NewMatchData(matcher, favourite, all, 5, 1) // Use single thread for deterministic ordering
	result := Find(md)

	// Verify that results have scores and ranks assigned
	if len(result) < 2 {
		t.Fatal("expected at least 2 results")
	}

	// Verify all results have valid scores
	for i, r := range result {
		if r.Score < 0 || r.Score > 1.0 {
			t.Fatalf("result[%d]: expected score in range [0, 1], got %f", i, r.Score)
		}
		if r.Rank < 1 || r.Rank > len(result) {
			t.Fatalf("result[%d]: expected rank in range [1, %d], got %d", i, len(result), r.Rank)
		}
		// Verify that each result has a perfume
		if r.Perfume.Brand == "" {
			t.Fatalf("result[%d]: expected perfume to have brand", i)
		}
	}
}

func TestFind_SingleThread(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	all := []models.Perfume{
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}

	md := NewMatchData(matcher, favourite, all, 5, 1)
	result := Find(md)

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestFind_MoreThreadsThanItems(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	all := []models.Perfume{
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
	}

	// More threads than items
	md := NewMatchData(matcher, favourite, all, 5, 10)
	result := Find(md)

	// Should use min(threads, len(all)) = 1 thread
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}

func TestFind_LargeDataset(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	// Create a larger dataset
	all := make([]models.Perfume, 100)
	for i := range all {
		all[i] = models.Perfume{
			Brand: "Brand",
			Name:  "Perfume",
			Sex:   "female",
		}
	}

	md := NewMatchData(matcher, favourite, all, 10, 4)
	result := Find(md)

	if len(result) > 10 {
		t.Fatalf("expected at most 10 results, got %d", len(result))
	}

	expectedRanks := make(map[int]bool)
	for i := 1; i <= len(result); i++ {
		expectedRanks[i] = true
	}

	for i, r := range result {
		if r.Rank < 1 || r.Rank > len(result) {
			t.Fatalf("result[%d]: expected rank in range [1, %d], got %d", i, len(result), r.Rank)
		}
		if !expectedRanks[r.Rank] {
			t.Fatalf("result[%d]: got duplicate or invalid rank %d", i, r.Rank)
		}
		delete(expectedRanks, r.Rank)
		if r.Score != 0.5 {
			t.Fatalf("result[%d]: expected score 0.5, got %f", i, r.Score)
		}
	}

	if len(expectedRanks) > 0 {
		t.Fatalf("not all ranks were assigned: missing %v", expectedRanks)
	}
}

func TestFind_ChunkDistribution(t *testing.T) {
	t.Parallel()

	callCount := 0
	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			callCount++
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	// Create dataset that should be split into chunks
	all := make([]models.Perfume, 10)
	for i := range all {
		all[i] = models.Perfume{
			Brand: "Brand",
			Name:  "Perfume",
			Sex:   "female",
		}
	}

	md := NewMatchData(matcher, favourite, all, 10, 3)
	result := Find(md)

	// Should process all items (excluding favourite if it's in the list)
	// With 3 threads and 10 items, each thread should process ~3-4 items
	// Total calls should be 10 (one per item in all)
	if callCount != 10 {
		t.Fatalf("expected 10 similarity score calls, got %d", callCount)
	}

	if len(result) != 10 {
		t.Fatalf("expected 10 results, got %d", len(result))
	}
}

func TestFind_AllItemsAreFavourite(t *testing.T) {
	t.Parallel()

	matcher := &MockMatcher{
		GetSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	favourite := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	// All items are the same as favourite
	all := []models.Perfume{
		favourite,
		favourite,
		favourite,
	}

	md := NewMatchData(matcher, favourite, all, 5, 2)
	result := Find(md)

	// Should return empty since all items are excluded
	if len(result) != 0 {
		t.Fatalf("expected empty result when all items are favourite, got %d", len(result))
	}
}
