package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/promptops/backend/pkg/utils"
	"github.com/promptops/backend/services"
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

	// Schema is an optional JSON Schema (as a string) to validate the
	// LLM's response against.
	Schema string `json:"schema,omitempty"`
}

// ChatEvent is a single Server-Sent Event pushed to the client during
// streaming. The frontend uses the "done" flag to know when to finalize
// the assistant message bubble.
type ChatEvent struct {
	Content    string `json:"content"`              // Text fragment
	Done       bool   `json:"done"`                 // True on the final event
	Status     string `json:"status,omitempty"`     // "validating", "valid", "invalid", "retrying"
	RetryCount int    `json:"retry_count,omitempty"` // Current retry attempt
}

// -----------------------------------------------------------------------
// Handler
// -----------------------------------------------------------------------

// ChatHandler returns an http.HandlerFunc that:
//  1. Reads a ChatRequest from the request body.
//  2. Calls the Ollama service to generate a streaming response.
//  3. If a schema is provided, validates the response and retries on failure.
//  4. Pushes each token (or the final valid JSON) as an SSE event.
//
// Parameters:
//   - ollamaClient: the Ollama HTTP client to delegate generation to
//   - validator: the JSON Schema validation service
//   - defaultModel: the fallback model name if not specified in the request
func ChatHandler(ollamaClient *services.OllamaClient, validator *services.JSONValidator, defaultModel string, maxTokens int) http.HandlerFunc {
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
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, `{"error":"streaming not supported"}`, http.StatusInternalServerError)
			return
		}

		// ── Schema Guard Loop ──────────────────────────────────────
		log.Printf("[chat] Request with model=%s schema=%t", model, req.Schema != "")

		format := ""
		if req.Schema != "" {
			format = "json"
		}

		maxRetries := 3
		retryCount := 0
		currentPrompt := req.Message

		for {
			var fullResponse strings.Builder
			
			// If we are retrying, notify the client
			if retryCount > 0 {
				utils.WriteSSEEvent(w, flusher, ChatEvent{Status: "retrying", RetryCount: retryCount})
			} else if req.Schema != "" {
				utils.WriteSSEEvent(w, flusher, ChatEvent{Status: "validating"})
			}

			err := ollamaClient.GenerateStream(model, currentPrompt, format, maxTokens, func(content string, done bool) {
				if req.Schema != "" {
					// In schema mode, we collect the response first to validate it
					fullResponse.WriteString(content)
				} else {
					// Normal mode: stream directly to client
					utils.WriteSSEEvent(w, flusher, ChatEvent{Content: content, Done: done})
				}
			})

			if err != nil {
				log.Printf("[chat] Ollama error: %v", err)
				utils.WriteSSEEvent(w, flusher, ChatEvent{Content: fmt.Sprintf("[error] %v", err), Done: true})
				return
			}

			// If no schema, we're done (streaming was already handled in callback)
			if req.Schema == "" {
				return
			}

			// ── Validate JSON Response ──────────────────────────────
			responseStr := fullResponse.String()
			log.Printf("[chat] Validating response (len=%d)", len(responseStr))
			
			validationErr := validator.Validate(responseStr, req.Schema)
			if validationErr == nil {
				// SUCCESS: Valid JSON
				log.Printf("[chat] Schema validation passed")
				utils.WriteSSEEvent(w, flusher, ChatEvent{
					Content: responseStr,
					Done:    true,
					Status:  "valid",
				})
				return
			}

			// FAILURE: Invalid JSON
			log.Printf("[chat] Validation failed: %v", validationErr)
			retryCount++
			if retryCount >= maxRetries {
				log.Printf("[chat] Max retries exhausted")
				utils.WriteSSEEvent(w, flusher, ChatEvent{
					Content: responseStr, // Show the last malformed output anyway
					Done:    true,
					Status:  "invalid",
				})
				return
			}

			// Prepare retry prompt
			currentPrompt = fmt.Sprintf("%s\n\nIMPORTANT: Your last response was invalid JSON. Error: %v. Please try again and return ONLY valid JSON matching the requested schema.", req.Message, validationErr)
		}
	}
}
