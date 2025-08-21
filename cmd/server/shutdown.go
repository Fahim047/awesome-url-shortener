package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fahim047/awesome-url-shortener/pkg/cache"
)

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