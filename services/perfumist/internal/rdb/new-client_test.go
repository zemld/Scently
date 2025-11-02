package rdb

import (
	"testing"
)

func TestGetRedisClient_ReusesInstance(t *testing.T) {
	t.Parallel()

	client = nil

	c1 := GetRedisClient()
	c2 := GetRedisClient()
	if c1 != c2 {
		t.Fatalf("expected same client instance")
	}
}
