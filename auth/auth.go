package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"sync"

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

var (
	// cfg holds the loaded configuration in memory
	cfg Config

	// tokenStore is an in-memory set of valid tokens
	tokenStore = make(map[string]bool)

	// mutex protects concurrent access to cfg and tokenStore
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

// GenerateToken creates a cryptographically secure random token string
func GenerateToken() (string, error) {
	// 32 bytes -> 64 hex chars
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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

// StoreToken adds the token to the in-memory token database
func StoreToken(token string) {
	if token == "" {
		return
	}
	mutex.Lock()
	tokenStore[token] = true
	mutex.Unlock()
}

// IsTokenValid returns true if the token exists in the in-memory token database
func IsTokenValid(token string) bool {
	if token == "" {
		return false
	}
	mutex.RLock()
	valid := tokenStore[token]
	mutex.RUnlock()
	return valid
}

// RevokeToken removes the token from the in-memory token database
func RevokeToken(token string) {
	if token == "" {
		return
	}
	mutex.Lock()
	delete(tokenStore, token)
	mutex.Unlock()
}
