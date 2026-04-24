package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/promptops/backend/handlers"
)

var _ = Describe("AuthHandler", func() {
	var jwtSecret = "test-secret"

	Describe("Register", func() {
		It("should return a token upon successful registration", func() {
			reqBody := handlers.AuthRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			handler := handlers.RegisterHandler(nil, jwtSecret)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var resp handlers.AuthResponse
			err := json.Unmarshal(rr.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Token).NotTo(BeEmpty())
			Expect(resp.User.Email).To(Equal("test@example.com"))
		})
	})

	Describe("Login", func() {
		It("should return a token upon successful login", func() {
			reqBody := handlers.AuthRequest{
				Email:    "test@example.com",
				Password: "password123",
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			handler := handlers.LoginHandler(nil, jwtSecret)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var resp handlers.AuthResponse
			err := json.Unmarshal(rr.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Token).NotTo(BeEmpty())
		})
	})
})
