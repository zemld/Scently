package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPerfumesToUpdate_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/perfumes/update", bytes.NewBufferString("{invalid"))
	_, err := getPerfumesToUpdate(r)
	if err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}

func TestGetPerfumesToUpdate_ValidEmpty(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/perfumes/update", bytes.NewBufferString("{\"perfumes\":[]}"))
	perfumes, err := getPerfumesToUpdate(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(perfumes) != 0 {
		t.Fatalf("expected empty perfumes, got %d", len(perfumes))
	}
}
