package auth

import (
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	yaml "gopkg.in/yaml.v2"
)

// User represents a single credential entry loaded from config.yaml
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Config holds application configuration loaded from config.yaml
type Config struct {
	Users []User `yaml:"users"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	// cfg holds the loaded configuration in memory
	cfg Config

	// JWT_SECRET_KEY holds the secret key for JWT signing
	JWT_SECRET_KEY string

	// denylist stores revoked JWT token IDs
	denylist = make(map[string]bool)

	// mutex protects concurrent access to cfg and denylist
	mutex sync.RWMutex
)

// LoadConfig reads a YAML file from the provided path and loads it into memory
func LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var newCfg Config
	if err := yaml.Unmarshal(data, &newCfg); err != nil {
		return err
	}

	mutex.Lock()
	cfg = newCfg
	mutex.Unlock()
	return nil
}

// LoadJWTSecret loads the JWT secret key from .env file
func LoadJWTSecret() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return errors.New("failed to load .env file: " + err.Error())
	}

	// Get JWT_SECRET_KEY from environment
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return errors.New("JWT_SECRET_KEY not found in environment variables")
	}

	JWT_SECRET_KEY = secret
	return nil
}

// GenerateJWT creates a new JWT token for the given username with 1-hour expiration
func GenerateJWT(username string) (string, error) {
	if JWT_SECRET_KEY == "" {
		return "", errors.New("JWT secret key not loaded")
	}

	// Create claims with 1-hour expiration
	claims := JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-flix-api",
			Subject:   username,
			ID:        generateTokenID(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT verifies the JWT token signature, expiration, and checks if it's not in denylist
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	if JWT_SECRET_KEY == "" {
		return nil, errors.New("JWT secret key not loaded")
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWT_SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check if token is in denylist
	mutex.RLock()
	if denylist[claims.ID] {
		mutex.RUnlock()
		return nil, errors.New("token has been revoked")
	}
	mutex.RUnlock()

	return claims, nil
}

// RevokeJWT adds the token ID to the denylist for logout functionality
func RevokeJWT(tokenString string) error {
	// Parse token to get claims without validation (since it might be expired)
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET_KEY), nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Add token ID to denylist
	mutex.Lock()
	denylist[claims.ID] = true
	mutex.Unlock()

	return nil
}

// ValidateUser checks if the provided username and password match the config
func ValidateUser(username, password string) (bool, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	if len(cfg.Users) == 0 {
		return false, errors.New("no users loaded in config")
	}

	for _, u := range cfg.Users {
		if u.Username == username && u.Password == password {
			return true, nil
		}
	}
	return false, nil
}

// generateTokenID creates a unique token ID for JWT claims
func generateTokenID() string {
	// Simple implementation using timestamp and random number
	// In production, you might want to use a more robust UUID generator
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
