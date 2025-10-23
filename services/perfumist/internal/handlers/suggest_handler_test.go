package handlers

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
)

func TestParseQuery(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/?brand=A&name=X", nil)
	p := parseQuery(req)
	if p.Brand != "A" || p.Name != "X" || !isValidQuery(p) {
		t.Fatalf("parseQuery should succeed: %+v", p)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/?brand=A", nil)
	p2 := parseQuery(req2)
	if p2.Brand != "A" || p2.Name != "" || isValidQuery(p2) {
		t.Fatalf("parseQuery should fail when name missing: %+v", p2)
	}
}

func TestFillResponseWithSuggestions(t *testing.T) {
	t.Parallel()

	suggestions := []models.GluedPerfumeWithScore{
		{GluedPerfume: models.GluedPerfume{Brand: "A", Name: "X"}, Score: 0.789},
		{GluedPerfume: models.GluedPerfume{Brand: "B", Name: "Y"}, Score: 0.0},
	}
	var resp SuggestResponse
	fillResponseWithSuggestions(&resp, suggestions)
	if !resp.Success {
		t.Fatalf("success should be true when there is at least one suggestion")
	}
	if resp.Suggested[0].Rank != 1 {
		t.Fatalf("rank should start at 1")
	}
	if math.Abs(resp.Suggested[0].Score-0.79) > 1e-9 {
		t.Fatalf("score should be rounded to 2 decimals, got %v", resp.Suggested[0].Score)
	}

	var respEmpty SuggestResponse
	fillResponseWithSuggestions(&respEmpty, []models.GluedPerfumeWithScore{})
	if respEmpty.Success {
		t.Fatalf("empty suggestions should set success=false")
	}
}
