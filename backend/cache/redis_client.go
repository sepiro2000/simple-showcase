package cache

import (
	"context"
	"encoding/json"
	"time"

	"backend/config"
	"backend/models"

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

// GetProductList retrieves the product list from Redis
// Returns (products, exists, error) where exists indicates if the key was found in Redis
func GetProductList(ctx context.Context, client *redis.Client) ([]models.Product, bool, error) {
	if client == nil {
		return nil, false, nil
	}

	key := getProductListKey()
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil // Key doesn't exist in Redis
	}
	if err != nil {
		return nil, false, err
	}

	var products []models.Product
	if err := json.Unmarshal([]byte(val), &products); err != nil {
		return nil, false, err
	}

	return products, true, nil
}

// SetProductList sets the product list in Redis
func SetProductList(ctx context.Context, client *redis.Client, products []models.Product) error {
	if client == nil {
		return nil
	}

	key := getProductListKey()
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}

	// Cache for 1 hour
	return client.Set(ctx, key, data, time.Hour).Err()
}

// InvalidateProductList removes the product list from Redis cache
func InvalidateProductList(ctx context.Context, client *redis.Client) error {
	if client == nil {
		return nil
	}

	key := getProductListKey()
	return client.Del(ctx, key).Err()
}

// getProductListKey returns the Redis key for product list
func getProductListKey() string {
	return "product:list"
}
