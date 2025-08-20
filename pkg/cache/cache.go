package cache

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var ctx = context.Background()

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

    // Test connection
    if _, err := Rdb.Ping(ctx).Result(); err != nil {
        return err
    }

    return nil
}

// CacheSet stores a mapping in Redis with TTL
func CacheSet(shortKey, longURL string, ttl time.Duration) error {
    return Rdb.Set(ctx, "url:"+shortKey, longURL, ttl).Err()
}

// CacheGet retrieves a longURL by shortKey
func CacheGet(shortKey string) (string, error) {
    return Rdb.Get(ctx, "url:"+shortKey).Result()
}

// CacheIncrClicks increments click counter for shortKey
func CacheIncrClicks(shortKey string) (int64, error) {
    return Rdb.Incr(ctx, "clicks:"+shortKey).Result()
}
