package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
)

var (
	once          sync.Once
	redisInstance *RedisCacher
)

type RedisCacher struct {
	client   *redis.Client
	cacheTTL time.Duration
}

type SuggestionsContextKey string

const (
	SuggestionsKey SuggestionsContextKey = "suggestions"
)

func GetOrCreateRedisCacher(host string, port string, password string, cacheTTL time.Duration) *RedisCacher {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
		})
		redisInstance = &RedisCacher{
			client:   client,
			cacheTTL: cacheTTL,
		}
	})
	return redisInstance
}

func (c *RedisCacher) Save(ctx context.Context, key string, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := c.client.Set(ctx, key, encoded, c.cacheTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (c *RedisCacher) Load(ctx context.Context, key string) (any, error) {
	encoded, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result perfume.Suggestions
	if err := json.Unmarshal([]byte(encoded), &result); err != nil {
		return nil, err
	}
	return result, nil
}
