// Package config provides centralised configuration management for the
// GhostAI Lite backend. It reads values from environment variables (with
// optional .env file support via godotenv) and exposes a single Config
// struct consumed by the rest of the application.
//
// Environment variables:
//
//	PORT          – HTTP server listen port        (default "8080")
//	OLLAMA_HOST   – Ollama API base URL            (default "http://localhost:11434")
//	OLLAMA_MODEL  – Model name for inference       (default "tinyllama")
//	FRONTEND_URL  – Allowed CORS origin            (default "http://localhost:3000")
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration values for the backend server.
// Fields are populated once at startup via Load() and should be treated
// as read-only for the lifetime of the process.
type Config struct {
	Port        string // HTTP listen port, e.g. "8080"
	OllamaHost  string // Base URL of the Ollama API server
	OllamaModel string // Default model name used for inference
	FrontendURL string // Allowed CORS origin for the frontend
}

// Load reads environment variables (falling back to .env if present) and
// returns a fully-populated Config. Missing variables are replaced with
// sensible development defaults so the server can start with zero config.
func Load() *Config {
	// Best-effort: ignore error if .env doesn't exist (e.g. in Docker)
	if err := godotenv.Load(); err != nil {
		log.Println("[config] No .env file found, using environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		OllamaHost:  getEnv("OLLAMA_HOST", "http://localhost:11434"),
		OllamaModel: getEnv("OLLAMA_MODEL", "tinyllama"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

// getEnv returns the value of the named environment variable, or fallback
// if the variable is empty or unset.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
