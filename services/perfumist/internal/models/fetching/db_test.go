package fetching

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

func TestNewDbFetcher(t *testing.T) {
	t.Parallel()

	url := "http://test-url:8080"
	token := "test-token"
	fetcher := NewDB(url, token)

	if fetcher == nil {
		t.Fatal("expected non-nil fetcher")
	}
	if fetcher.url != url {
		t.Fatalf("expected url %q, got %q", url, fetcher.url)
	}
	if fetcher.token != token {
		t.Fatalf("expected token %q, got %q", token, fetcher.token)
	}
	if fetcher.timeout != 2*time.Second {
		t.Fatalf("expected timeout %v, got %v", 2*time.Second, fetcher.timeout)
	}
}

func TestDbFetcher_getPerfumes_Success(t *testing.T) {
	expectedPerfumes := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "Sauvage", Sex: "male"},
	}
	resp := perfume.PerfumeResponse{Perfumes: expectedPerfumes}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}
	if len(perfumes) != len(expectedPerfumes) {
		t.Fatalf("expected %d perfumes, got %d", len(expectedPerfumes), len(perfumes))
	}
	for i, p := range expectedPerfumes {
		if !perfumes[i].Equal(p) {
			t.Fatalf("perfume %d: expected %+v, got %+v", i, p, perfumes[i])
		}
	}
}

func TestDbFetcher_getPerfumes_HTTPError(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_getPerfumes_BadStatus(t *testing.T) {
	resp := perfume.PerfumeResponse{Perfumes: []perfume.Perfume{}}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     http.StatusText(http.StatusInternalServerError),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_getPerfumes_ReadBodyError(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       errorReader{},
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_getPerfumes_NoContent(t *testing.T) {
	resp := perfume.PerfumeResponse{Perfumes: []perfume.Perfume{}}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}
	if perfumes == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(perfumes) != 0 {
		t.Fatalf("expected no perfumes, got %d", len(perfumes))
	}
}

func TestDbFetcher_getPerfumes_NotFound(t *testing.T) {
	resp := perfume.PerfumeResponse{Perfumes: []perfume.Perfume{}}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     http.StatusText(http.StatusNotFound),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, status)
	}
	if perfumes == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(perfumes) != 0 {
		t.Fatalf("expected no perfumes, got %d", len(perfumes))
	}
}

func TestDbFetcher_getPerfumes_Forbidden(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Status:     http.StatusText(http.StatusForbidden),
			Body:       io.NopCloser(strings.NewReader("Forbidden")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_getPerfumes_AddsAuthHeader(t *testing.T) {
	var capturedRequest *http.Request
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		capturedRequest = r
		resp := perfume.PerfumeResponse{Perfumes: []perfume.Perfume{{Brand: "Test", Name: "Test"}}}
		body, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if capturedRequest == nil {
		t.Fatal("expected request to be captured")
	}
	authHeader := capturedRequest.Header.Get("Authorization")
	expectedAuth := "Bearer test-token"
	if authHeader != expectedAuth {
		t.Fatalf("expected auth header %q, got %q", expectedAuth, authHeader)
	}
}

func TestDbFetcher_getPerfumesAsync_Success(t *testing.T) {
	expectedPerfumes := []perfume.Perfume{{Brand: "Chanel", Name: "No5"}}
	resp := perfume.PerfumeResponse{Perfumes: expectedPerfumes}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	results := make(chan perfumesFetchAndGlueResult, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	fetcher.getPerfumesAsync(context.Background(), parameters.RequestPerfume{}, results, wg)
	wg.Wait()
	close(results)

	result := <-results
	if result.Status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, result.Status)
	}
	if len(result.Perfumes) != len(expectedPerfumes) {
		t.Fatalf("expected %d perfumes, got %d", len(expectedPerfumes), len(result.Perfumes))
	}
}

func TestDbFetcher_getPerfumesAsync_Error(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	results := make(chan perfumesFetchAndGlueResult, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	fetcher.getPerfumesAsync(context.Background(), parameters.RequestPerfume{}, results, wg)
	wg.Wait()
	close(results)

	result := <-results
	if result.Status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, result.Status)
	}
	if result.Perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", result.Perfumes)
	}
}

func TestDbFetcher_fetchPerfumeResults_Success(t *testing.T) {
	t.Parallel()

	ch := make(chan perfumesFetchAndGlueResult, 3)
	ch <- perfumesFetchAndGlueResult{
		Status:   http.StatusOK,
		Perfumes: []perfume.Perfume{{Brand: "Chanel", Name: "No5"}},
	}
	ch <- perfumesFetchAndGlueResult{
		Status:   http.StatusOK,
		Perfumes: []perfume.Perfume{{Brand: "Dior", Name: "Sauvage"}},
	}
	ch <- perfumesFetchAndGlueResult{
		Status:   http.StatusNotFound,
		Perfumes: []perfume.Perfume{},
	}
	close(ch)

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.fetchPerfumeResults(context.Background(), ch)

	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}
	if len(perfumes) != 2 {
		t.Fatalf("expected 2 perfumes, got %d", len(perfumes))
	}
}

func TestDbFetcher_fetchPerfumeResults_ServerError(t *testing.T) {
	t.Parallel()

	ch := make(chan perfumesFetchAndGlueResult, 1)
	ch <- perfumesFetchAndGlueResult{Status: http.StatusInternalServerError}
	close(ch)

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.fetchPerfumeResults(context.Background(), ch)

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_fetchPerfumeResults_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fetcher := NewDB("http://test-url:8080", "test-token")
	perfumes, status := fetcher.fetchPerfumeResults(ctx, make(chan perfumesFetchAndGlueResult))

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_Fetch_Success(t *testing.T) {
	callCount := 0
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		callCount++
		perfumes := []perfume.Perfume{
			{Brand: "Chanel", Name: "No5"},
			{Brand: "Dior", Name: "Sauvage"},
		}
		resp := perfume.PerfumeResponse{Perfumes: perfumes}
		body, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	params := []parameters.RequestPerfume{
		{Brand: "Chanel"},
		{Brand: "Dior"},
	}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if !ok {
		t.Fatal("expected true on success")
	}
	if callCount != 2 {
		t.Fatalf("expected 2 HTTP calls, got %d", callCount)
	}
	if len(perfumes) != 4 {
		t.Fatalf("expected 4 perfumes (2 from each request), got %d", len(perfumes))
	}
}

func TestDbFetcher_Fetch_EmptyResults(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp := perfume.PerfumeResponse{Perfumes: []perfume.Perfume{}}
		body, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	params := []parameters.RequestPerfume{{Brand: "Chanel"}}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on empty results")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_Fetch_ServerError(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     http.StatusText(http.StatusInternalServerError),
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	params := []parameters.RequestPerfume{{Brand: "Chanel"}}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on server error")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_Fetch_Timeout(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		// Simulate timeout by waiting longer than the fetcher's timeout
		time.Sleep(3 * time.Second)
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader("{}")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewDB("http://test-url:8080", "test-token")
	params := []parameters.RequestPerfume{{Brand: "Chanel"}}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on timeout")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (errorReader) Close() error {
	return nil
}
