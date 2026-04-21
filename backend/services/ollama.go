// Package services contains external service clients for the PromptOps Engine
// handlers. Each service encapsulates a single external dependency, keeping
// the handler layer thin and testable.
//
// Currently provided services:
//   - OllamaClient: streams LLM inference via the Ollama REST API
package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/promptops/backend/pkg/metrics"
)

// -----------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------

// OllamaClient communicates with a running Ollama instance to generate
// text completions. It wraps the Ollama /api/generate endpoint and
// supports streaming NDJSON responses.
type OllamaClient struct {
	// BaseURL is the root URL of the Ollama server, e.g. "http://localhost:11434".
	BaseURL string

	// httpClient is the underlying HTTP client used for requests.
	// A default client is used if nil.
	httpClient *http.Client
}

// GenerateRequest is the payload sent to Ollama's /api/generate endpoint.
// Only the fields used by PromptOps Engine are included; Ollama accepts many
// more options (temperature, top_k, etc.) that can be added later.
type GenerateRequest struct {
	Model   string          `json:"model"`           // Model name, e.g. "tinyllama"
	Prompt  string          `json:"prompt"`          // The user's chat message
	Stream  bool            `json:"stream"`          // Must be true for streaming responses
	Format  string          `json:"format,omitempty"` // Set to "json" for structured output
	Options *GenerateOptions `json:"options,omitempty"` // Optional generation parameters
}

// GenerateOptions holds optional Ollama generation parameters.
// See https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values
type GenerateOptions struct {
	NumPredict int `json:"num_predict,omitempty"` // Max tokens to generate (0 = unlimited)
}

// GenerateChunk represents a single line of Ollama's NDJSON streaming
// response. Each chunk contains a fragment of the generated text and a
// flag indicating whether generation is complete.
type GenerateChunk struct {
	Response      string `json:"response"` // Text fragment (may be a single token)
	Done          bool   `json:"done"`     // True on the final chunk
	PromptEvalCount int    `json:"prompt_eval_count"` // Number of tokens in prompt
	EvalCount       int    `json:"eval_count"`        // Number of tokens in response
}

// -----------------------------------------------------------------------
// Constructor
// -----------------------------------------------------------------------

// NewOllamaClient creates an OllamaClient configured to talk to the given
// baseURL (e.g. "http://localhost:11434"). The client uses Go's default
// HTTP transport, which supports keep-alive for connection reuse.
func NewOllamaClient(baseURL string) *OllamaClient {
	return &OllamaClient{
		BaseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// -----------------------------------------------------------------------
// Streaming Generation
// -----------------------------------------------------------------------

// GenerateStream sends a prompt to Ollama and invokes the callback for
// every chunk of generated text. The callback receives the text fragment
// and a boolean indicating whether generation is complete.
//
// Flow:
//  1. POST to /api/generate with stream:true
//  2. Read response body line-by-line (NDJSON)
//  3. Parse each line into a GenerateChunk
//  4. Call onChunk(chunk.Response, chunk.Done)
//  5. Stop when chunk.Done == true
//
// Errors from the HTTP request, non-200 status codes, or JSON parse
// failures are returned immediately.
// GenerateStream sends a prompt to Ollama and invokes the callback for
// every chunk of generated text.
func (c *OllamaClient) GenerateStream(ctx context.Context, model, prompt, format string, maxTokens int, onChunk func(content string, done bool)) error {
	start := time.Now()
	// Build request payload
	reqBody := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
		Format: format,
	}

	// Apply token limit if specified
	if maxTokens > 0 {
		reqBody.Options = &GenerateOptions{NumPredict: maxTokens}
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	// Make the HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	// Read streaming NDJSON response line by line
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk GenerateChunk
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			return fmt.Errorf("parse chunk: %w", err)
		}

		onChunk(chunk.Response, chunk.Done)

		if chunk.Done {
			// Record metrics on completion
			duration := time.Since(start).Seconds()
			metrics.OllamaRequestDuration.WithLabelValues(model).Observe(duration)

			if chunk.PromptEvalCount > 0 {
				metrics.OllamaTokenUsageTotal.WithLabelValues(model, "prompt").Add(float64(chunk.PromptEvalCount))
			}
			if chunk.EvalCount > 0 {
				metrics.OllamaTokenUsageTotal.WithLabelValues(model, "completion").Add(float64(chunk.EvalCount))
			}
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading stream: %w", err)
	}

	return nil
}
