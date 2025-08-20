package main

import (
	"log"
	"net/http"

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
	mux.Handle("GET /", http.HandlerFunc(api.RedirectHandler))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}