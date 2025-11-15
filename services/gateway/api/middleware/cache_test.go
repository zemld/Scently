package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
)

// MockCacher is a mock implementation of cache for testing
type MockCacher struct {
	store map[string]interface{}
}

func NewMockCacher() *MockCacher {
	return &MockCacher{
		store: make(map[string]interface{}),
	}
}

func (m *MockCacher) Save(ctx context.Context, key string, value interface{}) error {
	m.store[key] = value
	return nil
}

func (m *MockCacher) Load(ctx context.Context, key string) (interface{}, error) {
	value, ok := m.store[key]
	if !ok {
		return nil, nil
	}
	return value, nil
}

func (m *MockCacher) Clear() {
	m.store = make(map[string]interface{})
}

func TestGetCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "all params",
			query:    "brand=Chanel&name=No.5&sex=female&use_ai=true",
			expected: "Chanel:No.5:female:true",
		},
		{
			name:     "missing params",
			query:    "brand=Chanel",
			expected: "Chanel:::",
		},
		{
			name:     "empty query",
			query:    "",
			expected: ":::",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/perfume/suggest?"+tt.query, nil)
			key := getCacheKey(*req)
			if key != tt.expected {
				t.Errorf("expected cache key '%s', got '%s'", tt.expected, key)
			}
		})
	}
}

func TestGetTTL(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{
			name:     "valid duration",
			envValue: "2h",
			expected: 2 * time.Hour,
		},
		{
			name:     "invalid duration",
			envValue: "invalid",
			expected: defaultTTL,
		},
		{
			name:     "empty env",
			envValue: "",
			expected: defaultTTL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalTTL := ttlEnv
			os.Setenv("TTL_SECONDS", tt.envValue)
			ttlEnv = tt.envValue
			defer func() {
				os.Unsetenv("TTL_SECONDS")
				ttlEnv = originalTTL
			}()

			ttl := getTTL()
			if ttl != tt.expected {
				t.Errorf("expected TTL %v, got %v", tt.expected, ttl)
			}
		})
	}
}

func TestCache_CacheHit(t *testing.T) {
	mockCache := NewMockCacher()

	// Pre-populate cache
	cachedSuggestions := perfume.Suggestions{
		Perfumes: []perfume.Ranked{
			{
				Perfume: perfume.Perfume{
					Brand: "Chanel",
					Name:  "No.5",
					Sex:   "female",
				},
				Rank:  1,
				Score: 0.95,
			},
		},
	}
	mockCache.Save(context.Background(), "Chanel:No.5:female:", cachedSuggestions)

	// Test the cache key generation and load logic
	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest?brand=Chanel&name=No.5&sex=female", nil)
	key := getCacheKey(*req)

	cached, err := mockCache.Load(context.Background(), key)
	if err != nil {
		t.Fatalf("unexpected error loading from cache: %v", err)
	}

	if cached == nil {
		t.Error("expected cached value, got nil")
	}

	suggestions, ok := cached.(perfume.Suggestions)
	if !ok {
		t.Error("expected perfume.Suggestions type")
	}

	if len(suggestions.Perfumes) != 1 {
		t.Errorf("expected 1 perfume, got %d", len(suggestions.Perfumes))
	}
}

func TestCache_CacheMiss(t *testing.T) {
	mockCache := NewMockCacher()

	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest?brand=Chanel&name=No.5", nil)
	key := getCacheKey(*req)

	cached, err := mockCache.Load(context.Background(), key)
	if err != nil {
		t.Fatalf("unexpected error loading from cache: %v", err)
	}

	if cached != nil {
		t.Error("expected nil for cache miss, got value")
	}
}

func TestCache_SaveAndLoad(t *testing.T) {
	mockCache := NewMockCacher()

	suggestions := perfume.Suggestions{
		Perfumes: []perfume.Ranked{
			{
				Perfume: perfume.Perfume{
					Brand: "Dior",
					Name:  "Sauvage",
					Sex:   "male",
				},
				Rank:  1,
				Score: 0.92,
			},
		},
	}

	key := "Dior:Sauvage:male:"

	// Save
	err := mockCache.Save(context.Background(), key, suggestions)
	if err != nil {
		t.Fatalf("unexpected error saving to cache: %v", err)
	}

	// Load
	loaded, err := mockCache.Load(context.Background(), key)
	if err != nil {
		t.Fatalf("unexpected error loading from cache: %v", err)
	}

	if loaded == nil {
		t.Fatal("expected loaded value, got nil")
	}

	loadedSuggestions, ok := loaded.(perfume.Suggestions)
	if !ok {
		t.Fatal("expected perfume.Suggestions type")
	}

	if len(loadedSuggestions.Perfumes) != 1 {
		t.Errorf("expected 1 perfume, got %d", len(loadedSuggestions.Perfumes))
	}

	if loadedSuggestions.Perfumes[0].Perfume.Brand != "Dior" {
		t.Errorf("expected brand 'Dior', got '%s'", loadedSuggestions.Perfumes[0].Perfume.Brand)
	}
}

func TestCache_JSONEncoding(t *testing.T) {
	suggestions := perfume.Suggestions{
		Perfumes: []perfume.Ranked{
			{
				Perfume: perfume.Perfume{
					Brand: "Test",
					Name:  "Perfume",
					Sex:   "unisex",
				},
				Rank:  1,
				Score: 0.9,
			},
		},
	}

	// Test JSON encoding
	data, err := json.Marshal(suggestions)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}

	// Test JSON decoding
	var decoded perfume.Suggestions
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	if len(decoded.Perfumes) != 1 {
		t.Errorf("expected 1 perfume, got %d", len(decoded.Perfumes))
	}

	if decoded.Perfumes[0].Perfume.Brand != "Test" {
		t.Errorf("expected brand 'Test', got '%s'", decoded.Perfumes[0].Perfume.Brand)
	}
}
