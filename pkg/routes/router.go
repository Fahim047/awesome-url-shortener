package routes

import (
	"net/http"

	"github.com/Fahim047/awesome-url-shortener/pkg/api"
	"github.com/Fahim047/awesome-url-shortener/pkg/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /api/v1/shorten", http.HandlerFunc(api.ShortenURLHandler))
	mux.Handle("GET /api/v1/analytics/{short_key}", http.HandlerFunc(api.AnalyticsHandler))
	mux.Handle("GET /", http.HandlerFunc(api.RedirectHandler))

	return middleware.Chain(mux,
		middleware.Logger,
	)
}
