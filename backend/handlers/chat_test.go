package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"
	"github.com/promptops/backend/handlers"
	"github.com/promptops/backend/middleware"
	"github.com/promptops/backend/services"
)

var _ = Describe("ChatHandler", func() {
	var (
		ollamaServer *httptest.Server
		ollamaClient *services.OllamaClient
		validator    *services.JSONValidator
		defaultModel string
		maxTokens    int
		mockUserID   uuid.UUID
	)

	BeforeEach(func() {
		defaultModel = "test-model"
		maxTokens = 100
		validator = services.NewJSONValidator()
		mockUserID = uuid.New()
	})

	// Helper to inject user into context
	withUser := func(req *http.Request) *http.Request {
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, mockUserID)
		return req.WithContext(ctx)
	}

	AfterEach(func() {
		if ollamaServer != nil {
			ollamaServer.Close()
		}
	})

	// Helper to parse SSE events from recorder
	parseEvents := func(body string) []handlers.ChatEvent {
		var events []handlers.ChatEvent
		lines := strings.Split(body, "\n\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "data: ") {
				var ev handlers.ChatEvent
				json.Unmarshal([]byte(line[6:]), &ev)
				events = append(events, ev)
			}
		}
		return events
	}

	Context("Schema Guard Mode", func() {
		It("should return a valid JSON response when schema is matched", func() {
			ollamaServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/x-ndjson")
				fmt.Fprintf(w, `{"response":"{\"name\":\"John\", \"age\":30}", "done":true}`+"\n")
			}))
			ollamaClient = services.NewOllamaClient(ollamaServer.URL)

			schema := `{
				"type": "object",
				"properties": {
					"name": { "type": "string" },
					"age": { "type": "number" }
				},
				"required": ["name", "age"]
			}`

			reqBody := handlers.ChatRequest{
				Message: "Generate a profile",
				Schema:  schema,
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/chat", bytes.NewReader(body))
			req = withUser(req)
			rr := httptest.NewRecorder()

			handler := handlers.ChatHandler(nil, ollamaClient, validator, defaultModel, maxTokens)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			
			events := parseEvents(rr.Body.String())
			
			var validEvent *handlers.ChatEvent
			for _, ev := range events {
				if ev.Status == "valid" {
					validEvent = &ev
					break
				}
			}
			
			Expect(validEvent).NotTo(BeNil())
			Expect(validEvent.Content).To(ContainSubstring(`"name":"John"`))
		})

		It("should retry when response is invalid JSON", func() {
			callCount := 0
			ollamaServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				w.Header().Set("Content-Type", "application/x-ndjson")
				if callCount == 1 {
					fmt.Fprintf(w, `{"response":"invalid-json", "done":true}`+"\n")
				} else {
					fmt.Fprintf(w, `{"response":"{\"name\":\"Jane\", \"age\":25}", "done":true}`+"\n")
				}
			}))
			ollamaClient = services.NewOllamaClient(ollamaServer.URL)

			schema := `{"type":"object","properties":{"name":{"type":"string"},"age":{"type":"number"}},"required":["name","age"]}`

			reqBody := handlers.ChatRequest{
				Message: "Generate a profile",
				Schema:  schema,
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/chat", bytes.NewReader(body))
			req = withUser(req)
			rr := httptest.NewRecorder()

			handler := handlers.ChatHandler(nil, ollamaClient, validator, defaultModel, maxTokens)
			handler.ServeHTTP(rr, req)

			Expect(callCount).To(Equal(2))
			
			events := parseEvents(rr.Body.String())
			
			retryFound := false
			validFound := false
			for _, ev := range events {
				if ev.Status == "retrying" {
					retryFound = true
				}
				if ev.Status == "valid" {
					validFound = true
					Expect(ev.Content).To(ContainSubstring(`"name":"Jane"`))
				}
			}
			Expect(retryFound).To(BeTrue())
			Expect(validFound).To(BeTrue())
		})
	})

	Context("Normal Mode (No Schema)", func() {
		It("should stream tokens directly to client", func() {
			ollamaServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/x-ndjson")
				fmt.Fprintf(w, `{"response":"Hello", "done":false}`+"\n")
				fmt.Fprintf(w, `{"response":" world", "done":true}`+"\n")
			}))
			ollamaClient = services.NewOllamaClient(ollamaServer.URL)

			reqBody := handlers.ChatRequest{
				Message: "Hi",
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/chat", bytes.NewReader(body))
			req = withUser(req)
			rr := httptest.NewRecorder()

			handler := handlers.ChatHandler(nil, ollamaClient, validator, defaultModel, maxTokens)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			
			events := parseEvents(rr.Body.String())
			Expect(len(events)).To(BeNumerically(">=", 2))
			Expect(events[0].Content).To(Equal("Hello"))
			Expect(events[len(events)-1].Done).To(BeTrue())
		})
	})
})

