package rdb

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client = nil

const (
	redisHost         = "redis_cache"
	redisPort         = "6379"
	redisPasswordFile = "REDIS_PASSWORD"
)

func GetRedisClient() *redis.Client {
	if Client != nil {
		return Client
	}
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: getRedisPassword(),
	})
	return Client
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
