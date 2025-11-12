package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors_AllowedOrigin(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	handler := Cors(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://frontend:3000")

	rr := httptest.NewRecorder()
	handler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	if body := rr.Body.String(); body != "success" {
		t.Errorf("expected body 'success', got '%s'", body)
	}
}

func TestCors_DisallowedOrigin(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called for disallowed origin")
	})

	handler := Cors(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://malicious-site.com")

	rr := httptest.NewRecorder()
	handler(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, status)
	}
}

func TestCors_DisallowedOrigin_ReturnsStatusForbidden(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Cors(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://unauthorized.com")

	rr := httptest.NewRecorder()
	handler(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("expected status code %d for disallowed origin, got %d", http.StatusForbidden, status)
	}
}

func TestCors_DisallowedOrigin_ReturnsErrorMessage(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Cors(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://evil.com")

	rr := httptest.NewRecorder()
	handler(rr, req)

	expectedBody := "CORS not allowed\n"
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("expected body '%s', got '%s'", expectedBody, body)
	}
}
