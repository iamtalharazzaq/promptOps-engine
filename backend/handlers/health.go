// Package handlers provides HTTP handler functions for the GhostAI Lite API.
//
// Available endpoints:
//
//	GET  /health – Liveness / readiness probe (see HealthHandler)
//	POST /chat   – Stream an LLM response via SSE (see ChatHandler)
//
// Handlers are intentionally thin: they validate input, delegate to a
// service, and format the response. All heavy lifting lives in the
// services package.
package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse is the JSON body returned by the /health endpoint.
// It confirms the server is running and reports the current server time.
type HealthResponse struct {
	Status    string `json:"status"`    // Always "ok" when the server is reachable
	Timestamp string `json:"timestamp"` // ISO-8601 server time
	Service   string `json:"service"`   // Service identifier ("ghostai-backend")
	Version   string `json:"version"`   // Semantic version of the backend
	MaxTokens int    `json:"maxTokens"` // Configured per-response token limit
}

// HealthHandler returns an http.HandlerFunc that responds to GET /health
// with a JSON object confirming the service is alive. This endpoint is
// used by Docker health-checks, load balancers, and the frontend's
// connectivity indicator.
//
// Response example:
//
//	{
//	  "status":    "ok",
//	  "timestamp": "2026-04-16T14:00:00Z",
//	  "service":   "ghostai-backend",
//	  "version":   "0.1.0",
//	  "maxTokens": 256
//	}
func HealthHandler(maxTokens int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Service:   "ghostai-backend",
			Version:   "0.1.0",
			MaxTokens: maxTokens,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
