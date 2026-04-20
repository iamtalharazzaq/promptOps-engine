package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteSSEEvent formats and writes an SSE event to the response writer.
func WriteSSEEvent(w http.ResponseWriter, flusher http.Flusher, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	// SSE format: "data: <json>\n\n"
	if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
		return fmt.Errorf("write sse: %w", err)
	}
	
	flusher.Flush()
	return nil
}

// ErrorJSON writes a JSON error response.
func ErrorJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
