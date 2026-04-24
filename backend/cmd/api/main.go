// PromptOps Engine Backend
//
// This is the entry point for the PromptOps Engine API server. It wires
// together configuration, middleware, services, and HTTP handlers into
// a single chi router and starts listening for requests.
//
// Architecture overview:
//
//	main.go          ← you are here (wiring + startup)
//	├─ config/       ← environment variable loading
//	├─ middleware/    ← CORS, request logging
//	├─ handlers/     ← HTTP handler functions (health, chat)
//	└─ services/     ← external service clients (Ollama)
//
// Quick start:
//
//	cd backend && go run main.go
//
// The server listens on the port defined by the PORT environment variable
// (default 8080). Make sure Ollama is running locally before hitting /chat.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/promptops/backend/config"
	"github.com/promptops/backend/handlers"
	"github.com/promptops/backend/middleware"
	"github.com/promptops/backend/pkg/db"
	"github.com/promptops/backend/services"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// ── Load configuration ──────────────────────────────────────
	cfg := config.Load()
	cfg.Validate()

	// ── Configure Structured Logging ────────────────────────────
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("PromptOps Engine Backend starting",
		"version", "v0.1.0",
		"ollama_host", cfg.OllamaHost,
		"ollama_model", cfg.OllamaModel,
		"max_tokens", cfg.MaxTokens,
	)

	// ── Initialise database ─────────────────────────────────────
	database, err := db.Init(context.Background(), cfg.DBURL)
	if err != nil {
		slog.Error("Database initialisation failed", "error", err)
		os.Exit(1)
	}

	// ── Initialise services ─────────────────────────────────────
	ollamaClient := services.NewOllamaClient(cfg.OllamaHost)
	jsonValidator := services.NewJSONValidator()

	// ── Build router ────────────────────────────────────────────
	r := chi.NewRouter()

	// Global middleware (applied to ALL routes)
	r.Use(middleware.CORS(cfg.FrontendURL))
	r.Use(middleware.Logger())
	r.Use(middleware.Metrics())

	// Routes
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health", handlers.HealthHandler(cfg.MaxTokens))

	// Auth routes
	r.Post("/auth/register", handlers.RegisterHandler(database, cfg.JWTSecret))
	r.Post("/auth/login", handlers.LoginHandler(database, cfg.JWTSecret))

	// Chat routes (Protected)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(cfg.JWTSecret))
		r.Get("/chats", handlers.GetChatsHandler(database))
		r.Post("/chat", handlers.ChatHandler(database, ollamaClient, jsonValidator, cfg.OllamaModel, cfg.MaxTokens))
	})

	// ── Start server ────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("Server listening", "addr", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
