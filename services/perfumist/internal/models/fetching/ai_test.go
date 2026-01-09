package fetching

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/config"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
)

func TestNewAIFetcher(t *testing.T) {
	t.Parallel()

	url := "http://test-url:8000/v1/advise"
	mockConfig := &config.MockConfigManager{}
	fetcher := NewAI(url, "test-folder", "test-model", "test-api-key", mockConfig)

	if fetcher == nil {
		t.Fatal("expected non-nil fetcher")
	}
	if fetcher.url != url {
		t.Fatalf("expected url %q, got %q", url, fetcher.url)
	}
}

func TestAIFetcher_FetchMany_EmptyParams(t *testing.T) {
	t.Parallel()

	mockConfig := &config.MockConfigManager{}
	fetcher := NewAI("http://test-url:8000/v1/advise", "test-folder", "test-model", "test-api-key", mockConfig)
	perfumesChan := fetcher.FetchMany(context.Background(), []parameters.RequestPerfume{})

	// Channel should be closed immediately for empty params
	_, ok := <-perfumesChan
	if ok {
		t.Fatal("expected closed channel on empty params")
	}
}

func TestAIFetcher_Fetch_HTTPError(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewAI("http://test-url:8000/v1/advise", "test-folder", "test-model", "test-api-key", mockConfig)
	param := parameters.RequestPerfume{Brand: "Chanel"}
	perfumesChan := fetcher.Fetch(context.Background(), param)

	// Channel should be closed with no perfumes on error
	_, ok := <-perfumesChan
	if ok {
		t.Fatal("expected closed channel on HTTP error")
	}
}

func TestAIFetcher_Fetch_Non200Status(t *testing.T) {
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
	fetcher := NewAI("http://test-url:8000/v1/advise", "test-folder", "test-model", "test-api-key", mockConfig)
	param := parameters.RequestPerfume{Brand: "Chanel"}
	perfumesChan := fetcher.Fetch(context.Background(), param)

	// Channel should be closed with no perfumes on non-200 status
	_, ok := <-perfumesChan
	if ok {
		t.Fatal("expected closed channel on non-200 status")
	}
}

func TestAIFetcher_Fetch_EmptyBody(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewAI("http://test-url:8000/v1/advise", "test-folder", "test-model", "test-api-key", mockConfig)
	param := parameters.RequestPerfume{Brand: "Chanel"}
	perfumesChan := fetcher.Fetch(context.Background(), param)

	// Channel should be closed with no perfumes on empty body
	_, ok := <-perfumesChan
	if ok {
		t.Fatal("expected closed channel on empty body")
	}
}

func TestAIFetcher_Fetch_InvalidJSON(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader("invalid json")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	mockConfig := &config.MockConfigManager{}
	fetcher := NewAI("http://test-url:8000/v1/advise", "test-folder", "test-model", "test-api-key", mockConfig)
	param := parameters.RequestPerfume{Brand: "Chanel"}
	perfumesChan := fetcher.Fetch(context.Background(), param)

	// Channel should be closed with no perfumes on invalid JSON
	_, ok := <-perfumesChan
	if ok {
		t.Fatal("expected closed channel on invalid JSON")
	}
}

func TestAIFetcher_Fetch_Success(t *testing.T) {
	expectedPerfumes := []models.Perfume{
		{Brand: "Chanel", Name: "No5", Sex: "female"},
		{Brand: "Dior", Name: "Sauvage", Sex: "male"},
	}
	textBytes, err := json.Marshal(expectedPerfumes)
	if err != nil {
		t.Fatalf("failed to marshal expected perfumes: %v", err)
	}
	ycResponse := map[string]any{
		"result": map[string]any{
			"alternatives": []any{
				map[string]any{
					"message": map[string]any{
						"text": string(textBytes),
					},
				},
			},
		},
	}
	body, err := json.Marshal(ycResponse)
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
	fetcher := NewAI("http://test-url:8000/v1/advise", "test-folder", "test-model", "test-api-key", mockConfig)
	param := parameters.RequestPerfume{Brand: "Chanel"}
	perfumesChan := fetcher.Fetch(context.Background(), param)

	perfumes := make([]models.Perfume, 0)
	for perfume := range perfumesChan {
		perfumes = append(perfumes, perfume)
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

func TestAIFetcher_Fetch_BuildsPOSTRequestToCompletionAPI(t *testing.T) {
	var capturedRequest *http.Request
	var capturedBody []byte

	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		capturedRequest = r
		if r.Body != nil {
			b, _ := io.ReadAll(r.Body)
			capturedBody = b
		}

		ycResponse := map[string]any{
			"result": map[string]any{
				"alternatives": []any{
					map[string]any{
						"message": map[string]any{
							"text": "[]",
						},
					},
				},
			},
		}
		body, _ := json.Marshal(ycResponse)
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

	url := "http://test-url:8000/v1/advise"
	mockConfig := &config.MockConfigManager{}
	fetcher := NewAI(url, "test-folder", "aliceai-llm/latest", "test-api-key", mockConfig)
	param := parameters.RequestPerfume{
		Brand: "Chanel",
		Name:  "No5",
		Sex:   models.Female,
	}
	perfumesChan := fetcher.Fetch(context.Background(), param)
	// Drain the channel to ensure request is made
	for range perfumesChan {
	}

	if capturedRequest == nil {
		t.Fatal("expected request to be captured")
	}

	if capturedRequest.Method != http.MethodPost {
		t.Fatalf("expected method %q, got %q", http.MethodPost, capturedRequest.Method)
	}
	if capturedRequest.URL.String() != url {
		t.Fatalf("expected url %q, got %q", url, capturedRequest.URL.String())
	}

	if got := capturedRequest.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type %q, got %q", "application/json", got)
	}
	if got := capturedRequest.Header.Get("Authorization"); got != "Bearer test-api-key" {
		t.Fatalf("expected Authorization %q, got %q", "Bearer test-api-key", got)
	}

	if len(capturedBody) == 0 {
		t.Fatal("expected non-empty request body")
	}
	var rb requestBody
	if err := json.Unmarshal(capturedBody, &rb); err != nil {
		t.Fatalf("failed to unmarshal request body: %v", err)
	}
	if rb.ModelUri != "gpt://test-folder/aliceai-llm/latest" {
		t.Fatalf("expected modelUri %q, got %q", "gpt://test-folder/aliceai-llm/latest", rb.ModelUri)
	}
	if len(rb.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(rb.Messages))
	}
	if rb.Messages[0].Role != "system" || rb.Messages[1].Role != "user" {
		t.Fatalf("unexpected message roles: %+v", rb.Messages)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
