package advising

import (
	"context"
	"testing"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/config"
	"github.com/zemld/Scently/perfumist/internal/models/matching"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
)

func TestNewTagsBased(t *testing.T) {
	t.Parallel()

	weights := matching.NewBaseWeights(0.3, 0.4, 0.3)
	matcher := matching.NewTagsBasedAdapter(*weights, []string{"floral", "sweet"})
	fetcher := &MockFetcher{}
	mockConfig := &config.MockConfigManager{}

	tagsBasedAdvisor := NewTagsBased(matcher, fetcher, mockConfig)

	if tagsBasedAdvisor == nil {
		t.Fatal("expected non-nil TagsBased")
	}
	if tagsBasedAdvisor.fetcher != fetcher {
		t.Fatal("expected fetcher to be set")
	}
	if tagsBasedAdvisor.cm != mockConfig {
		t.Fatal("expected config manager to be set")
	}
}

func TestTagsBased_Advise_Success(t *testing.T) {
	t.Parallel()

	perfumes := []models.Perfume{
		{
			Brand: "Dior",
			Name:  "J'adore",
			Sex:   "female",
			Properties: models.Properties{
				EnrichedUpperNotes: []models.EnrichedNote{
					{Name: "Ylang-Ylang", Tags: []string{"floral", "sweet"}},
				},
				EnrichedCoreNotes: []models.EnrichedNote{
					{Name: "Rose", Tags: []string{"floral", "romantic"}},
					{Name: "Jasmine", Tags: []string{"floral", "sweet"}},
				},
				EnrichedBaseNotes: []models.EnrichedNote{
					{Name: "Vanilla", Tags: []string{"sweet", "warm"}},
				},
			},
		},
		{
			Brand: "Tom Ford",
			Name:  "Black Orchid",
			Sex:   "unisex",
			Properties: models.Properties{
				EnrichedUpperNotes: []models.EnrichedNote{
					{Name: "Truffle", Tags: []string{"earthy", "spicy"}},
				},
				EnrichedCoreNotes: []models.EnrichedNote{
					{Name: "Fruit", Tags: []string{"fruity", "fresh"}},
				},
				EnrichedBaseNotes: []models.EnrichedNote{
					{Name: "Patchouli", Tags: []string{"woody", "earthy"}},
				},
			},
		},
	}

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, param parameters.RequestPerfume) <-chan models.Perfume {
			ch := make(chan models.Perfume)
			go func() {
				defer close(ch)
				if param.Sex == "female" && param.Brand == "" && param.Name == "" {
					for _, p := range perfumes {
						select {
						case <-ctx.Done():
							return
						case ch <- p:
						}
					}
				}
			}()
			return ch
		},
	}

	weights := matching.NewBaseWeights(0.3, 0.4, 0.3)
	matcher := matching.NewTagsBasedAdapter(*weights, []string{"floral", "sweet"})

	mockConfig := &config.MockConfigManager{
		GetIntWithDefaultFunc: func(key string, defaultValue int) int {
			switch key {
			case "suggest_count":
				return 4
			case "threads_count":
				return 8
			case "minimal_tag_count":
				return 3
			}
			return defaultValue
		},
	}

	advisor := NewTagsBased(matcher, fetcher, mockConfig)
	params := parameters.RequestPerfume{
		Sex: "female",
	}

	result, err := advisor.Advise(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result) == 0 {
		t.Fatal("expected at least one match")
	}
}

func TestTagsBased_Advise_FetcherFails(t *testing.T) {
	t.Parallel()

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, param parameters.RequestPerfume) <-chan models.Perfume {
			ch := make(chan models.Perfume)
			close(ch)
			return ch
		},
	}

	weights := matching.NewBaseWeights(0.3, 0.4, 0.3)
	matcher := matching.NewTagsBasedAdapter(*weights, []string{})
	mockConfig := &config.MockConfigManager{}

	advisor := NewTagsBased(matcher, fetcher, mockConfig)
	params := parameters.RequestPerfume{
		Sex: "female",
	}

	result, err := advisor.Advise(context.Background(), params)

	// When fetcher returns empty channel, Common.Advise returns empty array without error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestTagsBased_Advise_FetcherReturnsEmpty(t *testing.T) {
	t.Parallel()

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, param parameters.RequestPerfume) <-chan models.Perfume {
			ch := make(chan models.Perfume)
			close(ch)
			return ch
		},
	}

	weights := matching.NewBaseWeights(0.3, 0.4, 0.3)
	matcher := matching.NewTagsBasedAdapter(*weights, []string{})
	mockConfig := &config.MockConfigManager{}

	advisor := NewTagsBased(matcher, fetcher, mockConfig)
	params := parameters.RequestPerfume{
		Sex: "female",
	}

	result, err := advisor.Advise(context.Background(), params)

	// When fetcher returns empty channel, Common.Advise returns empty array without error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestTagsBased_Advise_RespectsConfigParams(t *testing.T) {
	t.Parallel()

	perfumes := make([]models.Perfume, 10)
	for i := range perfumes {
		perfumes[i] = models.Perfume{
			Brand: "Brand",
			Name:  "Perfume",
			Sex:   "female",
			Properties: models.Properties{
				EnrichedUpperNotes: []models.EnrichedNote{
					{Name: "Rose", Tags: []string{"floral"}},
				},
			},
		}
	}

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, param parameters.RequestPerfume) <-chan models.Perfume {
			ch := make(chan models.Perfume)
			go func() {
				defer close(ch)
				if param.Sex == "female" {
					for _, p := range perfumes {
						select {
						case <-ctx.Done():
							return
						case ch <- p:
						}
					}
				}
			}()
			return ch
		},
	}

	weights := matching.NewBaseWeights(0.3, 0.4, 0.3)
	matcher := matching.NewTagsBasedAdapter(*weights, []string{"floral", "sweet"})

	var suggestCount int
	var threadsCount int

	mockConfig := &config.MockConfigManager{
		GetIntWithDefaultFunc: func(key string, defaultValue int) int {
			switch key {
			case "suggest_count":
				suggestCount = 3
				return 3
			case "threads_count":
				threadsCount = 2
				return 2
			}
			return defaultValue
		},
	}

	advisor := NewTagsBased(matcher, fetcher, mockConfig)
	params := parameters.RequestPerfume{
		Sex: "female",
	}

	result, err := advisor.Advise(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if suggestCount != 3 {
		t.Fatalf("expected suggest_count to be called, got %d", suggestCount)
	}
	if threadsCount != 2 {
		t.Fatalf("expected threads_count to be called, got %d", threadsCount)
	}
	if len(result) > 3 {
		t.Fatalf("expected at most 3 results, got %d", len(result))
	}
}

