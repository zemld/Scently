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
	var resp SuggestResponse
	p, ok := parseQuery(req, &resp)
	if !ok || p.Brand != "A" || p.Name != "X" || !resp.Input.Ok {
		t.Fatalf("parseQuery should succeed: %+v, ok=%v", p, ok)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/?brand=A", nil)
	var resp2 SuggestResponse
	if _, ok := parseQuery(req2, &resp2); ok || resp2.Input.Ok || resp2.Success {
		t.Fatalf("parseQuery should fail when name missing: %+v", resp2)
	}
}

func TestFillResponseWithSuggestions(t *testing.T) {
	t.Parallel()

	suggestions := []gluedPerfumeWithScore{
		{GluedPerfume: models.GluedPerfume{Brand: "A", Name: "X"}, Score: 0.789},
		{GluedPerfume: models.GluedPerfume{Brand: "B", Name: "Y"}, Score: 0.0},
	}
	var resp SuggestResponse
	fillResponseWithSuggestions(&resp, suggestions)
	if !resp.Success {
		t.Fatalf("success should be true when there is at least one suggestion")
	}
	if len(resp.Suggested) != 1 {
		t.Fatalf("should stop at score==0, got %d", len(resp.Suggested))
	}
	if resp.Suggested[0].Rank != 1 {
		t.Fatalf("rank should start at 1")
	}
	if math.Abs(resp.Suggested[0].Score-0.79) > 1e-9 {
		t.Fatalf("score should be rounded to 2 decimals, got %v", resp.Suggested[0].Score)
	}

	var respEmpty SuggestResponse
	fillResponseWithSuggestions(&respEmpty, []gluedPerfumeWithScore{})
	if respEmpty.Success {
		t.Fatalf("empty suggestions should set success=false")
	}
}

func TestUpdateMostSimilarIfNeeded(t *testing.T) {
	t.Parallel()

	arr := make([]gluedPerfumeWithScore, 4)
	for i, s := range []float64{0.1, 0.2, 0.3, 0.4} {
		updateMostSimilarIfNeeded(arr, models.GluedPerfume{Name: "n"}, s)
		if arr[i].Score == 0 {
			t.Fatalf("position %d should be filled", i)
		}
	}
	updateMostSimilarIfNeeded(arr, models.GluedPerfume{Name: "m"}, 0.25)
	if !(arr[1].Score >= 0.25 && arr[2].Score >= 0.25) {
		t.Fatalf("middle insertion not reflected: %+v", arr)
	}
}

func TestFetchPerfumeResults_StatusAggregationBug(t *testing.T) {
	t.Parallel()

	fav := make(chan perfumesFetchAndGlueResult, 1)
	all := make(chan perfumesFetchAndGlueResult, 1)

	fav <- perfumesFetchAndGlueResult{Status: http.StatusInternalServerError}
	close(fav)
	all <- perfumesFetchAndGlueResult{Status: http.StatusOK}
	close(all)

	_, _, status := fetchPerfumeResults(httptest.NewRequest(http.MethodGet, "/", nil).Context(), fav, all)
	if status != http.StatusInternalServerError {
		t.Fatalf("expected aggregated status 500, got %d", status)
	}
}
