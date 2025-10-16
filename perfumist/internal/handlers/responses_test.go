package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteResponse_ContentTypeAndBody(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	payload := map[string]string{"ok": "true"}
	WriteResponse(rr, payload, http.StatusTeapot)

	if rr.Code != http.StatusTeapot {
		t.Fatalf("status code = %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type = %q", ct)
	}

	var got map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got["ok"] != "true" {
		t.Fatalf("unexpected body: %v", got)
	}
}
