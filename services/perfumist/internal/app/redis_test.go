package app

import (
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

func TestGetCacheKey(t *testing.T) {
	t.Parallel()

	key := getCacheKey(parameters.RequestPerfume{Brand: "A", Name: "X", UseAI: false, Sex: "male"})
	if key != "A:X:Comparision:male" {
		t.Fatalf("unexpected key: %q", key)
	}
}