func TestTagsBased_Advise_CalculatesPerfumeTags(t *testing.T) {
	t.Parallel()

	perfumes := []models.Perfume{
		{
			Brand: "Dior",
			Name:  "J'adore",
			Sex:   "female",
			Properties: models.Properties{
				EnrichedUpperNotes: []models.EnrichedNote{
					{Name: "Rose", Tags: []string{"floral"}},
					{Name: "Jasmine", Tags: []string{"floral", "sweet"}},
				},
				EnrichedCoreNotes: []models.EnrichedNote{
					{Name: "Vanilla", Tags: []string{"sweet", "warm"}},
					{Name: "Amber", Tags: []string{"warm"}},
				},
				EnrichedBaseNotes: []models.EnrichedNote{
					{Name: "Musk", Tags: []string{"woody"}},
				},
			},
		},
	}

	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, param parameters.RequestPerfume) <-chan models.Perfume {
			ch := make(chan models.Perfume)
			go func() {
				defer close(ch)
				if param.Sex == "female" {
					for _, p := range perfumes {
						select {
						case <-ctx.Done():
							return
						case ch <- p:
						}
					}
				}
			}()
			return ch
		},
	}

	weights := matching.NewBaseWeights(0.3, 0.4, 0.3)
	matcher := matching.NewTagsBasedAdapter(*weights, []string{"floral", "sweet"})

	mockConfig := &config.MockConfigManager{
		GetIntWithDefaultFunc: func(key string, defaultValue int) int {
			switch key {
			case "suggest_count":
				return 4
			case "threads_count":
				return 8
			case "minimal_tag_count":
				return 1 // floral appears 2 times (> 1), should be included
			}
			return defaultValue
		},
	}

	advisor := NewTagsBased(matcher, fetcher, mockConfig)
	params := parameters.RequestPerfume{
		Sex: "female",
	}

	result, err := advisor.Advise(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected at least one result")
	}
	// Verify that tags were calculated
	found := false
	for _, note := range result[0].Perfume.Properties.EnrichedUpperNotes {
		if note.Name == "Rose" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'floral' tag to be calculated (appears 2 times)")
	}
}

func TestTagsBased_Advise_VerifyFetchParams(t *testing.T) {
	t.Parallel()

	perfumes := []models.Perfume{
		{
			Brand: "Dior",
			Name:  "J'adore",
			Sex:   "female",
		},
	}

	var fetchedSex models.Sex
	fetcher := &MockFetcher{
		FetchFunc: func(ctx context.Context, param parameters.RequestPerfume) <-chan models.Perfume {
			ch := make(chan models.Perfume)
			go func() {
				defer close(ch)
				fetchedSex = param.Sex
				if param.Brand == "" && param.Name == "" {
					for _, p := range perfumes {
						select {
						case <-ctx.Done():
							return
						case ch <- p:
						}
					}
				}
			}()
			return ch
		},
	}

	weights := matching.NewBaseWeights(0.3, 0.4, 0.3)
	matcher := matching.NewTagsBasedAdapter(*weights, []string{"floral", "sweet"})

	mockConfig := &config.MockConfigManager{
		GetIntWithDefaultFunc: func(key string, defaultValue int) int {
			return defaultValue
		},
	}

	advisor := NewTagsBased(matcher, fetcher, mockConfig)
	params := parameters.RequestPerfume{
		Sex: "female",
	}

	_, err := advisor.Advise(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedSex != "female" {
		t.Fatalf("expected sex 'female' to be fetched, got %q", fetchedSex)
	}
}
