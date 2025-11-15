package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacher struct {
	client   *redis.Client
	cacheTTL time.Duration
}

type SuggestionsContextKey string

const (
	SuggestionsKey SuggestionsContextKey = "suggestions"
)

func NewRedisCacher(host string, port string, password string, cacheTTL time.Duration) (*RedisCacher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisCacher{
		client:   client,
		cacheTTL: cacheTTL,
	}, nil
}

func (c *RedisCacher) Close() error {
	return c.client.Close()
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
