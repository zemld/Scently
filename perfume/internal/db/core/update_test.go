package core

import "testing"

func TestGetSavepointQuery(t *testing.T) {
	if got := getSavepointQuery("SAVEPOINT sp_", 3); got != "SAVEPOINT sp_3" {
		t.Fatalf("getSavepointQuery() = %q, want %q", got, "SAVEPOINT sp_3")
	}
	if got := getSavepointQuery("ROLLBACK TO SAVEPOINT sp_", 0); got != "ROLLBACK TO SAVEPOINT sp_0" {
		t.Fatalf("getSavepointQuery() = %q, want %q", got, "ROLLBACK TO SAVEPOINT sp_0")
	}
}
