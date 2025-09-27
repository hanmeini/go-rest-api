package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-flix-api/auth"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	valid, err := auth.ValidateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !valid {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(req.Username)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}

// LogoutHandler handles POST /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing Authorization header", http.StatusUnauthorized)
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		http.Error(w, "empty token", http.StatusUnauthorized)
		return
	}

	if err := auth.RevokeJWT(token); err != nil {
		http.Error(w, "failed to revoke token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logout successful"})
}
