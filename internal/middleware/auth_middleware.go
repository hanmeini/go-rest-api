package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type DenylistChecker func(jti string) bool

// AuthMiddleware memproteksi endpoint hanya untuk user login
// Param: secret JWT, fungsi cek denylist
func AuthMiddleware(secret string, isTokenRevoked DenylistChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(*JWTClaims)
			if !ok || (isTokenRevoked != nil && isTokenRevoked(claims.ID)) {
				http.Error(w, "Token revoked", http.StatusUnauthorized)
				return
			}
			// Inject username ke context/header
			r = r.WithContext(context.WithValue(r.Context(), "username", claims.Username))
			r.Header.Set("X-Username", claims.Username)
			next.ServeHTTP(w, r)
		})
	}
}
