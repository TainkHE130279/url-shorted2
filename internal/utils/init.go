package utils

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRedis(redisUrl string) *redis.Client {
	if redisUrl == "" {
		redisUrl = "redis://localhost:6379"
	}
	var opts, err = redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	// Tối ưu Redis connection cho high concurrency
	opts.PoolSize = 100    // Tăng pool size
	opts.MinIdleConns = 20 // Tăng min idle connections
	opts.MaxRetries = 3    // Tăng số lần retry
	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second
	opts.PoolTimeout = 4 * time.Second

	return redis.NewClient(opts)
}
