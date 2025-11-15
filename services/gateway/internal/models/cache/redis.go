package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
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

func (c *RedisCacher) Save(ctx context.Context, key string, value []byte) error {
	if err := c.client.Set(ctx, key, value, c.cacheTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (c *RedisCacher) Load(ctx context.Context, key string) ([]byte, error) {
	encoded, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return []byte(encoded), nil
}
