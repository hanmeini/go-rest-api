package middleware

import (
	"net/http"
	"strings"

	"go-flix-api/auth"
)

// AuthMiddleware validates Bearer token and passes request to next handler if valid
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		if !auth.IsTokenValid(token) {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
