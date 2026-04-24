// Package middleware provides reusable HTTP middleware for the PromptOps Engine
// backend, implementing cross-cutting concerns such as request logging and
// CORS header injection.
//
// Middleware functions follow the standard Go http.Handler wrapper pattern
// and are designed to be composed via chi's Use() method or any compatible
// router/mux.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// RequestIDKey is the key used to store the request ID in the context.
	RequestIDKey contextKey = "request_id"
	// UserIDKey is the key used to store the user ID in the context.
	UserIDKey contextKey = "user_id"
	// UserEmailKey is the key used to store the user email in the context.
	UserEmailKey contextKey = "user_email"
)

// responseWriter is a thin wrapper around http.ResponseWriter that captures
// the HTTP status code written by downstream handlers.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Logger returns middleware that logs every incoming request using structured logging.
// It also generates a unique Request ID for each request and injects it into the context.
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate or retrieve Request ID
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Add Request ID to context
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
			r = r.WithContext(ctx)

			// Set X-Request-ID header on response
			w.Header().Set("X-Request-ID", requestID)

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Serve the request
			next.ServeHTTP(wrapped, r)

			// Log request details using slog
			slog.Info("request completed",
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapped.statusCode),
				slog.Duration("latency", time.Since(start)),
				slog.String("ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
		})
	}
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
