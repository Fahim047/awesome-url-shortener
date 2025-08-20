package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load DB config from env vars
	dbHost := getenv("DB_HOST", "localhost")
	dbPort := getenv("DB_PORT", "5432")
	dbUser := getenv("DB_USER", "postgres")
	dbPass := getenv("DB_PASSWORD", "postgres")
	dbName := getenv("DB_NAME", "awesome_url_shortener")

	// Build connection string
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	// Connect using pgx connection pool
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}
	defer pool.Close()

	// Test query
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("DB unreachable:", err)
	}

	// Health check endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// helper function for env vars with defaults
func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
