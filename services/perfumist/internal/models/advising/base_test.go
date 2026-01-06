package advising

import (
	"context"
	"testing"

	"github.com/zemld/Scently/perfumist/internal/config"
	"github.com/zemld/Scently/perfumist/internal/errors"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
	"github.com/zemld/Scently/models"
)

type MockFetcher struct {
	FetchFunc func(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool)
}

func (m *MockFetcher) Fetch(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
	if m.FetchFunc != nil {
		return m.FetchFunc(ctx, params)
	}
	return nil, false
}

type MockMatcher struct {
	GetPerfumeSimilarityScoreFunc func(first models.Properties, second models.Properties) float64
}

func (m *MockMatcher) GetPerfumeSimilarityScore(first models.Properties, second models.Properties) float64 {
	if m.GetPerfumeSimilarityScoreFunc != nil {
		return m.GetPerfumeSimilarityScoreFunc(first, second)
	}
	return 0.0
}

func TestBase_Advise_Success(t *testing.T) {
	t.Parallel()

	favouritePerfume := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
		Properties: models.Properties{
			Type:       "Eau de Parfum",
			Family:     []string{"Floral"},
			UpperNotes: []string{"Aldehydes", "Ylang-Ylang"},
			CoreNotes:  []string{"Rose", "Jasmine"},
			BaseNotes:  []string{"Vanilla", "Amber"},
		},
	}

	allPerfumes := []models.Perfume{
		{
			Brand: "Dior",
			Name:  "J'adore",
			Sex:   "female",
			Properties: models.Properties{
				Type:       "Eau de Parfum",
				Family:     []string{"Floral"},
				UpperNotes: []string{"Ylang-Ylang"},
				CoreNotes:  []string{"Rose", "Jasmine"},
				BaseNotes:  []string{"Vanilla"},
			},
		},
		{
			Brand: "Tom Ford",
			Name:  "Black Orchid",
			Sex:   "unisex",
			Properties: models.Properties{
				Type:       "Eau de Parfum",
				Family:     []string{"Oriental"},
				UpperNotes: []string{"Truffle"},
				CoreNotes:  []string{"Fruit"},
				BaseNotes:  []string{"Patchouli"},
			},
		},
	}

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
			if len(params) == 1 {
				// First fetch - favourite perfume (has Brand and Name)
				if params[0].Brand == "Chanel" && params[0].Name == "No5" {
					return []models.Perfume{favouritePerfume}, true
				}
				// Second fetch - all perfumes with same sex (only Sex is set, Brand and Name are empty)
				if params[0].Brand == "" && params[0].Name == "" && params[0].Sex == "female" {
					return allPerfumes, true
				}
			}
			return nil, false
		},
	}

	matcher := &MockMatcher{
		GetPerfumeSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			// Simple similarity: count matching notes
			score := 0.0
			if first.Type == second.Type {
				score += 0.3
			}
			// Check for matching notes
			firstNotes := make(map[string]bool)
			for _, n := range first.UpperNotes {
				firstNotes[n] = true
			}
			for _, n := range first.CoreNotes {
				firstNotes[n] = true
			}
			for _, n := range first.BaseNotes {
				firstNotes[n] = true
			}
			for _, n := range second.UpperNotes {
				if firstNotes[n] {
					score += 0.1
				}
			}
			for _, n := range second.CoreNotes {
				if firstNotes[n] {
					score += 0.1
				}
			}
			for _, n := range second.BaseNotes {
				if firstNotes[n] {
					score += 0.1
				}
			}
			return score
		},
	}

	mockConfig := &config.MockConfigManager{}
	base := NewBase(fetcher, matcher, mockConfig)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := base.Advise(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// Should return matches (excluding the favourite perfume itself)
	if len(result) == 0 {
		t.Fatal("expected at least one match")
	}
	// Verify that favourite perfume is not in results
	for _, r := range result {
		if r.Perfume.Equal(favouritePerfume) {
			t.Fatal("favourite perfume should not be in results")
		}
	}
}

