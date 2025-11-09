package advising

import (
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

func TestNewAiAdvisor(t *testing.T) {
	t.Parallel()

	adviseFetcher := &MockFetcher{}
	enrichFetcher := &MockFetcher{}

	advisor := NewAI(adviseFetcher, enrichFetcher)

	if advisor == nil {
		t.Fatal("expected non-nil advisor")
	}
	if advisor.adviseFetcher != adviseFetcher {
		t.Fatal("expected adviseFetcher to be set")
	}
	if advisor.enrichFetcher != enrichFetcher {
		t.Fatal("expected enrichFetcher to be set")
	}
}

func TestAiAdvisor_Advise_Success(t *testing.T) {
	t.Parallel()

	aiSuggestions := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
	}

	enrichedPerfumes := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female", ImageUrl: "http://example.com/no5.jpg"},
		{Brand: "Dior", Name: "J'adore", Sex: "female", ImageUrl: "http://example.com/jadore.jpg"},
	}

	adviseFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			if len(params) == 1 && params[0].Brand == "Chanel" {
				return aiSuggestions, true
			}
			return nil, false
		},
	}

	enrichFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			// Verify enrichment params are constructed correctly
			if len(params) != len(aiSuggestions) {
				t.Fatalf("expected %d enrichment params, got %d", len(aiSuggestions), len(params))
			}
			for i, p := range params {
				if p.Brand != aiSuggestions[i].Brand {
					t.Fatalf("param[%d]: expected brand %q, got %q", i, aiSuggestions[i].Brand, p.Brand)
				}
				if p.Name != aiSuggestions[i].Name {
					t.Fatalf("param[%d]: expected name %q, got %q", i, aiSuggestions[i].Name, p.Name)
				}
				if p.Sex != "female" {
					t.Fatalf("param[%d]: expected sex %q, got %q", i, "female", p.Sex)
				}
			}
			return enrichedPerfumes, true
		},
	}

	advisor := NewAI(adviseFetcher, enrichFetcher)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != len(enrichedPerfumes) {
		t.Fatalf("expected %d results, got %d", len(enrichedPerfumes), len(result))
	}

	// Verify ranks are assigned correctly
	for i, r := range result {
		if r.Rank != i+1 {
			t.Fatalf("result[%d]: expected rank %d, got %d", i, i+1, r.Rank)
		}
		if !r.Perfume.Equal(enrichedPerfumes[i]) {
			t.Fatalf("result[%d]: expected perfume %+v, got %+v", i, enrichedPerfumes[i], r.Perfume)
		}
	}
}

func TestAiAdvisor_Advise_AdviseFetcherFails(t *testing.T) {
	t.Parallel()

	adviseFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return nil, false
		},
	}

	enrichFetcher := &MockFetcher{}
	advisor := NewAI(adviseFetcher, enrichFetcher)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err == nil {
		t.Fatal("expected error when adviseFetcher fails")
	}
	if err.Error() != "failed to get AI suggestions" {
		t.Fatalf("expected error 'failed to get AI suggestions', got %q", err.Error())
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestAiAdvisor_Advise_AdviseFetcherReturnsEmpty(t *testing.T) {
	t.Parallel()

	adviseFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return []perfume.Perfume{}, true
		},
	}

	enrichFetcher := &MockFetcher{}
	advisor := NewAI(adviseFetcher, enrichFetcher)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err == nil {
		t.Fatal("expected error when adviseFetcher returns empty")
	}
	if err.Error() != "failed to get AI suggestions" {
		t.Fatalf("expected error 'failed to get AI suggestions', got %q", err.Error())
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestAiAdvisor_Advise_EnrichFetcherFails(t *testing.T) {
	t.Parallel()

	aiSuggestions := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
	}

	adviseFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return aiSuggestions, true
		},
	}

	enrichFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return nil, false
		},
	}

	advisor := NewAI(adviseFetcher, enrichFetcher)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err == nil {
		t.Fatal("expected error when enrichFetcher fails")
	}
	if err.Error() != "failed to get enrichment results" {
		t.Fatalf("expected error 'failed to get enrichment results', got %q", err.Error())
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestAiAdvisor_Advise_EnrichFetcherReturnsEmpty(t *testing.T) {
	t.Parallel()

	aiSuggestions := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
	}

	adviseFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return aiSuggestions, true
		},
	}

	enrichFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return []perfume.Perfume{}, true
		},
	}

	advisor := NewAI(adviseFetcher, enrichFetcher)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err == nil {
		t.Fatal("expected error when enrichFetcher returns empty")
	}
	if err.Error() != "failed to get enrichment results" {
		t.Fatalf("expected error 'failed to get enrichment results', got %q", err.Error())
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestAiAdvisor_Advise_MultipleSuggestions(t *testing.T) {
	t.Parallel()

	aiSuggestions := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}

	enrichedPerfumes := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "J'adore", Sex: "female"},
		{Brand: "Tom Ford", Name: "Black Orchid", Sex: "unisex"},
	}

	adviseFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return aiSuggestions, true
		},
	}

	enrichFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return enrichedPerfumes, true
		},
	}

	advisor := NewAI(adviseFetcher, enrichFetcher)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female",
	}

	result, err := advisor.Advise(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != len(enrichedPerfumes) {
		t.Fatalf("expected %d results, got %d", len(enrichedPerfumes), len(result))
	}

	// Verify all results have correct ranks
	for i, r := range result {
		if r.Rank != i+1 {
			t.Fatalf("result[%d]: expected rank %d, got %d", i, i+1, r.Rank)
		}
	}
}

func TestAiAdvisor_Advise_EnrichmentParamsConstruction(t *testing.T) {
	t.Parallel()

	aiSuggestions := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "Sauvage", Sex: "male"},
	}

	enrichedPerfumes := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "Sauvage", Sex: "male"},
	}

	adviseFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			return aiSuggestions, true
		},
	}

	enrichFetcher := &MockFetcher{
		FetchFunc: func(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
			// Verify that enrichment params preserve brand and name from suggestions
			// but use sex from original params
			if len(params) != 2 {
				t.Fatalf("expected 2 enrichment params, got %d", len(params))
			}

			// First param should be Chanel/No5 with sex from original params
			if params[0].Brand != "Chanel" || params[0].Name != "No5" {
				t.Fatalf("expected first param to be Chanel/No5, got %+v", params[0])
			}

			// Second param should be Dior/Sauvage with sex from original params
			if params[1].Brand != "Dior" || params[1].Name != "Sauvage" {
				t.Fatalf("expected second param to be Dior/Sauvage, got %+v", params[1])
			}

			// Both should have sex from original params (not from suggestions)
			// Note: WithSex only accepts "male" or "female", so "unisex" won't be set
			// But the params should still be constructed with the sex from original request
			for i, p := range params {
				if p.Sex != "female" {
					t.Fatalf("param[%d]: expected sex 'female' from original params, got %q", i, p.Sex)
				}
			}

			return enrichedPerfumes, true
		},
	}

	advisor := NewAI(adviseFetcher, enrichFetcher)
	params := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   "female", // Will be used in enrichment params
	}

	result, err := advisor.Advise(params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != len(enrichedPerfumes) {
		t.Fatalf("expected %d results, got %d", len(enrichedPerfumes), len(result))
	}
}

