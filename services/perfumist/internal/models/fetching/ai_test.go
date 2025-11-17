package fetching

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

func TestNewAIFetcher(t *testing.T) {
	t.Parallel()

	url := "http://test-url:8000/v1/advise"
	fetcher := NewAI(url)

	if fetcher == nil {
		t.Fatal("expected non-nil fetcher")
	}
	if fetcher.url != url {
		t.Fatalf("expected url %q, got %q", url, fetcher.url)
	}
}

func TestAIFetcher_Fetch_EmptyParams(t *testing.T) {
	t.Parallel()

	fetcher := NewAI("http://test-url:8000/v1/advise")
	perfumes, ok := fetcher.Fetch(context.Background(), []parameters.RequestPerfume{})

	if ok {
		t.Fatal("expected false on empty params")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestAIFetcher_Fetch_HTTPError(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewAI("http://test-url:8000/v1/advise")
	params := []parameters.RequestPerfume{{Brand: "Chanel"}}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on HTTP error")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestAIFetcher_Fetch_Non200Status(t *testing.T) {
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

	fetcher := NewAI("http://test-url:8000/v1/advise")
	params := []parameters.RequestPerfume{{Brand: "Chanel"}}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on non-200 status")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestAIFetcher_Fetch_EmptyBody(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewAI("http://test-url:8000/v1/advise")
	params := []parameters.RequestPerfume{{Brand: "Chanel"}}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on empty body")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestAIFetcher_Fetch_InvalidJSON(t *testing.T) {
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader("invalid json")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		config.HTTPClient.Transport = origTransport
	})

	fetcher := NewAI("http://test-url:8000/v1/advise")
	params := []parameters.RequestPerfume{{Brand: "Chanel"}}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if ok {
		t.Fatal("expected false on invalid JSON")
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestAIFetcher_Fetch_Success(t *testing.T) {
	expectedPerfumes := []perfume.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "Sauvage", Sex: "male"},
	}
	suggestion := aISuggestion{Perfumes: expectedPerfumes}
	body, err := json.Marshal(suggestion)
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

	fetcher := NewAI("http://test-url:8000/v1/advise")
	params := []parameters.RequestPerfume{{Brand: "Chanel"}}
	perfumes, ok := fetcher.Fetch(context.Background(), params)

	if !ok {
		t.Fatal("expected true on success")
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

func TestAIFetcher_Fetch_AddsQueryParams(t *testing.T) {
	var capturedRequest *http.Request
	origTransport := config.HTTPClient.Transport
	config.HTTPClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		capturedRequest = r
		suggestion := aISuggestion{Perfumes: []perfume.Perfume{}}
		body, _ := json.Marshal(suggestion)
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

	fetcher := NewAI("http://test-url:8000/v1/advise")
	params := []parameters.RequestPerfume{
		{Brand: "Chanel", Name: "No5", Sex: parameters.SexFemale},
	}
	fetcher.Fetch(context.Background(), params)

	if capturedRequest == nil {
		t.Fatal("expected request to be captured")
	}
	query := capturedRequest.URL.Query()
	if got := query.Get("brand"); got != "Chanel" {
		t.Fatalf("expected brand %q, got %q", "Chanel", got)
	}
	if got := query.Get("name"); got != "No5" {
		t.Fatalf("expected name %q, got %q", "No5", got)
	}
	if got := query.Get("sex"); got != parameters.SexFemale {
		t.Fatalf("expected sex %q, got %q", parameters.SexFemale, got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
