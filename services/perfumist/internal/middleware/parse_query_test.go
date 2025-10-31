package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

func TestParseAndValidateQuery_ValidParams_SetsContextAndCallsNext(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		val := r.Context().Value(parameters.ParamsKey)
		params, ok := val.(parameters.RequestPerfume)
		if !ok {
			t.Fatalf("expected parameters.RequestPerfume in context, got %T", val)
		}
		if params.Brand != "Chanel" || params.Name != "No5" {
			t.Fatalf("unexpected brand/name in params: %+v", params)
		}
		if params.Sex != parameters.SexMale {
			t.Fatalf("expected sex %q, got %q", parameters.SexMale, params.Sex)
		}
		if params.UseAI != true {
			t.Fatalf("expected UseAI=true, got %v", params.UseAI)
		}
		w.WriteHeader(http.StatusOK)
	})

	h := ParseAndValidateQuery(next)
	req := httptest.NewRequest(http.MethodGet, "/?brand=Chanel&name=No5&sex=male&use_ai=true", nil)
	rr := httptest.NewRecorder()

	h(rr, req)

	if !nextCalled {
		t.Fatalf("expected next to be called")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestParseAndValidateQuery_DefaultsAndParsing(t *testing.T) {
	cases := []struct {
		name      string
		query     string
		expectSex string
		expectAI  bool
	}{
		{
			name:      "female sex, no use_ai defaults false",
			query:     "/?brand=Dior&name=Sauvage&sex=female",
			expectSex: parameters.SexFemale,
			expectAI:  false,
		},
		{
			name:      "invalid sex defaults to unisex",
			query:     "/?brand=Dior&name=Sauvage&sex=unknown&use_ai=true",
			expectSex: parameters.SexUnisex,
			expectAI:  true,
		},
		{
			name:      "invalid use_ai parses to false",
			query:     "/?brand=YSL&name=Libre&sex=male&use_ai=notabool",
			expectSex: parameters.SexMale,
			expectAI:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				val := r.Context().Value(parameters.ParamsKey)
				params, ok := val.(parameters.RequestPerfume)
				if !ok {
					t.Fatalf("expected parameters.RequestPerfume in context, got %T", val)
				}
				if params.Sex != tc.expectSex {
					t.Fatalf("expected sex %q, got %q", tc.expectSex, params.Sex)
				}
				if params.UseAI != tc.expectAI {
					t.Fatalf("expected UseAI=%v, got %v", tc.expectAI, params.UseAI)
				}
				w.WriteHeader(http.StatusOK)
			})

			h := ParseAndValidateQuery(next)
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			rr := httptest.NewRecorder()

			h(rr, req)

			if !nextCalled {
				t.Fatalf("expected next to be called")
			}
			if rr.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
			}
		})
	}
}

func TestParseAndValidateQuery_MissingRequiredParams_ReturnsBadRequest(t *testing.T) {
	cases := []struct {
		name string
		path string
	}{
		{name: "missing brand", path: "/?name=No5&sex=male"},
		{name: "missing name", path: "/?brand=Chanel&sex=female"},
		{name: "missing both", path: "/?sex=male"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			h := ParseAndValidateQuery(next)
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()

			h(rr, req)

			if nextCalled {
				t.Fatalf("did not expect next to be called")
			}
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
			}
			expectedBody := "Brand and name are required\n"
			if rr.Body.String() != expectedBody {
				t.Fatalf("expected body %q, got %q", expectedBody, rr.Body.String())
			}
		})
	}
}
