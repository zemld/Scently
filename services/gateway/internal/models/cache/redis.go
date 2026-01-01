package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

var (
	once   sync.Once
	client *redis.Client
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
	var initErr error
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:       fmt.Sprintf("%s:%s", host, port),
			Password:   password,
			MaxRetries: 1,
			MaintNotificationsConfig: &maintnotifications.Config{
				Mode: maintnotifications.ModeDisabled,
			},
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %w", err)
			client = nil
		}
	})

	if initErr != nil {
		return &RedisCacher{client: nil, cacheTTL: cacheTTL}, initErr
	}

	if client == nil {
		return &RedisCacher{client: nil, cacheTTL: cacheTTL}, fmt.Errorf("redis client is not initialized")
	}

	return &RedisCacher{
		client:   client,
		cacheTTL: cacheTTL,
	}, nil
}

func (c *RedisCacher) Save(ctx context.Context, key string, value []byte) error {
	if c.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	if err := c.client.Set(ctx, key, value, c.cacheTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (c *RedisCacher) Load(ctx context.Context, key string) ([]byte, error) {
	if c.client == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}
	encoded, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return []byte(encoded), nil
}
