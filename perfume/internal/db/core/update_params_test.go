package core

import "testing"

func TestUpdateParameters_WithTruncate(t *testing.T) {
	p := NewUpdateParameters()
	if p.IsTruncate {
		t.Fatalf("IsTruncate default true, want false")
	}
	p.WithTruncate()
	if !p.IsTruncate {
		t.Fatalf("WithTruncate() did not set IsTruncate to true")
	}
}
