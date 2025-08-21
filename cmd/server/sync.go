package main

import (
	"context"
	"log"
	"time"

	"github.com/Fahim047/awesome-url-shortener/pkg/cache"
)

func startClickSync() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := cache.SyncClicks(context.Background()); err != nil {
			log.Println("Error syncing clicks:", err)
		}
	}
}
