package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	// Redis is the Redis client instance
	Redis *redis.Client
	// Ctx is the background context
	Ctx = context.Background()
)

// InitRedis initializes and tests the Redis connection
func InitRedis() error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // Default Redis address
	}

	// Parse pool size from environment or use default
	poolSize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	if poolSize <= 0 {
		poolSize = 10 // Default pool size
	}

	// Parse min idle connections
	minIdleConns, _ := strconv.Atoi(os.Getenv("REDIS_MIN_IDLE_CONN"))
	if minIdleConns <= 0 {
		minIdleConns = 3 // Default min idle connections
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     os.Getenv("REDIS_PASSWORD"), // Optional password
		DB:           0,                           // Use default DB
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		// Connection and operation timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
		// Reconnection settings
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})

	// Test the connection with retry
	var lastErr error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		_, lastErr = Redis.Ping(Ctx).Result()
		if lastErr == nil {
			log.Println("Successfully connected to Redis")
			return nil
		}

		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * time.Second
			log.Printf("Failed to connect to Redis (attempt %d/%d), retrying in %v: %v", 
				i+1, maxRetries, waitTime, lastErr)
			time.Sleep(waitTime)
		}
	}

	return fmt.Errorf("failed to connect to Redis after %d attempts: %v", maxRetries, lastErr)
}