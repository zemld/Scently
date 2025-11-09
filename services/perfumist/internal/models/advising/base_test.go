package advising

import (
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

// MockFetcher is a mock implementation of fetching.Fetcher
type MockFetcher struct {
	FetchFunc func([]parameters.RequestPerfume) ([]perfume.Perfume, bool)
}

func (m *MockFetcher) Fetch(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
	if m.FetchFunc != nil {
		return m.FetchFunc(params)
	}
	return nil, false
}

// MockMatcher is a mock implementation of matching.Matcher
type MockMatcher struct {
	FindFunc func(favourite perfume.Perfume, all []perfume.Perfume, matchesCount int) []perfume.Ranked
}

func (m *MockMatcher) Find(favourite perfume.Perfume, all []perfume.Perfume, matchesCount int) []perfume.Ranked {
	if m.FindFunc != nil {
		return m.FindFunc(favourite, all, matchesCount)
	}
	return nil
}

func TestNewBaseAdvisor(t *testing.T) {
	t.Parallel()

	fetcher := &MockFetcher{}
	matcher := &MockMatcher{}
	adviseCount := 5

	advisor := NewBaseAdvisor(fetcher, matcher, adviseCount)

	if advisor == nil {
		t.Fatal("expected non-nil advisor")
	}
	if advisor.fetcher != fetcher {
		t.Fatal("expected fetcher to be set")
	}
	if advisor.matcher != matcher {
		t.Fatal("expected matcher to be set")
	}
	if advisor.adviseCount != adviseCount {
		t.Fatalf("expected adviseCount %d, got %d", adviseCount, advisor.adviseCount)
	}
}

func TestBaseAdvisor_Advise_Success(t *testing.T) {
	t.Parallel()

	favouritePerfume := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	allPerfumes := []perfume.Perfume{
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}

	expectedRanked := []perfume.Ranked{
		{Perfume: allPerfumes[0], Rank: 1, Score: 0.9},
		{Perfume: allPerfumes[1], Rank: 2, Score: 0.7},
	}

	fetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" && params[0].Name == "No5" {
				// First call: fetch favourite
				return []perfume.Perfume{favouritePerfume}, true
			}
			// Second call: fetch all perfumes by sex
			return allPerfumes, true
		},
	}

	matcher := &MockMatcher{
		FindFunc: func(fav perfume.Perfume, all []perfume.Perfume, count int) []perfume.Ranked {
			if !fav.Equal(favouritePerfume) {
				t.Fatalf("expected favourite perfume %+v, got %+v", favouritePerfume, fav)
			}
			if len(all) != len(allPerfumes) {
				t.Fatalf("expected %d perfumes, got %d", len(allPerfumes), len(all))
			}
			if count != 5 {
				t.Fatalf("expected count 5, got %d", count)
			}
			return expectedRanked
		},
	}

	advisor := NewBaseAdvisor(fetcher, matcher, 5)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != len(expectedRanked) {
		t.Fatalf("expected %d results, got %d", len(expectedRanked), len(result))
	}
	for i, r := range result {
		if r.Rank != expectedRanked[i].Rank {
			t.Fatalf("result[%d]: expected rank %d, got %d", i, expectedRanked[i].Rank, r.Rank)
		}
		if r.Score != expectedRanked[i].Score {
			t.Fatalf("result[%d]: expected score %f, got %f", i, expectedRanked[i].Score, r.Score)
		}
	}
}

func TestBaseAdvisor_Advise_FetcherFailsOnFavourite(t *testing.T) {
	t.Parallel()

	fetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return nil, false
		},
	}

	matcher := &MockMatcher{}
	advisor := NewBaseAdvisor(fetcher, matcher, 5)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err == nil {
		t.Fatal("expected error when fetcher fails")
	}
	if err.Error() != "failed to get favourite perfumes" {
		t.Fatalf("expected error 'failed to get favourite perfumes', got %q", err.Error())
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestBaseAdvisor_Advise_FetcherReturnsEmptyFavourite(t *testing.T) {
	t.Parallel()

	fetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" {
				return []perfume.Perfume{}, true
			}
			return []perfume.Perfume{{Brand: "Dior"}}, true
		},
	}

	matcher := &MockMatcher{}
	advisor := NewBaseAdvisor(fetcher, matcher, 5)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err == nil {
		t.Fatal("expected error when favourite is empty")
	}
	if err.Error() != "failed to get favourite perfumes" {
		t.Fatalf("expected error 'failed to get favourite perfumes', got %q", err.Error())
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestBaseAdvisor_Advise_FetcherFailsOnAllPerfumes(t *testing.T) {
	t.Parallel()

	favouritePerfume := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	fetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" && params[0].Name == "No5" {
				// First call: fetch favourite - succeeds
				return []perfume.Perfume{favouritePerfume}, true
			}
			// Second call: fetch all perfumes - fails
			return nil, false
		},
	}

	matcher := &MockMatcher{}
	advisor := NewBaseAdvisor(fetcher, matcher, 5)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err == nil {
		t.Fatal("expected error when fetcher fails on all perfumes")
	}
	if err.Error() != "failed to get all perfumes" {
		t.Fatalf("expected error 'failed to get all perfumes', got %q", err.Error())
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestBaseAdvisor_Advise_FetcherReturnsEmptyAllPerfumes(t *testing.T) {
	t.Parallel()

	favouritePerfume := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	fetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" && params[0].Name == "No5" {
				// First call: fetch favourite - succeeds
				return []perfume.Perfume{favouritePerfume}, true
			}
			// Second call: fetch all perfumes - returns empty
			return []perfume.Perfume{}, true
		},
	}

	matcher := &MockMatcher{}
	advisor := NewBaseAdvisor(fetcher, matcher, 5)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err == nil {
		t.Fatal("expected error when all perfumes is empty")
	}
	if err.Error() != "failed to get all perfumes" {
		t.Fatalf("expected error 'failed to get all perfumes', got %q", err.Error())
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestBaseAdvisor_Advise_WithSexFilter(t *testing.T) {
	t.Parallel()

	favouritePerfume := perfume.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	allPerfumes := []perfume.Perfume{
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}

	fetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" {
				return []perfume.Perfume{favouritePerfume}, true
			}
			// Verify that second call uses sex filter
			if len(params) == 1 && params[0].Sex == "female" && params[0].Brand == "" {
				return allPerfumes, true
			}
			return nil, false
		},
	}

	matcher := &MockMatcher{
		FindFunc: func(fav perfume.Perfume, all []perfume.Perfume, count int) []perfume.Ranked {
			return []perfume.Ranked{
				{Perfume: all[0], Rank: 1, Score: 0.9},
			}
		},
	}

	advisor := NewBaseAdvisor(fetcher, matcher, 5)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}
