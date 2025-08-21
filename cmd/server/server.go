package main

import (
	"log"
	"net/http"
)

func StartServer(router http.Handler) {
	log.Println("🚀 Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
