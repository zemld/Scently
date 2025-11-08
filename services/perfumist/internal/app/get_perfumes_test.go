package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

func TestFetchPerfumeResults_ServerErrorShortCircuit(t *testing.T) {
	t.Parallel()

	ch := make(chan perfumesFetchAndGlueResult, 2)
	ch <- perfumesFetchAndGlueResult{Status: http.StatusInternalServerError}
	close(ch)

	_, status := fetchPerfumeResults(httptest.NewRequest(http.MethodGet, "/", nil).Context(), ch)
	if status != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", status)
	}
}

func TestFetchPerfumeResults_AccumulatesOKResults(t *testing.T) {
	t.Parallel()

	ch := make(chan perfumesFetchAndGlueResult, 3)
	ch <- perfumesFetchAndGlueResult{Status: http.StatusOK, Perfumes: []models.Perfume{{Name: "A"}}}
	ch <- perfumesFetchAndGlueResult{Status: http.StatusOK, Perfumes: []models.Perfume{{Name: "B"}, {Name: "C"}}}
	close(ch)

	perfumes, status := fetchPerfumeResults(httptest.NewRequest(http.MethodGet, "/", nil).Context(), ch)
	if status != 0 {
		t.Fatalf("expected status 0, got %d", status)
	}
	if len(perfumes) != 3 {
		t.Fatalf("expected 3 perfumes, got %d", len(perfumes))
	}
}

func TestFetchPerfumeResults_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	perfumes, status := fetchPerfumeResults(ctx, make(chan perfumesFetchAndGlueResult))
	if status != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes on context cancel, got %v", perfumes)
	}
}

func TestAddQueryParameterAddsValue(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	addQueryParameter(req, "brand", "Chanel")

	if got := req.URL.Query().Get("brand"); got != "Chanel" {
		t.Fatalf("expected brand to be set, got %q", got)
	}
}

func TestAddQueryParameterSkipsEmptyValue(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	addQueryParameter(req, "brand", "")

	if got := req.URL.Query().Get("brand"); got != "" {
		t.Fatalf("expected empty brand value, got %q", got)
	}
}

func TestUpdateQuerySetsHeadersAndFilters(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	params := parameters.RequestPerfume{Brand: "Chanel", Name: "No5", Sex: parameters.SexFemale}
	t.Setenv(perfumeInternalTokenEnv, "secret-token")

	updateQuery(req, params)

	query := req.URL.Query()
	if got := query.Get("brand"); got != "Chanel" {
		t.Fatalf("expected brand to be %q, got %q", "Chanel", got)
	}
	if got := query.Get("name"); got != "No5" {
		t.Fatalf("expected name to be %q, got %q", "No5", got)
	}
	if got := query.Get("sex"); got != parameters.SexFemale {
		t.Fatalf("expected sex to be %q, got %q", parameters.SexFemale, got)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer secret-token" {
		t.Fatalf("expected authorization header, got %q", got)
	}
}

func TestUpdateQuerySkipsInvalidSex(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	params := parameters.RequestPerfume{Sex: parameters.SexUnisex}

	updateQuery(req, params)

	if got := req.URL.Query().Get("sex"); got != "" {
		t.Fatalf("expected no sex parameter, got %q", got)
	}
}

func TestGetPerfumesDoError(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("boom")
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	perfumes, status := getPerfumes(context.Background(), parameters.RequestPerfume{})
	if status != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestGetPerfumesServerErrorStatus(t *testing.T) {
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

	perfumes, status := getPerfumes(context.Background(), parameters.RequestPerfume{})
	if status != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", status)
	}
	if perfumes != nil {
		t.Fatalf("expected nil perfumes, got %v", perfumes)
	}
}

func TestGetPerfumesNoContent(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		respBody := `{"perfumes":[]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       io.NopCloser(strings.NewReader(respBody)),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	t.Cleanup(func() {
		http.DefaultClient.Transport = origTransport
	})

	perfumes, status := getPerfumes(context.Background(), parameters.RequestPerfume{})
	if status != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", status)
	}
	if perfumes == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(perfumes) != 0 {
		t.Fatalf("expected no perfumes, got %d", len(perfumes))
	}
}

func TestGetPerfumesSuccess(t *testing.T) {
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp := models.PerfumeResponse{
			Perfumes: []models.Perfume{{Name: "Test", Brand: "Brand"}},
		}
		body, err := json.Marshal(resp)
		if err != nil {
			return nil, err
		}
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

	perfumes, status := getPerfumes(context.Background(), parameters.RequestPerfume{})
	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}
	if len(perfumes) != 1 {
		t.Fatalf("expected 1 perfume, got %d", len(perfumes))
	}
	if perfumes[0].Name != "Test" || perfumes[0].Brand != "Brand" {
		t.Fatalf("unexpected perfume data: %+v", perfumes[0])
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
