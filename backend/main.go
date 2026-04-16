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
	"fmt"
	"log"
	"net/http"

	"github.com/promptops/backend/config"
	"github.com/promptops/backend/handlers"
	"github.com/promptops/backend/middleware"
	"github.com/promptops/backend/services"

	"github.com/go-chi/chi/v5"
)

func main() {
	// ── Load configuration ──────────────────────────────────────
	cfg := config.Load()
	cfg.Validate()

	log.Printf("[main] PromptOps Engine Backend v0.1.0")
	log.Printf("[main] Ollama host  : %s", cfg.OllamaHost)
	log.Printf("[main] Ollama model : %s", cfg.OllamaModel)
	log.Printf("[main] Max tokens   : %d", cfg.MaxTokens)
	log.Printf("[main] Frontend URL : %s", cfg.FrontendURL)

	// ── Initialise services ─────────────────────────────────────
	ollamaClient := services.NewOllamaClient(cfg.OllamaHost)

	// ── Build router ────────────────────────────────────────────
	r := chi.NewRouter()

	// Global middleware (applied to ALL routes)
	r.Use(middleware.CORS(cfg.FrontendURL))
	r.Use(middleware.Logger())

	// Routes
	r.Get("/health", handlers.HealthHandler(cfg.MaxTokens))
	r.Post("/chat", handlers.ChatHandler(ollamaClient, cfg.OllamaModel, cfg.MaxTokens))

	// ── Start server ────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("[main] Server listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("[main] Server failed: %v", err)
	}
}
