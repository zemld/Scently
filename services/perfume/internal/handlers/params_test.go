package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestGetSelectionParameters(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/perfumes/get?brand=Dior&name=Sauvage", nil)
	p := getSelectionParameters(r)
	if p.Brand != "Dior" || p.Name != "Sauvage" {
		t.Fatalf("unexpected params: %+v", p)
	}

	r2 := httptest.NewRequest(http.MethodGet, "/v1/perfumes/get?brand=Chanel", nil)
	p2 := getSelectionParameters(r2)
	if p2.Brand != "Chanel" || p2.Name != "" {
		t.Fatalf("unexpected params: %+v", p2)
	}

	r3 := httptest.NewRequest(http.MethodGet, "/v1/perfumes/get", nil)
	p3 := getSelectionParameters(r3)
	if p3.Brand != "" || p3.Name != "" {
		t.Fatalf("unexpected params: %+v", p3)
	}
}

func TestGetUpdateParametersFromRequest_HardAndPassword(t *testing.T) {
	dir := t.TempDir()
	passPath := filepath.Join(dir, "pass.txt")
	if err := os.WriteFile(passPath, []byte("secret\n"), 0600); err != nil {
		t.Fatalf("write pass file: %v", err)
	}
	t.Setenv("HARD_UPDATE_PASSWORD_FILE", passPath)

	q := make(url.Values)
	q.Set("hard", "true")
	q.Set("password", "secret")
	r := httptest.NewRequest(http.MethodPost, "/v1/perfumes/update?"+q.Encode(), nil)
	up := getUpdateParametersFromRequest(r)
	if !up.IsTruncate {
		t.Fatalf("expected IsTruncate=true when hard=true and password correct")
	}

	q2 := make(url.Values)
	q2.Set("hard", "true")
	q2.Set("password", "wrong")
	r2 := httptest.NewRequest(http.MethodPost, "/v1/perfumes/update?"+q2.Encode(), nil)
	up2 := getUpdateParametersFromRequest(r2)
	if up2.IsTruncate {
		t.Fatalf("expected IsTruncate=false when password wrong")
	}

	q3 := make(url.Values)
	q3.Set("hard", "notabool")
	r3 := httptest.NewRequest(http.MethodPost, "/v1/perfumes/update?"+q3.Encode(), nil)
	up3 := getUpdateParametersFromRequest(r3)
	if up3.IsTruncate {
		t.Fatalf("expected IsTruncate=false when hard invalid")
	}
}
