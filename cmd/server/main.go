package main

import (
	"fmt"
	"log"
	"net/http"

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

	// Health check endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}