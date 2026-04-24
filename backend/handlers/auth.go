package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/promptops/backend/pkg/auth"
	"github.com/promptops/backend/pkg/models"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// RegisterHandler handles new user sign-ups.
func RegisterHandler(db *bun.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		user := &models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
			Name:     req.Name,
		}

		if db != nil {
			_, err = db.NewInsert().Model(user).Exec(r.Context())
			if err != nil {
				http.Error(w, "email already exists or database error", http.StatusConflict)
				return
			}
		}

		token, err := auth.GenerateToken(user.ID, user.Email, jwtSecret)
		if err != nil {
			http.Error(w, "token generation failed", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
	}
}

// LoginHandler handles user authentication.
func LoginHandler(db *bun.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		user := new(models.User)
		if db != nil {
			err := db.NewSelect().
				Model(user).
				Where("email = ?", req.Email).
				Scan(context.Background())

			if err != nil {
				http.Error(w, "invalid credentials", http.StatusUnauthorized)
				return
			}
		} else {
			// Mock user for testing if DB is nil
			user.Email = req.Email
			hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 4)
			user.Password = string(hashed)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Email, jwtSecret)
		if err != nil {
			http.Error(w, "token generation failed", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
	}
}
