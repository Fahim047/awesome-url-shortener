package cache

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

// Connect initializes Redis client
func Connect() error {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(redisDBStr)
	if err != nil {
		db = 0
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       db,
	})

	return Rdb.Ping(context.Background()).Err()
}

// CacheSet stores a mapping in Redis with TTL
func CacheSet(ctx context.Context, shortKey, longURL string, ttl time.Duration) error {
	return Rdb.Set(ctx, "url:"+shortKey, longURL, ttl).Err()
}

// CacheGet retrieves a longURL by shortKey
func CacheGet(ctx context.Context, shortKey string) (string, error) {
	return Rdb.Get(ctx, "url:"+shortKey).Result()
}

// CacheIncrClicks increments click counter for shortKey
func CacheIncrClicks(ctx context.Context, shortKey string) (int64, error) {
	return Rdb.Incr(ctx, "clicks:"+shortKey).Result()
}
