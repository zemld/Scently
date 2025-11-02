package rdb

import (
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client = nil
	once   sync.Once
)

const (
	redisHost        = "redis_cache"
	redisPort        = "6379"
	redisPasswordEnv = "REDIS_PASSWORD"
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: os.Getenv(redisPasswordEnv),
		})
	})
	return client
}
