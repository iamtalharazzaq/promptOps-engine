// Package config provides centralised configuration management for the
// PromptOps Engine backend. It reads values from environment variables (with
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
	"log/slog"
	"os"
	"strconv"

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
	MaxTokens   int    // Max tokens per response (0 = unlimited)
	DBURL       string // Supabase / PostgreSQL connection string
	JWTSecret   string // Secret key for signing JWT tokens
}

// Load reads environment variables (falling back to .env if present) and
// returns a fully-populated Config. Missing variables are replaced with
// sensible development defaults so the server can start with zero config.
func Load() *Config {
	// Best-effort: ignore error if .env doesn't exist (e.g. in Docker)
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		OllamaHost:  getEnv("OLLAMA_HOST", ""), // Require explicit config
		OllamaModel: getEnv("OLLAMA_MODEL", "tinyllama"),
		FrontendURL: getEnv("FRONTEND_URL", ""), // Require explicit config
		MaxTokens:   getEnvInt("MAX_TOKENS", 256),
		DBURL:       getEnv("DB_URL", ""),     // Require explicit config for Week 5
		JWTSecret:   getEnv("JWT_SECRET", "change-me-hackery-secret"),
	}
}

// Validate checks if mandatory configuration values are present.
func (c *Config) Validate() {
	if c.OllamaHost == "" {
		slog.Error("OLLAMA_HOST is required but not set in environment")
		os.Exit(1)
	}
	if c.FrontendURL == "" {
		slog.Error("FRONTEND_URL is required but not set in environment")
		os.Exit(1)
	}
	if c.DBURL == "" {
		slog.Error("DB_URL is required for Supabase persistence")
		os.Exit(1)
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

// getEnvInt returns the integer value of the named environment variable,
// or fallback if the variable is empty, unset, or not a valid integer.
func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	if i, err := strconv.Atoi(v); err == nil {
		return i
	}
	return fallback
}
