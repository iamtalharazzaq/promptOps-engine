package middleware

import "net/http"

// CORS returns middleware that injects Cross-Origin Resource Sharing headers
// on every response. This is required because the Next.js frontend
// (http://localhost:3000 in dev) makes cross-origin requests to the Go
// backend (http://localhost:8080).
//
// Headers set:
//   - Access-Control-Allow-Origin      – restricted to the configured frontend URL
//   - Access-Control-Allow-Methods     – POST, GET, OPTIONS
//   - Access-Control-Allow-Headers     – Content-Type
//   - Access-Control-Allow-Credentials – true  (for future session cookie support)
//
// Preflight (OPTIONS) requests are short-circuited with 204 No Content.
//
// Parameters:
//   - allowedOrigin: the frontend URL to permit, e.g. "http://localhost:3000"
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight request without touching downstream handlers
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
