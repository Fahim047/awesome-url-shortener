package main

import (
	"log"
	"net/http"
)

func StartServer(mux *http.ServeMux) {
	log.Println("🚀 Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
