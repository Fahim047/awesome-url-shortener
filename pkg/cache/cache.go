package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Fahim047/awesome-url-shortener/pkg/db"
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
func CacheIncrClicks(ctx context.Context, shortKey string) error {
	key := "clicks:" + shortKey

	// Check if Redis already has the counter
	_, err := Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist → load from DB
		mapping, dbErr := db.GetMapping(ctx, shortKey)
		if dbErr != nil {
			return dbErr
		}
		if mapping == nil {
			return fmt.Errorf("mapping not found for %s", shortKey)
		}

		// Initialize Redis with DB click count
		if setErr := Rdb.Set(ctx, key, mapping.ClickCount, 0).Err(); setErr != nil {
			return setErr
		}
	} else if err != nil {
		return err
	}

	// Always increment after ensuring base is correct
	return Rdb.Incr(ctx, key).Err()
}

// SyncClicks updates click counts from Redis to the database.
// This function should be called periodically to sync click counts
// from Redis to the database.
func SyncClicks(ctx context.Context) error {
	log.Printf(">>> SyncClicks triggered at %s\n", time.Now().Format(time.RFC3339))
	keys, err := Rdb.Keys(ctx, "clicks:*").Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		shortKey := key[len("clicks:"):]
		countStr, err := Rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			continue
		}
		db.UpdateClickCount(ctx, shortKey, count)
	}
	return nil

}
