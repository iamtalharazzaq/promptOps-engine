// Package middleware provides reusable HTTP middleware for the PromptOps Engine
// backend, implementing cross-cutting concerns such as request logging and
// CORS header injection.
//
// Middleware functions follow the standard Go http.Handler wrapper pattern
// and are designed to be composed via chi's Use() method or any compatible
// router/mux.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter is a thin wrapper around http.ResponseWriter that captures
// the HTTP status code written by downstream handlers. This lets the
// logging middleware report the final status without altering the response.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before delegating to the embedded
// ResponseWriter so it can be logged after the request completes.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Flush delegates to the underlying ResponseWriter's Flush method if it
// implements http.Flusher. This is critical for SSE streaming endpoints
// like /chat — without it, the type assertion `w.(http.Flusher)` in the
// chat handler would fail because the wrapper hides the interface.
func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Logger returns middleware that logs every incoming request with:
//   - HTTP method and path
//   - Response status code
//   - Round-trip latency
//
// Example output:
//
//	[HTTP] POST /chat → 200 (12.34ms)
//
// This is automatically applied to all routes via the router setup in main.go.
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the writer to capture the status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Serve the request
			next.ServeHTTP(wrapped, r)

			// Log request details
			log.Printf(
				"[HTTP] %s %s → %d (%s)",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				time.Since(start).Round(time.Microsecond),
			)
		})
	}
}
