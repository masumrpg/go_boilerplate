package database

import (
	"context"
	"fmt"
	"time"

	"go_boilerplate/internal/shared/config"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// InitRedis initializes the Redis client
func InitRedis(cfg *config.Config, logger *logrus.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("✓ Connected to Redis")
	return rdb, nil
}
