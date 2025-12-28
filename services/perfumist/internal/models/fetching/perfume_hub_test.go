package fetching

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
	protoModels "github.com/zemld/Scently/generated/proto/perfume-hub/models"
	"github.com/zemld/Scently/generated/proto/perfume-hub/requests"
	"google.golang.org/grpc"
)

// mockPerfumeStorageClient is a mock implementation of PerfumeStorageClient
type mockPerfumeStorageClient struct {
	GetPerfumesFunc func(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error)
}

func (m *mockPerfumeStorageClient) GetPerfumes(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error) {
	if m.GetPerfumesFunc != nil {
		return m.GetPerfumesFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *mockPerfumeStorageClient) UpdatePerfumes(ctx context.Context, in *requests.UpdatePerfumesRequest, opts ...grpc.CallOption) (*requests.UpdatePerfumesResponse, error) {
	return nil, nil
}

func TestNewPerfumeHub(t *testing.T) {
	t.Parallel()

	client := &mockPerfumeStorageClient{}
	hub := NewPerfumeHub(client)

	if hub == nil {
		t.Fatal("expected non-nil PerfumeHub")
	}
	if hub.client != client {
		t.Fatal("expected client to be set")
	}
}

func TestPerfumeHub_Fetch_Success(t *testing.T) {
	t.Parallel()

	expectedPerfumes := []*protoModels.Perfume{
		{
			Brand: "Chanel",
			Name:  "No5",
			Sex:   "female",
		},
		{
			Brand: "Dior",
			Name:  "Sauvage",
			Sex:   "male",
		},
	}

	callCount := 0
	client := &mockPerfumeStorageClient{
		GetPerfumesFunc: func(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error) {
			callCount++
			return &requests.GetPerfumesResponse{
				Perfumes: expectedPerfumes,
			}, nil
		},
	}

	hub := NewPerfumeHub(client)
	params := []parameters.RequestPerfume{
		{Brand: "Chanel", Name: "No5"},
		{Brand: "Dior", Name: "Sauvage"},
	}

	perfumes, ok := hub.Fetch(context.Background(), params)

	if !ok {
		t.Fatal("expected true on success")
	}
	if callCount != 2 {
		t.Fatalf("expected 2 calls, got %d", callCount)
	}
	if len(perfumes) != 4 {
		t.Fatalf("expected 4 perfumes (2 from each request), got %d", len(perfumes))
	}
}

func TestPerfumeHub_Fetch_EmptyResults(t *testing.T) {
	t.Parallel()

	client := &mockPerfumeStorageClient{
		GetPerfumesFunc: func(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error) {
			return &requests.GetPerfumesResponse{
				Perfumes: []*protoModels.Perfume{},
			}, nil
		},
	}

	hub := NewPerfumeHub(client)
	params := []parameters.RequestPerfume{
		{Brand: "Chanel", Name: "No5"},
	}

	perfumes, ok := hub.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on empty results")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestPerfumeHub_Fetch_ClientError(t *testing.T) {
	t.Parallel()

	client := &mockPerfumeStorageClient{
		GetPerfumesFunc: func(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error) {
			return nil, context.DeadlineExceeded
		},
	}

	hub := NewPerfumeHub(client)
	params := []parameters.RequestPerfume{
		{Brand: "Chanel", Name: "No5"},
	}

	perfumes, ok := hub.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on client error")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestPerfumeHub_Fetch_Timeout(t *testing.T) {
	client := &mockPerfumeStorageClient{
		GetPerfumesFunc: func(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error) {
			// Simulate timeout by waiting longer than the fetcher's timeout
			time.Sleep(config.PerfumeHubFetcherTimeout + 100*time.Millisecond)
			return &requests.GetPerfumesResponse{
				Perfumes: []*protoModels.Perfume{{Brand: "Chanel", Name: "No5"}},
			}, nil
		},
	}

	hub := NewPerfumeHub(client)
	params := []parameters.RequestPerfume{
		{Brand: "Chanel", Name: "No5"},
	}

	perfumes, ok := hub.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on timeout")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestPerfumeHub_getPerfumesAsync_Success(t *testing.T) {
	t.Parallel()

	expectedResponse := &requests.GetPerfumesResponse{
		Perfumes: []*protoModels.Perfume{
			{Brand: "Chanel", Name: "No5", Sex: "female"},
		},
	}

	client := &mockPerfumeStorageClient{
		GetPerfumesFunc: func(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error) {
			return expectedResponse, nil
		},
	}

	hub := NewPerfumeHub(client)
	results := make(chan *requests.GetPerfumesResponse, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	hub.getPerfumesAsync(context.Background(), parameters.RequestPerfume{Brand: "Chanel", Name: "No5"}, results, wg)
	wg.Wait()
	close(results)

	result := <-results
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Perfumes) != 1 {
		t.Fatalf("expected 1 perfume, got %d", len(result.Perfumes))
	}
	if result.Perfumes[0].Brand != "Chanel" {
		t.Fatalf("expected brand Chanel, got %s", result.Perfumes[0].Brand)
	}
}

func TestPerfumeHub_getPerfumesAsync_Error(t *testing.T) {
	t.Parallel()

	client := &mockPerfumeStorageClient{
		GetPerfumesFunc: func(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error) {
			return nil, context.DeadlineExceeded
		},
	}

	hub := NewPerfumeHub(client)
	results := make(chan *requests.GetPerfumesResponse, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	hub.getPerfumesAsync(context.Background(), parameters.RequestPerfume{Brand: "Chanel"}, results, wg)
	wg.Wait()
	close(results)

	result := <-results
	if result != nil {
		t.Fatalf("expected nil result on error, got %v", result)
	}
}

func TestPerfumeHub_getPerfumesAsync_NilWaitGroup(t *testing.T) {
	t.Parallel()

	client := &mockPerfumeStorageClient{
		GetPerfumesFunc: func(ctx context.Context, in *requests.GetPerfumesRequest, opts ...grpc.CallOption) (*requests.GetPerfumesResponse, error) {
			return &requests.GetPerfumesResponse{Perfumes: []*protoModels.Perfume{}}, nil
		},
	}

	hub := NewPerfumeHub(client)
	results := make(chan *requests.GetPerfumesResponse, 1)

	// Should not panic with nil waitgroup
	hub.getPerfumesAsync(context.Background(), parameters.RequestPerfume{Brand: "Chanel"}, results, nil)
	close(results)

	result := <-results
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestPerfumeHub_fetchPerfumeResults_Success(t *testing.T) {
	t.Parallel()

	ch := make(chan *requests.GetPerfumesResponse, 3)
	ch <- &requests.GetPerfumesResponse{
		Perfumes: []*protoModels.Perfume{
			{Brand: "Chanel", Name: "No5", Sex: "female"},
		},
	}
	ch <- &requests.GetPerfumesResponse{
		Perfumes: []*protoModels.Perfume{
			{Brand: "Dior", Name: "Sauvage", Sex: "male"},
		},
	}
	ch <- &requests.GetPerfumesResponse{
		Perfumes: []*protoModels.Perfume{},
	}
	close(ch)

	hub := NewPerfumeHub(&mockPerfumeStorageClient{})
	perfumes, status := hub.fetchPerfumeResults(context.Background(), ch)

	if status != 0 {
		t.Fatalf("expected status 0, got %d", status)
	}
	if len(perfumes) != 2 {
		t.Fatalf("expected 2 perfumes, got %d", len(perfumes))
	}
}

func TestPerfumeHub_fetchPerfumeResults_NilResponse(t *testing.T) {
	t.Parallel()

	ch := make(chan *requests.GetPerfumesResponse, 2)
	ch <- nil
	ch <- &requests.GetPerfumesResponse{
		Perfumes: []*protoModels.Perfume{
			{Brand: "Chanel", Name: "No5"},
		},
	}
	close(ch)

	hub := NewPerfumeHub(&mockPerfumeStorageClient{})
	perfumes, status := hub.fetchPerfumeResults(context.Background(), ch)

	if status != 0 {
		t.Fatalf("expected status 0, got %d", status)
	}
	if len(perfumes) != 1 {
		t.Fatalf("expected 1 perfume, got %d", len(perfumes))
	}
}

func TestPerfumeHub_fetchPerfumeResults_NilPerfumeInResponse(t *testing.T) {
	t.Parallel()

	ch := make(chan *requests.GetPerfumesResponse, 1)
	ch <- &requests.GetPerfumesResponse{
		Perfumes: []*protoModels.Perfume{
			nil,
			{Brand: "Chanel", Name: "No5"},
		},
	}
	close(ch)

	hub := NewPerfumeHub(&mockPerfumeStorageClient{})
	perfumes, status := hub.fetchPerfumeResults(context.Background(), ch)

	if status != 0 {
		t.Fatalf("expected status 0, got %d", status)
	}
	if len(perfumes) != 1 {
		t.Fatalf("expected 1 perfume, got %d", len(perfumes))
	}
}

func TestPerfumeHub_fetchPerfumeResults_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ch := make(chan *requests.GetPerfumesResponse)
	go func() {
		time.Sleep(10 * time.Millisecond)
		close(ch)
	}()

	hub := NewPerfumeHub(&mockPerfumeStorageClient{})
	perfumes, status := hub.fetchPerfumeResults(ctx, ch)

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestConvertPerfumeToModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *protoModels.Perfume
		expected perfume.Perfume
	}{
		{
			name: "full perfume",
			input: &protoModels.Perfume{
				Brand:    "Chanel",
				Name:     "No5",
				Sex:      "female",
				ImageUrl: stringPtr("http://example.com/image.jpg"),
				Properties: &protoModels.Perfume_Properties{
					PerfumeType: "EDP",
					Family:      []string{"floral"},
					UpperNotes:  []string{"bergamot"},
					CoreNotes:   []string{"rose"},
					BaseNotes:   []string{"sandalwood"},
				},
				Shops: []*protoModels.Perfume_ShopInfo{
					{
						ShopName: "Shop1",
						Domain:   "shop1.com",
						ImageUrl: stringPtr("http://shop1.com/image.jpg"),
						Variants: []*protoModels.Perfume_ShopInfo_Variant{
							{Volume: 50, Link: "http://shop1.com/product", Price: 1000},
						},
					},
				},
			},
			expected: perfume.Perfume{
				Brand:    "Chanel",
				Name:     "No5",
				Sex:      "female",
				ImageUrl: "http://example.com/image.jpg",
				Properties: perfume.Properties{
					Type:       "EDP",
					Family:     []string{"floral"},
					UpperNotes: []string{"bergamot"},
					CoreNotes:  []string{"rose"},
					BaseNotes:  []string{"sandalwood"},
				},
				Shops: []perfume.ShopInfo{
					{
						ShopName: "Shop1",
						Domain:   "shop1.com",
						ImageUrl: "http://shop1.com/image.jpg",
						Variants: []perfume.Variant{
							{Volume: 50, Link: "http://shop1.com/product", Price: 1000},
						},
					},
				},
			},
		},
		{
			name:     "nil perfume",
			input:    nil,
			expected: perfume.Perfume{},
		},
		{
			name: "minimal perfume",
			input: &protoModels.Perfume{
				Brand: "Chanel",
				Name:  "No5",
				Sex:   "female",
			},
			expected: perfume.Perfume{
				Brand:      "Chanel",
				Name:       "No5",
				Sex:        "female",
				ImageUrl:    "",
				Properties: perfume.Properties{},
				Shops:      []perfume.ShopInfo{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertPerfumeToModel(tt.input)
			if result.Brand != tt.expected.Brand {
				t.Errorf("Brand: expected %q, got %q", tt.expected.Brand, result.Brand)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Name: expected %q, got %q", tt.expected.Name, result.Name)
			}
			if result.Sex != tt.expected.Sex {
				t.Errorf("Sex: expected %q, got %q", tt.expected.Sex, result.Sex)
			}
			if result.ImageUrl != tt.expected.ImageUrl {
				t.Errorf("ImageUrl: expected %q, got %q", tt.expected.ImageUrl, result.ImageUrl)
			}
		})
	}
}

func TestConvertPropertiesToModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *protoModels.Perfume_Properties
		expected perfume.Properties
	}{
		{
			name: "full properties",
			input: &protoModels.Perfume_Properties{
				PerfumeType: "EDP",
				Family:      []string{"floral", "woody"},
				UpperNotes:  []string{"bergamot", "lemon"},
				CoreNotes:   []string{"rose", "jasmine"},
				BaseNotes:   []string{"sandalwood", "musk"},
			},
			expected: perfume.Properties{
				Type:       "EDP",
				Family:     []string{"floral", "woody"},
				UpperNotes: []string{"bergamot", "lemon"},
				CoreNotes:  []string{"rose", "jasmine"},
				BaseNotes:  []string{"sandalwood", "musk"},
			},
		},
		{
			name:     "nil properties",
			input:    nil,
			expected: perfume.Properties{},
		},
		{
			name: "empty properties",
			input: &protoModels.Perfume_Properties{
				PerfumeType: "",
				Family:      []string{},
				UpperNotes:  []string{},
				CoreNotes:   []string{},
				BaseNotes:   []string{},
			},
			expected: perfume.Properties{
				Type:       "",
				Family:     []string{},
				UpperNotes: []string{},
				CoreNotes:  []string{},
				BaseNotes:  []string{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertPropertiesToModel(tt.input)
			if result.Type != tt.expected.Type {
				t.Errorf("Type: expected %q, got %q", tt.expected.Type, result.Type)
			}
			if len(result.Family) != len(tt.expected.Family) {
				t.Errorf("Family length: expected %d, got %d", len(tt.expected.Family), len(result.Family))
			}
		})
	}
}

func TestConvertShopInfoToModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []*protoModels.Perfume_ShopInfo
		expected []perfume.ShopInfo
	}{
		{
			name: "full shop info",
			input: []*protoModels.Perfume_ShopInfo{
				{
					ShopName: "Shop1",
					Domain:   "shop1.com",
					ImageUrl: stringPtr("http://shop1.com/image.jpg"),
					Variants: []*protoModels.Perfume_ShopInfo_Variant{
						{Volume: 50, Link: "http://shop1.com/product1", Price: 1000},
						{Volume: 100, Link: "http://shop1.com/product2", Price: 1800},
					},
				},
				{
					ShopName: "Shop2",
					Domain:   "shop2.com",
					ImageUrl: nil,
					Variants: []*protoModels.Perfume_ShopInfo_Variant{},
				},
			},
			expected: []perfume.ShopInfo{
				{
					ShopName: "Shop1",
					Domain:   "shop1.com",
					ImageUrl: "http://shop1.com/image.jpg",
					Variants: []perfume.Variant{
						{Volume: 50, Link: "http://shop1.com/product1", Price: 1000},
						{Volume: 100, Link: "http://shop1.com/product2", Price: 1800},
					},
				},
				{
					ShopName: "Shop2",
					Domain:   "shop2.com",
					ImageUrl: "",
					Variants: []perfume.Variant{},
				},
			},
		},
		{
			name:     "nil shops",
			input:    nil,
			expected: []perfume.ShopInfo{},
		},
		{
			name:     "empty shops",
			input:    []*protoModels.Perfume_ShopInfo{},
			expected: []perfume.ShopInfo{},
		},
		{
			name: "shop with nil variant",
			input: []*protoModels.Perfume_ShopInfo{
				{
					ShopName: "Shop1",
					Domain:   "shop1.com",
					Variants: []*protoModels.Perfume_ShopInfo_Variant{
						nil,
						{Volume: 50, Link: "http://shop1.com/product", Price: 1000},
					},
				},
			},
			expected: []perfume.ShopInfo{
				{
					ShopName: "Shop1",
					Domain:   "shop1.com",
					ImageUrl: "",
					Variants: []perfume.Variant{
						{Volume: 0, Link: "", Price: 0}, // nil variant creates zero value
						{Volume: 50, Link: "http://shop1.com/product", Price: 1000},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertShopInfoToModel(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("length: expected %d, got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i].ShopName != tt.expected[i].ShopName {
					t.Errorf("ShopName[%d]: expected %q, got %q", i, tt.expected[i].ShopName, result[i].ShopName)
				}
				if result[i].Domain != tt.expected[i].Domain {
					t.Errorf("Domain[%d]: expected %q, got %q", i, tt.expected[i].Domain, result[i].Domain)
				}
			}
		})
	}
}

func TestConvertVariantsToModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []*protoModels.Perfume_ShopInfo_Variant
		expected []perfume.Variant
	}{
		{
			name: "full variants",
			input: []*protoModels.Perfume_ShopInfo_Variant{
				{Volume: 50, Link: "http://shop.com/product1", Price: 1000},
				{Volume: 100, Link: "http://shop.com/product2", Price: 1800},
			},
			expected: []perfume.Variant{
				{Volume: 50, Link: "http://shop.com/product1", Price: 1000},
				{Volume: 100, Link: "http://shop.com/product2", Price: 1800},
			},
		},
		{
			name:     "nil variants",
			input:    nil,
			expected: []perfume.Variant{},
		},
		{
			name:     "empty variants",
			input:    []*protoModels.Perfume_ShopInfo_Variant{},
			expected: []perfume.Variant{},
		},
		{
			name: "variants with nil",
			input: []*protoModels.Perfume_ShopInfo_Variant{
				nil,
				{Volume: 50, Link: "http://shop.com/product", Price: 1000},
			},
			expected: []perfume.Variant{
				{Volume: 0, Link: "", Price: 0}, // nil variant creates zero value
				{Volume: 50, Link: "http://shop.com/product", Price: 1000},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := convertVariantsToModel(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("length: expected %d, got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i].Volume != tt.expected[i].Volume {
					t.Errorf("Volume[%d]: expected %d, got %d", i, tt.expected[i].Volume, result[i].Volume)
				}
				if result[i].Link != tt.expected[i].Link {
					t.Errorf("Link[%d]: expected %q, got %q", i, tt.expected[i].Link, result[i].Link)
				}
				if result[i].Price != tt.expected[i].Price {
					t.Errorf("Price[%d]: expected %d, got %d", i, tt.expected[i].Price, result[i].Price)
				}
			}
		})
	}
}

func TestTryConvertPointer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "non-nil pointer",
			input:    stringPtr("test"),
			expected: "test",
		},
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty string",
			input:    stringPtr(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tryConvertPointer(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

