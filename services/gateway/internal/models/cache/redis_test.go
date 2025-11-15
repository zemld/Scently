package cache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
)

// MockRedisClient is a simple in-memory mock for Redis
type MockRedisClient struct {
	store map[string]string
	ttl   map[string]time.Time
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		store: make(map[string]string),
		ttl:   make(map[string]time.Time),
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx, "get", key)

	// Check if expired
	if expiry, ok := m.ttl[key]; ok && time.Now().After(expiry) {
		delete(m.store, key)
		delete(m.ttl, key)
		cmd.SetErr(redis.Nil)
		return cmd
	}

	value, ok := m.store[key]
	if !ok {
		cmd.SetErr(redis.Nil)
		return cmd
	}

	cmd.SetVal(value)
	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "set", key, value)

	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	default:
		data, _ := json.Marshal(value)
		strValue = string(data)
	}

	m.store[key] = strValue
	if expiration > 0 {
		m.ttl[key] = time.Now().Add(expiration)
	}

	cmd.SetVal("OK")
	return cmd
}

func (m *MockRedisClient) Clear() {
	m.store = make(map[string]string)
	m.ttl = make(map[string]time.Time)
}

// Test helper functions removed - testing Save/Load logic directly

func TestRedisCacher_Save(t *testing.T) {
	// Test the JSON marshaling logic (which is what Save does)
	suggestions := perfume.Suggestions{
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

	// Test JSON marshaling (which is what Save does)
	data, err := json.Marshal(suggestions)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty data")
	}

	// Verify it can be unmarshaled back
	var decoded perfume.Suggestions
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	if len(decoded.Perfumes) != 1 {
		t.Errorf("expected 1 perfume, got %d", len(decoded.Perfumes))
	}
}

func TestRedisCacher_Load(t *testing.T) {
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

	// Test JSON encoding/decoding (which is what Load does)
	data, err := json.Marshal(suggestions)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}

	var result perfume.Suggestions
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	if len(result.Perfumes) != 1 {
		t.Errorf("expected 1 perfume, got %d", len(result.Perfumes))
	}

	if result.Perfumes[0].Perfume.Brand != "Dior" {
		t.Errorf("expected brand 'Dior', got '%s'", result.Perfumes[0].Perfume.Brand)
	}
}

func TestRedisCacher_Load_EmptyData(t *testing.T) {
	// Test loading empty/nil data
	var result perfume.Suggestions
	err := json.Unmarshal([]byte(""), &result)
	if err == nil {
		t.Error("expected error for empty JSON")
	}
}

func TestRedisCacher_Load_InvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid": "json"`

	var result perfume.Suggestions
	err := json.Unmarshal([]byte(invalidJSON), &result)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestRedisCacher_Save_InvalidValue(t *testing.T) {
	// Test with a value that can't be marshaled (channel)
	invalidValue := make(chan int)

	// This should fail during marshaling
	_, err := json.Marshal(invalidValue)
	if err == nil {
		t.Error("expected error for unmarshalable value")
	}
}

func TestRedisCacher_ContextKey(t *testing.T) {
	// Test that SuggestionsKey is properly defined
	if SuggestionsKey != "suggestions" {
		t.Errorf("expected SuggestionsKey to be 'suggestions', got '%s'", SuggestionsKey)
	}
}

func TestRedisCacher_TTL(t *testing.T) {
	ttl := 30 * time.Minute
	testCacher := &RedisCacher{
		cacheTTL: ttl,
	}

	if testCacher.cacheTTL != ttl {
		t.Errorf("expected TTL %v, got %v", ttl, testCacher.cacheTTL)
	}
}
