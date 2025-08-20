package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fahim047/awesome-url-shortener/pkg/api"
	"github.com/Fahim047/awesome-url-shortener/pkg/cache"
	"github.com/Fahim047/awesome-url-shortener/pkg/db"
)

func main() {
	// 1. Connect to services
	if err := db.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", err)
	}
	defer db.Pool.Close()

	if err := cache.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	defer cache.Rdb.Close()

	// 2. Setup routes
	mux := http.NewServeMux()
	mux.Handle("POST /api/v1/shorten", http.HandlerFunc(api.ShortenURLHandler))
	mux.Handle("GET /api/v1/analytics/{short_key}", http.HandlerFunc(api.AnalyticsHandler))
	mux.Handle("GET /", http.HandlerFunc(api.RedirectHandler))

	// 3. Background tasks
	go startClickSync()

	// 4. Server startup
	go func() {
		log.Println("🚀 Server running on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 5. Wait for shutdown signal
	waitForShutdown()
}

// Runs sync every 30s
func startClickSync() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := cache.SyncClicks(context.Background()); err != nil {
			log.Println("Error syncing clicks:", err)
		}
	}
}

// Handles shutdown (Ctrl+C / docker stop)
func waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down...")

	if err := cache.SyncClicks(context.Background()); err != nil {
		log.Println("❌ Error syncing clicks on shutdown:", err)
	} else {
		log.Println("✅ Final sync complete")
	}
}
