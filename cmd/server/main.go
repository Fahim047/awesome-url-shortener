package main

import (
	"log"

	"github.com/Fahim047/awesome-url-shortener/pkg/cache"
	"github.com/Fahim047/awesome-url-shortener/pkg/db"
	"github.com/Fahim047/awesome-url-shortener/pkg/routes"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", err)
	}
	defer db.Pool.Close()

	if err := cache.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	defer cache.Rdb.Close()

	mux := routes.NewRouter()

	go startClickSync()

	go StartServer(mux)

	waitForShutdown()
}
