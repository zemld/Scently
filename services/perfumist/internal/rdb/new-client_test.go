package rdb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRedisPassword(t *testing.T) {
	t.Setenv("REDIS_PASSWORD", "")
	if pwd := getRedisPassword(); pwd != "" {
		t.Fatalf("expected empty password when env unset, got %q", pwd)
	}

	t.Setenv("REDIS_PASSWORD", "/non/existent/file")
	if pwd := getRedisPassword(); pwd != "" {
		t.Fatalf("expected empty password for missing file, got %q", pwd)
	}

	dir := t.TempDir()
	file := filepath.Join(dir, "pwd.txt")
	if err := os.WriteFile(file, []byte("secret"), 0600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	t.Setenv("REDIS_PASSWORD", file)
	if pwd := getRedisPassword(); pwd != "secret" {
		t.Fatalf("expected 'secret', got %q", pwd)
	}
}

func TestGetRedisClient_ReusesInstance(t *testing.T) {
	t.Parallel()

	client = nil

	c1 := GetRedisClient()
	c2 := GetRedisClient()
	if c1 != c2 {
		t.Fatalf("expected same client instance")
	}
}
