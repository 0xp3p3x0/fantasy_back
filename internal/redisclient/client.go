package redisclient

import (
	"context"
	"fmt"

	"back/internal/config"

	"github.com/redis/go-redis/v9"
)

// New connects to Redis when cfg.RedisAddr is non-empty. Returns (nil, nil) if Redis is disabled.
func New(cfg *config.Config) (*redis.Client, error) {
	if cfg.RedisAddr == "" {
		return nil, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return rdb, nil
}
