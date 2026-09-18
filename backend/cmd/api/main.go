package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/foolxdev/url-shortener/internal/handler"
	"github.com/foolxdev/url-shortener/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	//nolint:staticcheck // chi middleware.RealIP is intentionally used.
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(securityHeaders)
	if origins := allowedOrigins(); len(origins) > 0 {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: origins,
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			MaxAge:         300,
		}))
	}

	db, err := repository.OpenDB(
		env("DATABASE_URL", "postgres://web:web@localhost:5432/urlshortener?sslmode=disable"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := handler.NewApplication(db)

	// Create application/database model

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{
			"response": "ok",
		})
	})
	// Pass the initialized model to the handler

	r.Post("/api/urls", handler.CreateShortURL(app))
	r.Get("/{hash}", handler.RedirectURL(app))

	addr := env("APP_ADDR", ":4000")
	server := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("Server running on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func allowedOrigins() []string {
	value := os.Getenv("CORS_ALLOWED_ORIGINS")
	if value == "" {
		return []string{"http://localhost:5173"}
	}
	var origins []string
	for _, origin := range strings.Split(value, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
