package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Fahim047/awesome-url-shortener/pkg/cache"
	"github.com/Fahim047/awesome-url-shortener/pkg/db"
	"github.com/Fahim047/awesome-url-shortener/pkg/shortener"
)

type shortenRequest struct {
	LongURL     string     `json:"long_url"`
	CustomAlias string     `json:"custom_key,omitempty"`
	ExpireAt    *time.Time `json:"expire_at,omitempty"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

// POST /api/v1/shorten
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.LongURL == "" {
		http.Error(w, "long_url is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Handle custom alias or generate random key
	shortKey := req.CustomAlias
	if shortKey != "" {
		existingMapping, err := db.GetMapping(ctx, shortKey)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if existingMapping != nil {
			http.Error(w, "Custom alias already exists", http.StatusConflict)
			return
		}
	} else {
		key, err := shortener.GenerateShortKey() // implement Base62 or random key
		if err != nil {
			http.Error(w, "Failed to generate short key", http.StatusInternalServerError)
			return
		}
		shortKey = key
	}

	// 2. Save mapping to DB
	mapping := &db.URLMapping{
		ShortKey: shortKey,
		LongURL:  req.LongURL,
		ExpireAt: req.ExpireAt,
	}

	if err := db.CreateMapping(ctx, mapping); err != nil {
		http.Error(w, "Could not save mapping", http.StatusInternalServerError)
		return
	}

	// 4. Return short URL
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	resp := shortenResponse{
		ShortURL: fmt.Sprintf("%s/%s", baseURL, shortKey),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /:shortKey
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[1:]
	if shortKey == "" {
		http.Error(w, "short key required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Try Redis cache first
	longURL, err := cache.CacheGet(ctx, shortKey)
	if err == nil && longURL != "" {
		cache.CacheIncrClicks(ctx, shortKey) // increment click count in Redis
		http.Redirect(w, r, longURL, http.StatusFound)
		return
	}

	// 2. Fallback to DB
	m, err := db.GetMapping(ctx, shortKey)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if m == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// 3. Check expiration
	if m.ExpireAt != nil && m.ExpireAt.Before(time.Now()) {
		http.Error(w, "Link expired", http.StatusGone)
		return
	}

	// 4. Increment click count in Redis
	cache.CacheIncrClicks(ctx, shortKey)

	// 5. Cache long URL with TTL if exists
	if m.ExpireAt != nil {
		ttl := time.Until(*m.ExpireAt)
		if ttl > 0 {
			cache.CacheSet(ctx, shortKey, m.LongURL, ttl)
		}
	} else {
		cache.CacheSet(ctx, shortKey, m.LongURL, 24*time.Hour)
	}

	// 6. Redirect
	http.Redirect(w, r, m.LongURL, http.StatusFound)
}
