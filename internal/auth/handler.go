package auth

import (
	"encoding/json"
	"errors"
	"go-flix-api/config"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// --- Struct untuk Data & Dependensi ---

// Service adalah tempat semua logika bisnis inti.
type Service struct {
	cfg      *config.Config
	denylist map[string]time.Time // Database token yang sudah di-logout
	mu       sync.RWMutex
}

// Handler adalah lapisan HTTP yang menerima request.
type Handler struct {
	service *Service
}

// JWTClaims adalah data yang kita simpan di dalam token.
type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// --- Konstruktor (Fungsi "Pabrik") ---

// NewService membuat instance baru dari Service.
func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg:      cfg,
		denylist: make(map[string]time.Time),
	}
}

// NewHandler membuat instance baru dari Handler.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// --- Method-Method Service (Logika Inti) ---

// ValidateUser memeriksa username & password.
func (s *Service) ValidateUser(username, password string) bool {
	for _, u := range s.cfg.Users {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

// GenerateJWT membuat token JWT baru.
func (s *Service) GenerateJWT(username string) (string, error) {
	claims := JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(), // ID unik untuk setiap token
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

// RevokeToken menambahkan token ke denylist (untuk logout).
func (s *Service) RevokeToken(tokenStr string) error {
	claims := &JWTClaims{}
	// PERBAIKAN PENTING: Kita tetap validasi token sebelum di-logout untuk keamanan.
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	// Simpan ID token dan waktu kedaluwarsanya.
	s.denylist[claims.ID] = claims.ExpiresAt.Time
	return nil
}

// IsTokenRevoked memeriksa apakah token ada di denylist.
// Ini akan dipanggil oleh middleware.
func (s *Service) IsTokenRevoked(jti string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, revoked := s.denylist[jti]
	return revoked
}

// --- Method-Method Handler (Lapisan HTTP) ---

// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "Login credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "Invalid JSON"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Failed to generate token"
// @Router /api/login [post]
// Login menangani POST /api/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Memanggil service untuk validasi.
	if !h.service.ValidateUser(req.Username, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Memanggil service untuk membuat token.
	tokenStr, err := h.service.GenerateJWT(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

// @Summary User logout
// @Description Revoke JWT token (logout)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "message"
// @Failure 401 {object} map[string]string "Missing or invalid Authorization header"
// @Router /api/logout [post]
// Logout menangani POST /api/logout.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Memanggil service untuk me-revoke token.
	if err := h.service.RevokeToken(tokenStr); err != nil {
		http.Error(w, "Invalid token for logout", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logout successful"})
}
