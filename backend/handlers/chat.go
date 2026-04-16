package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ghostai-lite/backend/services"
)

// -----------------------------------------------------------------------
// Request / Response types
// -----------------------------------------------------------------------

// ChatRequest is the JSON body expected by POST /chat.
type ChatRequest struct {
	// Message is the user's natural-language prompt.
	Message string `json:"message"`

	// Model overrides the default Ollama model. If empty, the server's
	// configured default model is used (see config.OllamaModel).
	Model string `json:"model,omitempty"`
}

// ChatEvent is a single Server-Sent Event pushed to the client during
// streaming. The frontend uses the "done" flag to know when to finalize
// the assistant message bubble.
type ChatEvent struct {
	Content string `json:"content"` // Text fragment (may be a single token)
	Done    bool   `json:"done"`    // True on the final event
}

// -----------------------------------------------------------------------
// Handler
// -----------------------------------------------------------------------

// ChatHandler returns an http.HandlerFunc that:
//  1. Reads a ChatRequest from the request body.
//  2. Calls the Ollama service to generate a streaming response.
//  3. Pushes each token as an SSE event to the client.
//
// The response uses Server-Sent Events (SSE) with Content-Type
// "text/event-stream". Each event has the format:
//
//	data: {"content":"token","done":false}\n\n
//
// The final event sets done=true. If Ollama is unreachable or returns
// an error, a JSON error response is returned with status 502.
//
// Parameters:
//   - ollamaClient: the Ollama HTTP client to delegate generation to
//   - defaultModel: the fallback model name if not specified in the request
func ChatHandler(ollamaClient *services.OllamaClient, defaultModel string, maxTokens int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ── Parse request ──────────────────────────────────────────
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}

		if req.Message == "" {
			http.Error(w, `{"error":"message is required"}`, http.StatusBadRequest)
			return
		}

		// Use default model if none specified
		model := req.Model
		if model == "" {
			model = defaultModel
		}

		// ── Set SSE headers ────────────────────────────────────────
		// These headers tell the browser (and any proxy) to keep the
		// connection open and not buffer the response.
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no") // Nginx compatibility

		// Ensure we can flush
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, `{"error":"streaming not supported"}`, http.StatusInternalServerError)
			return
		}

		// ── Stream from Ollama ─────────────────────────────────────
		log.Printf("[chat] Streaming with model=%s maxTokens=%d prompt=%q", model, maxTokens, req.Message)

		err := ollamaClient.GenerateStream(model, req.Message, maxTokens, func(content string, done bool) {
			event := ChatEvent{Content: content, Done: done}
			data, _ := json.Marshal(event)

			// SSE format: "data: <json>\n\n"
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		})

		if err != nil {
			log.Printf("[chat] Ollama error: %v", err)
			// If headers are already sent, we can only log — the client
			// will see a truncated stream. For fresh errors, send 502.
			event := ChatEvent{Content: fmt.Sprintf("[error] %v", err), Done: true}
			data, _ := json.Marshal(event)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
