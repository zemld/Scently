package app

import (
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
)

func TestGetCacheKey(t *testing.T) {
	t.Parallel()

	key := getCacheKey(models.Perfume{Brand: "A", Name: "X"})
	if key != "A:X" {
		t.Fatalf("unexpected key: %q", key)
	}
}
