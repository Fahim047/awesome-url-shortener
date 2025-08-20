package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// Connect initializes the database connection pool
func Connect() error {
	dbHost := getenv("DB_HOST", "localhost")
	dbPort := getenv("DB_PORT", "5432")
	dbUser := getenv("DB_USER", "postgres")
	dbPass := getenv("DB_PASSWORD", "postgres")
	dbName := getenv("DB_NAME", "awesome_url_shortener")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("DB unreachable: %w", err)
	}

	Pool = pool
	return nil
}

// helper function for env vars with default values
func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

type URLMapping struct {
	ID         int64
	ShortKey   string
	LongURL    string
	ClickCount int64
	CreatedAt  time.Time
	ExpireAt   *time.Time
}

func CreateMapping(ctx context.Context, mapping *URLMapping) error {
	query := `
		INSERT INTO url_mappings (short_key, long_url, expire_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return Pool.QueryRow(ctx, query, mapping.ShortKey, mapping.LongURL, mapping.ExpireAt).
		Scan(&mapping.ID, &mapping.CreatedAt)
}

func GetMapping(ctx context.Context, shortKey string) (*URLMapping, error) {
	query := `
		SELECT id, short_key, long_url, click_count, created_at, expire_at
		FROM url_mappings
		WHERE short_key=$1
	`
	row := Pool.QueryRow(ctx, query, shortKey)
	var m URLMapping
	if err := row.Scan(&m.ID, &m.ShortKey, &m.LongURL, &m.ClickCount, &m.CreatedAt, &m.ExpireAt); err != nil {
		return nil, errors.New("mapping not found")
	}
	return &m, nil
}

func IncrementClickCount(ctx context.Context, shortKey string) error {
	query := `
		UPDATE url_mappings
		SET click_count = click_count + 1
		WHERE short_key=$1
	`
	_, err := Pool.Exec(ctx, query, shortKey)
	return err
}

