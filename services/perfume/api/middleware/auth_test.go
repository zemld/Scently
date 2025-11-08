package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthAllowsRequestWithValidToken(t *testing.T) {
	previous := perfumeToken
	defer func() { perfumeToken = previous }()
	perfumeToken = "secret"

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	res := httptest.NewRecorder()

	nextCalled := false
	Auth(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})(res, req)

	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
	if res.Result().StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Result().StatusCode, http.StatusOK)
	}
}

func TestAuthRejectsRequestWithoutToken(t *testing.T) {
	previous := perfumeToken
	defer func() { perfumeToken = previous }()
	perfumeToken = "secret"

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	nextCalled := false
	Auth(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})(res, req)

	if nextCalled {
		t.Fatalf("next handler should not be called")
	}
	if res.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Result().StatusCode, http.StatusUnauthorized)
	}
}
