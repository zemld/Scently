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
	redisHost         = "redis_cache"
	redisPort         = "6379"
	redisPasswordFile = "REDIS_PASSWORD"
)

type PerfumeCacheKey struct {
	Brand      string
	Name       string
	AdviseType string
	Sex        string
}

func GetRedisClient() *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: getRedisPassword(),
		})
	})
	return client
}

func getRedisPassword() string {
	filename := os.Getenv(redisPasswordFile)
	if filename == "" {
		return ""
	}
	password, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(password)
}
