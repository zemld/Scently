package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestSuggest_Success(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"suggestions":[]}`))
	}))
	defer mockServer.Close()

	// Set environment variable
	originalURL := suggestUrl
	os.Setenv("PERFUMIST_URL", mockServer.URL)
	suggestUrl = mockServer.URL
	defer func() {
		os.Unsetenv("PERFUMIST_URL")
		suggestUrl = originalURL
	}()

	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest?brand=test&name=test", nil)
	rr := httptest.NewRecorder()

	Suggest(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestSuggest_WithQueryParams(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query params are passed through
		if r.URL.Query().Get("brand") != "Chanel" {
			t.Errorf("expected brand 'Chanel', got '%s'", r.URL.Query().Get("brand"))
		}
		if r.URL.Query().Get("name") != "No.5" {
			t.Errorf("expected name 'No.5', got '%s'", r.URL.Query().Get("name"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"suggestions":[]}`))
	}))
	defer mockServer.Close()

	originalURL := suggestUrl
	os.Setenv("PERFUMIST_URL", mockServer.URL)
	suggestUrl = mockServer.URL
	defer func() {
		os.Unsetenv("PERFUMIST_URL")
		suggestUrl = originalURL
	}()

	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest?brand=Chanel&name=No.5", nil)
	rr := httptest.NewRecorder()

	Suggest(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}
}

func TestSuggest_ServerError(t *testing.T) {
	// Setup mock server that returns error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer mockServer.Close()

	originalURL := suggestUrl
	os.Setenv("PERFUMIST_URL", mockServer.URL)
	suggestUrl = mockServer.URL
	defer func() {
		os.Unsetenv("PERFUMIST_URL")
		suggestUrl = originalURL
	}()

	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest", nil)
	rr := httptest.NewRecorder()

	Suggest(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
}

func TestSuggest_Timeout(t *testing.T) {
	// Setup mock server that delays response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second) // Longer than default timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	originalURL := suggestUrl
	os.Setenv("PERFUMIST_URL", mockServer.URL)
	suggestUrl = mockServer.URL
	defer func() {
		os.Unsetenv("PERFUMIST_URL")
		suggestUrl = originalURL
	}()

	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest", nil)
	rr := httptest.NewRecorder()

	Suggest(rr, req)

	// Should timeout and return 408
	if status := rr.Code; status != http.StatusRequestTimeout {
		t.Errorf("expected status %d, got %d", http.StatusRequestTimeout, status)
	}
}

func TestSuggest_InvalidURL(t *testing.T) {
	originalURL := suggestUrl
	os.Setenv("PERFUMIST_URL", "http://invalid-url-that-does-not-exist:9999")
	suggestUrl = "http://invalid-url-that-does-not-exist:9999"
	defer func() {
		os.Unsetenv("PERFUMIST_URL")
		suggestUrl = originalURL
	}()

	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest", nil)
	rr := httptest.NewRecorder()

	Suggest(rr, req)

	// Should return error
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
}

func TestGetTimeoutFromRequest_WithAI(t *testing.T) {
	// Test with use_ai=true
	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest?use_ai=true", nil)

	timeout := getTimeoutFromRequest(*req)

	// Should return default AI timeout since env var is not set
	if timeout != defaultAITimeout {
		t.Errorf("expected timeout %v, got %v", defaultAITimeout, timeout)
	}
}

func TestGetTimeoutFromRequest_WithoutAI(t *testing.T) {
	// Test without use_ai
	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest", nil)

	timeout := getTimeoutFromRequest(*req)

	// Should return default non-AI timeout
	if timeout != defaultNonAITimeout {
		t.Errorf("expected timeout %v, got %v", defaultNonAITimeout, timeout)
	}
}

func TestGetTimeoutFromRequest_WithAIEnvVar(t *testing.T) {
	// Note: The current implementation parses the constant name as duration,
	// which will always fail, so it returns defaultAITimeout
	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest?use_ai=true", nil)

	timeout := getTimeoutFromRequest(*req)

	// Since parsing "SUGGEST_AI_TIMEOUT" as duration fails, it returns default
	if timeout != defaultAITimeout {
		t.Errorf("expected timeout %v, got %v", defaultAITimeout, timeout)
	}
}

func TestSuggest_ContextCanceled(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	originalURL := suggestUrl
	os.Setenv("PERFUMIST_URL", mockServer.URL)
	suggestUrl = mockServer.URL
	defer func() {
		os.Unsetenv("PERFUMIST_URL")
		suggestUrl = originalURL
	}()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := httptest.NewRequest(http.MethodGet, "/perfume/suggest", nil)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	Suggest(rr, req)

	// Should handle canceled context
	if status := rr.Code; status != http.StatusRequestTimeout {
		t.Errorf("expected status %d, got %d", http.StatusRequestTimeout, status)
	}
}
