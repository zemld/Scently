package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

func TestParseQuerySuccess(t *testing.T) {
	requestBody := `{"perfumes":[{"brand":"Brand","name":"Name","sex":"female"}]}`
	req := httptest.NewRequest(http.MethodGet, "/?brand=Amouage&name=Search&sex=male&hard=true", strings.NewReader(requestBody))
	res := httptest.NewRecorder()

	nextCalled := false
	handler := ParseQuery(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		sp, ok := r.Context().Value(core.SelectParametersContextKey).(*core.SelectParameters)
		if !ok {
			t.Fatalf("select parameters missing in context")
		}
		if sp.Brand != "Amouage" || sp.Name != "Search" || sp.Sex != "male" {
			t.Fatalf("unexpected select parameters: %#v", sp)
		}

		up, ok := r.Context().Value(models.UpdateParametersContextKey).(*models.UpdateParameters)
		if !ok {
			t.Fatalf("update parameters missing in context")
		}
		if len(up.Perfumes) != 1 {
			t.Fatalf("expected 1 perfume, got %d", len(up.Perfumes))
		}
	})

	handler(res, req)

	if !nextCalled {
		t.Fatalf("next handler was not called")
	}
}

func TestParseQueryInvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader("not-json"))
	res := httptest.NewRecorder()

	var capturedUpdateParams *models.UpdateParameters
	nextCalled := false
	ParseQuery(func(_ http.ResponseWriter, r *http.Request) {
		nextCalled = true
		capturedUpdateParams = r.Context().Value(models.UpdateParametersContextKey).(*models.UpdateParameters)
	})(res, req)

	if !nextCalled {
		t.Fatalf("next handler was not called for invalid JSON body")
	}
	if capturedUpdateParams == nil {
		t.Fatalf("expected update parameters in context")
	}
	if len(capturedUpdateParams.Perfumes) != 0 {
		t.Fatalf("expected no perfumes parsed, got %d", len(capturedUpdateParams.Perfumes))
	}
}

func TestGetPerfumesToUpdateInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader("not-json"))
	content, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("failed to read request body: %v", err)
	}
	if string(content) != "not-json" {
		t.Fatalf("unexpected content: %q", content)
	}
	req.Body = io.NopCloser(bytes.NewReader(content))
	up := models.NewUpdateParameters()
	if err := getPerfumesToUpdate(*req, up); err != nil {
		t.Fatalf("unexpected error when parsing invalid JSON body: %v", err)
	}
	if len(up.Perfumes) != 0 {
		t.Fatalf("expected zero perfumes on invalid JSON, got %d", len(up.Perfumes))
	}
}
