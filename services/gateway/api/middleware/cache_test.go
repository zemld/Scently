package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
)

// MockCacher is a mock implementation of cache for testing
type MockCacher struct {
	store map[string][]byte
}

func NewMockCacher() *MockCacher {
	return &MockCacher{
		store: make(map[string][]byte),
	}
}

func (m *MockCacher) Save(ctx context.Context, key string, value []byte) error {
	m.store[key] = value
	return nil
}

func (m *MockCacher) Load(ctx context.Context, key string) ([]byte, error) {
	value, ok := m.store[key]
	if !ok {
		return nil, nil
	}
	return value, nil
}

func (m *MockCacher) Close() error {
	return nil
}

func (m *MockCacher) Clear() {
	m.store = make(map[string][]byte)
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
			expected: "chanelno5femaletrue", // canonizer нормализует строки
		},
		{
			name:     "missing params",
			query:    "brand=Chanel",
			expected: "chanel", // canonizer убирает пустые значения
		},
		{
			name:     "empty query",
			query:    "",
			expected: "", // canonizer убирает пустые значения
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
	// Тест проверяет, что getTTL() возвращает валидное значение
	// Без изменения глобальных переменных или переменных окружения
	ttl := getTTL()

	// Проверяем, что возвращается валидная длительность
	if ttl <= 0 {
		t.Errorf("expected positive TTL, got %v", ttl)
	}

	// Проверяем, что значение соответствует либо defaultTTL, либо значению из окружения
	// (но не меняем окружение для проверки)
	if ttl != defaultTTL && ttl < time.Second {
		t.Errorf("TTL seems invalid: %v", ttl)
	}
}

func TestTryLoadFromCache_CacheHit(t *testing.T) {
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
	cachedData, err := json.Marshal(cachedSuggestions)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}
	mockCache.Save(context.Background(), "Chanel:No.5:female:", cachedData)

	// Test tryLoadFromCache
	w := httptest.NewRecorder()

	cacheHit := tryLoadFromCache(context.Background(), mockCache, "Chanel:No.5:female:", w)

	if !cacheHit {
		t.Error("expected cache hit, got cache miss")
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response perfume.Suggestions
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected error unmarshaling response: %v", err)
	}

	if len(response.Perfumes) != 1 {
		t.Errorf("expected 1 perfume, got %d", len(response.Perfumes))
	}
}

func TestTryLoadFromCache_CacheMiss(t *testing.T) {
	mockCache := NewMockCacher()

	w := httptest.NewRecorder()
	cacheHit := tryLoadFromCache(context.Background(), mockCache, "nonexistent:key:", w)

	if cacheHit {
		t.Error("expected cache miss, got cache hit")
	}

	// httptest.NewRecorder() по умолчанию имеет Code = 200, но если WriteHeader не вызывался,
	// то это означает, что ответ не был записан
	if w.Body.Len() > 0 {
		t.Errorf("expected no response body, got %d bytes", w.Body.Len())
	}
}

func TestTryLoadFromCache_EmptySuggestions(t *testing.T) {
	mockCache := NewMockCacher()

	// Cache with empty suggestions
	emptySuggestions := perfume.Suggestions{
		Perfumes: []perfume.Ranked{},
	}
	cachedData, err := json.Marshal(emptySuggestions)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}
	mockCache.Save(context.Background(), "empty:key:", cachedData)

	w := httptest.NewRecorder()
	cacheHit := tryLoadFromCache(context.Background(), mockCache, "empty:key:", w)

	if cacheHit {
		t.Error("expected cache miss for empty suggestions, got cache hit")
	}
}

// Note: Full middleware tests require dependency injection or build tags
// These tests focus on testing individual functions that can be tested in isolation

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
	cachedData, err := json.Marshal(suggestions)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}

	err = mockCache.Save(context.Background(), key, cachedData)
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

	var loadedSuggestions perfume.Suggestions
	if err := json.Unmarshal(loaded, &loadedSuggestions); err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
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
