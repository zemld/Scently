package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteResponse_JSONAndStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	body := map[string]any{"ok": true}
	WriteResponse(rr, http.StatusTeapot, body)

	if rr.Code != http.StatusTeapot {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusTeapot)
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}

	var decoded map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("body is not valid JSON: %v", err)
	}
	if decoded["ok"] != true {
		t.Fatalf("unexpected body: %v", decoded)
	}
}
