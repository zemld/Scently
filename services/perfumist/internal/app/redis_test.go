package app

import (
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/rdb"
)

func TestGetCacheKey(t *testing.T) {
	t.Parallel()

	key := getCacheKey(rdb.PerfumeCacheKey{Brand: "A", Name: "X", AdviseType: "Comparision"})
	if key != "A:X" {
		t.Fatalf("unexpected key: %q", key)
	}
}
