package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/rdb"
)

const cacheTTL = 3600 * time.Second

func LookupCache(ctx context.Context, requestedPerfume parameters.RequestPerfume) ([]models.RankedPerfumeWithProps, error) {
	key := getCacheKey(requestedPerfume)
	client := rdb.GetRedisClient()

	cached, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result []models.RankedPerfumeWithProps
	if err := json.Unmarshal([]byte(cached), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func Cache(ctx context.Context, requestedPerfume parameters.RequestPerfume, toCache []models.RankedPerfumeWithProps) error {
	encoded, err := json.Marshal(toCache)
	if err != nil {
		return err
	}

	client := rdb.GetRedisClient()
	key := getCacheKey(requestedPerfume)

	if err := client.Set(ctx, key, encoded, cacheTTL).Err(); err != nil {
		return err
	}
	return nil
}

func getCacheKey(perfume parameters.RequestPerfume) string {
	return fmt.Sprintf("%s:%s:%t:%s", perfume.Brand, perfume.Name, perfume.UseAI, perfume.Sex)
}
