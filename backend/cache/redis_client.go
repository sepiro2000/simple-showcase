package cache

import (
	"context"
	"fmt"
	"time"

	"backend/config"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis establishes a connection to Redis if configured
// Returns nil, nil if Redis is not configured (no error)
func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	if cfg.RedisHost == "" {
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

// GetProductLikes retrieves the like count for a product from Redis
func GetProductLikes(ctx context.Context, client *redis.Client, productID int64) (int64, error) {
	if client == nil {
		return 0, nil
	}

	key := getProductLikesKey(productID)
	val, err := client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// IncrementProductLikes increments the like count for a product in Redis
func IncrementProductLikes(ctx context.Context, client *redis.Client, productID int64) (int64, error) {
	if client == nil {
		return 0, nil
	}

	key := getProductLikesKey(productID)
	return client.Incr(ctx, key).Result()
}

// getProductLikesKey returns the Redis key for a product's likes
func getProductLikesKey(productID int64) string {
	return fmt.Sprintf("product:%d:likes", productID)
}
