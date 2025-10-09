package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/rdb"
)

const cacheTTL = 3600 * time.Second
const timeout = time.Second

func LookupCache(requestedPerfume models.Perfume) ([]models.RankedPerfumeWithProps, error) {
	key := getCacheKey(requestedPerfume)
	client := rdb.GetRedisClient()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cached, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var result []models.RankedPerfumeWithProps
	if err := json.Unmarshal([]byte(cached), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func Cache(requestedPerfume models.Perfume, toCache []models.RankedPerfumeWithProps) error {
	encoded, err := json.Marshal(toCache)
	if err != nil {
		return err
	}

	client := rdb.GetRedisClient()
	key := getCacheKey(requestedPerfume)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := client.Set(ctx, key, encoded, cacheTTL).Err(); err != nil {
		return err
	}
	return nil
}

func getCacheKey(perfume models.Perfume) string {
	return fmt.Sprintf("%s:%s", perfume.Brand, perfume.Name)
}
