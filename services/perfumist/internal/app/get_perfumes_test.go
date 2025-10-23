package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
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
	ch <- perfumesFetchAndGlueResult{Status: http.StatusOK, Perfumes: []models.GluedPerfume{{Name: "A"}}}
	ch <- perfumesFetchAndGlueResult{Status: http.StatusOK, Perfumes: []models.GluedPerfume{{Name: "B"}, {Name: "C"}}}
	close(ch)

	perfumes, status := fetchPerfumeResults(httptest.NewRequest(http.MethodGet, "/", nil).Context(), ch)
	if status != 0 {
		t.Fatalf("expected status 0, got %d", status)
	}
	if len(perfumes) != 3 {
		t.Fatalf("expected 3 perfumes, got %d", len(perfumes))
	}
}
