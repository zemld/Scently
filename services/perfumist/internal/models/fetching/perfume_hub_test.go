package fetching

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/config"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
	"github.com/zemld/Scently/perfumist/internal/models/perfume"
)

func TestDbFetcher_getPerfumes_Success(t *testing.T) {
	expectedPerfumes := []models.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "Sauvage", Sex: "male"},
	}
	resp := perfume.PerfumeResponse{Perfumes: expectedPerfumes}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
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
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_getPerfumes_BadStatus(t *testing.T) {
	resp := perfume.PerfumeResponse{Perfumes: []models.Perfume{}}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     http.StatusText(http.StatusInternalServerError),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_getPerfumes_ReadBodyError(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       errorReader{},
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
	perfumes, status := fetcher.getPerfumes(context.Background(), parameters.RequestPerfume{})

	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestDbFetcher_getPerfumes_NoContent(t *testing.T) {
	resp := perfume.PerfumeResponse{Perfumes: []models.Perfume{}}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
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
	resp := perfume.PerfumeResponse{Perfumes: []models.Perfume{}}
	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     http.StatusText(http.StatusNotFound),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
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
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Status:     http.StatusText(http.StatusForbidden),
			Body:       io.NopCloser(strings.NewReader("Forbidden")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
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
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		capturedRequest = r
		resp := perfume.PerfumeResponse{Perfumes: []models.Perfume{{Brand: "Test", Name: "Test"}}}
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
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
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

func TestDbFetcher_FetchMany_Success(t *testing.T) {
	callCount := 0
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		callCount++
		perfumes := []models.Perfume{
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
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
	params := []parameters.RequestPerfume{
		{Brand: "Chanel"},
		{Brand: "Dior"},
	}
	perfumesChan := fetcher.FetchMany(context.Background(), params)

	perfumes := make([]models.Perfume, 0)
	for p := range perfumesChan {
		perfumes = append(perfumes, p)
	}

	if callCount != 2 {
		t.Fatalf("expected 2 HTTP calls, got %d", callCount)
	}
	if len(perfumes) != 4 {
		t.Fatalf("expected 4 perfumes (2 from each request), got %d", len(perfumes))
	}
}

func TestDbFetcher_Fetch_EmptyResults(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp := perfume.PerfumeResponse{Perfumes: []models.Perfume{}}
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
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
	param := parameters.RequestPerfume{Brand: "Chanel"}
	perfumesChan := fetcher.Fetch(context.Background(), param)

	perfumes := make([]models.Perfume, 0)
	for p := range perfumesChan {
		perfumes = append(perfumes, p)
	}

	if len(perfumes) != 0 {
		t.Fatalf("expected 0 perfumes, got %d", len(perfumes))
	}
}

func TestDbFetcher_Fetch_ServerError(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     http.StatusText(http.StatusInternalServerError),
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
	param := parameters.RequestPerfume{Brand: "Chanel"}
	perfumesChan := fetcher.Fetch(context.Background(), param)

	perfumes := make([]models.Perfume, 0)
	for p := range perfumesChan {
		perfumes = append(perfumes, p)
	}

	if len(perfumes) != 0 {
		t.Fatalf("expected 0 perfumes on server error, got %d", len(perfumes))
	}
}

func TestDbFetcher_Fetch_Timeout(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
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
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewPerfumeHub("http://test-url:8080", "test-token", mockConfig)
	param := parameters.RequestPerfume{Brand: "Chanel"}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	perfumesChan := fetcher.Fetch(ctx, param)

	perfumes := make([]models.Perfume, 0)
	for p := range perfumesChan {
		perfumes = append(perfumes, p)
	}

	if len(perfumes) != 0 {
		t.Fatalf("expected 0 perfumes on timeout, got %d", len(perfumes))
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (errorReader) Close() error {
	return nil
}
