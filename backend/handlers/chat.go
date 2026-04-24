package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/promptops/backend/middleware"
	"github.com/promptops/backend/pkg/models"
	"github.com/promptops/backend/pkg/utils"
	"github.com/promptops/backend/services"
	"github.com/uptrace/bun"
)

// -----------------------------------------------------------------------
// Request / Response types
// -----------------------------------------------------------------------

// ChatRequest is the JSON body expected by POST /chat.
type ChatRequest struct {
	// ChatID is the optional ID of an existing chat session.
	ChatID uuid.UUID `json:"chat_id,omitempty"`

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
	Content    string    `json:"content"`              // Text fragment
	Done       bool      `json:"done"`                 // True on the final event
	Status     string    `json:"status,omitempty"`     // "validating", "valid", "invalid", "retrying"
	RetryCount int       `json:"retry_count,omitempty"` // Current retry attempt
	ChatID     uuid.UUID `json:"chat_id,omitempty"`     // ID of the chat session
}

// -----------------------------------------------------------------------
// Handler
// -----------------------------------------------------------------------

// ChatHandler returns an http.HandlerFunc that:
//  1. Reads a ChatRequest from the request body.
//  2. Calls the Ollama service to generate a streaming response.
//  3. If a schema is provided, validates the response and retries on failure.
//  4. Pushes each token (or the final valid JSON) as an SSE event.
//  5. Saves the conversation to the database.
func ChatHandler(db *bun.DB, ollamaClient *services.OllamaClient, validator *services.JSONValidator, defaultModel string, maxTokens int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ── Extract User from Context ─────────────────────────────
		userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

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

		// ── Chat Session Setup ─────────────────────────────────────
		ctx := r.Context()
		requestID := middleware.GetRequestID(ctx)

		chat := &models.Chat{
			ID:     req.ChatID,
			UserID: userID,
			Model:  model,
		}

		// If new chat, set title and initial history
		if chat.ID == uuid.Nil {
			chat.ID = uuid.New()
			chat.Title = req.Message
			if len(chat.Title) > 30 {
				chat.Title = chat.Title[:27] + "..."
			}
		} else if db != nil {
			// Find existing chat history (simplified logic)
			_ = db.NewSelect().Model(chat).Where("id = ? AND user_id = ?", chat.ID, userID).Scan(ctx)
		}

		slog.Info("chat request initiated",
			"request_id", requestID,
			"chat_id", chat.ID,
			"user_id", userID,
			"has_schema", req.Schema != "",
		)

		format := ""
		if req.Schema != "" {
			format = "json"
		}

		maxRetries := 3
		retryCount := 0
		currentPrompt := req.Message
		var fullAssistantResponse strings.Builder

		for {
			var currentIterationResponse strings.Builder

			// If we are retrying, notify the client
			if retryCount > 0 {
				utils.WriteSSEEvent(w, flusher, ChatEvent{Status: "retrying", RetryCount: retryCount, ChatID: chat.ID})
			} else if req.Schema != "" {
				utils.WriteSSEEvent(w, flusher, ChatEvent{Status: "validating", ChatID: chat.ID})
			}

			err := ollamaClient.GenerateStream(ctx, model, currentPrompt, format, maxTokens, func(content string, done bool) {
				currentIterationResponse.WriteString(content)
				if req.Schema == "" {
					// Normal mode: stream directly to client
					utils.WriteSSEEvent(w, flusher, ChatEvent{Content: content, Done: done, ChatID: chat.ID})
				}
			})

			if err != nil {
				slog.Error("ollama generation failed", "request_id", requestID, "error", err)
				utils.WriteSSEEvent(w, flusher, ChatEvent{Content: fmt.Sprintf("[error] %v", err), Done: true, ChatID: chat.ID})
				return
			}

			// If no schema, we're done
			if req.Schema == "" {
				fullAssistantResponse.WriteString(currentIterationResponse.String())
				break
			}

			// ── Validate JSON Response ──────────────────────────────
			responseStr := currentIterationResponse.String()
			slog.Info("validating response", "request_id", requestID, "len", len(responseStr))

			validationErr := validator.Validate(responseStr, req.Schema)
			if validationErr == nil {
				// SUCCESS: Valid JSON
				slog.Info("schema validation passed", "request_id", requestID)
				utils.WriteSSEEvent(w, flusher, ChatEvent{
					Content: responseStr,
					Done:    true,
					Status:  "valid",
					ChatID:  chat.ID,
				})
				fullAssistantResponse.WriteString(responseStr)
				break
			}

			// FAILURE: Invalid JSON
			slog.Warn("schema validation failed", "request_id", requestID, "error", validationErr)
			retryCount++
			if retryCount >= maxRetries {
				slog.Error("max retries exhausted", "request_id", requestID)
				utils.WriteSSEEvent(w, flusher, ChatEvent{
					Content: responseStr, // Show last malformed output
					Done:    true,
					Status:  "invalid",
					ChatID:  chat.ID,
				})
				fullAssistantResponse.WriteString(responseStr)
				break
			}

			currentPrompt = fmt.Sprintf("%s\n\nIMPORTANT: Your last response was invalid JSON. Error: %v. Please try again and return ONLY valid JSON matching the requested schema.", req.Message, validationErr)
		}

		// ── Persist Chat History ───────────────────────────────────
		// Simplified: append message to history string (Better to use a JSON array of messages)
		newMessage := fmt.Sprintf("\nUser: %s\nAssistant: %s", req.Message, fullAssistantResponse.String())
		chat.History += newMessage

		if db != nil {
			_, err := db.NewInsert().
				Model(chat).
				On("CONFLICT (id) DO UPDATE").
				Set("history = EXCLUDED.history").
				Set("updated_at = current_timestamp").
				Exec(ctx)

			if err != nil {
				slog.Error("failed to save chat history", "error", err)
			}
		}
	}
}

// GetChatsHandler returns all chats for the authenticated user.
func GetChatsHandler(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

		if db == nil {
			json.NewEncoder(w).Encode([]models.Chat{})
			return
		}

		var chats []models.Chat
		err := db.NewSelect().
			Model(&chats).
			Where("user_id = ?", userID).
			Order("updated_at DESC").
			Scan(r.Context())

		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(chats)
	}
}

func GetChatHandler(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for single chat retrieval if needed
	}
}