func TestBase_Advise_FetcherFailsOnFirstFetch(t *testing.T) {
	t.Parallel()

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
			return nil, false
		},
	}

	matcher := &MockMatcher{}
	mockConfig := &config.MockConfigManager{}
	base := NewBase(fetcher, matcher, mockConfig)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := base.Advise(context.Background(), params)

	if err == nil {
		t.Fatal("expected error when fetcher fails on first fetch")
	}
	serviceErr, ok := err.(*errors.ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if serviceErr.Message != "failed to interact with perfume service" {
		t.Fatalf("expected error message 'failed to interact with perfume service', got %q", serviceErr.Message)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestBase_Advise_FetcherReturnsEmptyOnFirstFetch(t *testing.T) {
	t.Parallel()

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
			return []models.Perfume{}, true
		},
	}

	matcher := &MockMatcher{}
	mockConfig := &config.MockConfigManager{}
	base := NewBase(fetcher, matcher, mockConfig)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := base.Advise(context.Background(), params)

	if err == nil {
		t.Fatal("expected error when fetcher returns empty on first fetch")
	}
	notFoundErr, ok := err.(*errors.NotFoundError)
	if !ok {
		t.Fatalf("expected NotFoundError, got %T", err)
	}
	if notFoundErr.Message != "perfume not found" {
		t.Fatalf("expected error message 'perfume not found', got %q", notFoundErr.Message)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestBase_Advise_FetcherFailsOnSecondFetch(t *testing.T) {
	t.Parallel()

	favouritePerfume := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" && params[0].Name == "No5" {
				// First fetch succeeds
				return []models.Perfume{favouritePerfume}, true
			}
			// Second fetch fails
			return nil, false
		},
	}

	matcher := &MockMatcher{}
	mockConfig := &config.MockConfigManager{}
	base := NewBase(fetcher, matcher, mockConfig)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := base.Advise(context.Background(), params)

	if err == nil {
		t.Fatal("expected error when fetcher fails on second fetch")
	}
	serviceErr, ok := err.(*errors.ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if serviceErr.Message != "failed to interact with perfume service" {
		t.Fatalf("expected error message 'failed to interact with perfume service', got %q", serviceErr.Message)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestBase_Advise_FetcherReturnsEmptyOnSecondFetch(t *testing.T) {
	t.Parallel()

	favouritePerfume := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" && params[0].Name == "No5" {
				// First fetch succeeds
				return []models.Perfume{favouritePerfume}, true
			}
			// Second fetch returns empty
			return []models.Perfume{}, true
		},
	}

	matcher := &MockMatcher{}
	mockConfig := &config.MockConfigManager{}
	base := NewBase(fetcher, matcher, mockConfig)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := base.Advise(context.Background(), params)

	if err == nil {
		t.Fatal("expected error when fetcher returns empty on second fetch")
	}
	serviceErr, ok := err.(*errors.ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if serviceErr.Message != "no perfumes available in database" {
		t.Fatalf("expected error message 'no perfumes available in database', got %q", serviceErr.Message)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestBase_Advise_VerifySecondFetchParams(t *testing.T) {
	t.Parallel()

	favouritePerfume := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	allPerfumes := []models.Perfume{
		{
			Brand: "Dior",
			Name:  "J'adore",
			Sex:   "female",
		},
	}

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" && params[0].Name == "No5" {
				// First fetch - favourite perfume
				return []models.Perfume{favouritePerfume}, true
			}
			// Second fetch - verify params
			if len(params) != 1 {
				t.Fatalf("expected 1 param for second fetch, got %d", len(params))
			}
			if params[0].Sex != "female" {
				t.Fatalf("expected sex 'female' in second fetch params, got %q", params[0].Sex)
			}
			if params[0].Brand != "" {
				t.Fatalf("expected empty brand in second fetch params, got %q", params[0].Brand)
			}
			if params[0].Name != "" {
				t.Fatalf("expected empty name in second fetch params, got %q", params[0].Name)
			}
			return allPerfumes, true
		},
	}

	matcher := &MockMatcher{
		GetPerfumeSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	mockConfig := &config.MockConfigManager{}
	base := NewBase(fetcher, matcher, mockConfig)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := base.Advise(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestBase_Advise_RespectsAdviseCount(t *testing.T) {
	t.Parallel()

	favouritePerfume := models.Perfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	// Create more perfumes than adviseCount
	allPerfumes := make([]models.Perfume, 10)
	for i := range allPerfumes {
		allPerfumes[i] = models.Perfume{
			Brand: "Brand",
			Name:  "Perfume",
			Sex:   "female",
		}
	}

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" && params[0].Name == "No5" {
				return []models.Perfume{favouritePerfume}, true
			}
			return allPerfumes, true
		},
	}

	matcher := &MockMatcher{
		GetPerfumeSimilarityScoreFunc: func(first models.Properties, second models.Properties) float64 {
			return 0.5
		},
	}

	mockConfig := &config.MockConfigManager{
		GetIntWithDefaultFunc: func(key string, defaultValue int) int {
			if key == "suggest_count" {
				return 3
			}
			return defaultValue
		},
	}
	base := NewBase(fetcher, matcher, mockConfig)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := base.Advise(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expectedCount := 3
	if len(result) > expectedCount {
		t.Fatalf("expected at most %d results, got %d", expectedCount, len(result))
	}
}
