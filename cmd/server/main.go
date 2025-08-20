package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Fahim047/awesome-url-shortener/pkg/api"
	"github.com/Fahim047/awesome-url-shortener/pkg/cache"
	"github.com/Fahim047/awesome-url-shortener/pkg/db"
)

func main() {
	// Initialize the postgres database connection
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Pool.Close()

	// Initialize the Redis cache connection
	if err := cache.Connect(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer cache.Rdb.Close()

	
	mux := http.NewServeMux()

	mux.Handle("POST /api/v1/shorten", http.HandlerFunc(api.ShortenURLHandler))
	mux.Handle("GET /api/v1/analytics/{short_key}", http.HandlerFunc(api.AnalyticsHandler))
	mux.Handle("GET /", http.HandlerFunc(api.RedirectHandler))


	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for range ticker.C {
        if err := cache.SyncClicks(context.Background()); err != nil {
            log.Println("Error syncing clicks:", err)
        }
    }
	}()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}