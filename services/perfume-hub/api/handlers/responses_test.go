package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteResponse(t *testing.T) {
	res := httptest.NewRecorder()
	payload := map[string]string{"status": "ok"}

	WriteResponse(res, http.StatusCreated, payload)

	result := res.Result()
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d", result.StatusCode, http.StatusCreated)
	}

	var decoded map[string]string
	if err := json.NewDecoder(result.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if decoded["status"] != "ok" {
		t.Fatalf("decoded body = %#v", decoded)
	}
}

func TestWriteResponse_JSONAndStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	body := map[string]any{"ok": true}
	WriteResponse(rr, http.StatusTeapot, body)

	if rr.Code != http.StatusTeapot {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusTeapot)
	}

	var decoded map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("body is not valid JSON: %v", err)
	}
	if decoded["ok"] != true {
		t.Fatalf("unexpected body: %v", decoded)
	}
}
