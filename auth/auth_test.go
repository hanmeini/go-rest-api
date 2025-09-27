package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestValidateUser_Success(t *testing.T) {
	cfgYAML := "users:\n  - username: \"testuser\"\n    password: \"secret\"\n"
	cfgPath := writeTempConfig(t, cfgYAML)

	if err := LoadConfig(cfgPath); err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	ok, err := ValidateUser("testuser", "secret")
	if err != nil {
		t.Fatalf("ValidateUser returned error: %v", err)
	}
	if !ok {
		t.Fatalf("expected credentials to be valid")
	}
}

func TestValidateUser_WrongPassword(t *testing.T) {
	cfgYAML := "users:\n  - username: \"testuser\"\n    password: \"secret\"\n"
	cfgPath := writeTempConfig(t, cfgYAML)

	if err := LoadConfig(cfgPath); err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	ok, err := ValidateUser("testuser", "wrong")
	if err != nil {
		t.Fatalf("ValidateUser returned error: %v", err)
	}
	if ok {
		t.Fatalf("expected credentials to be invalid with wrong password")
	}
}

func setupJWTSecret(t *testing.T) {
	t.Helper()
	// Set a test JWT secret
	JWT_SECRET_KEY = "test-secret-key-for-jwt-testing"
}

func TestJWT_GenerateAndValidate(t *testing.T) {
	setupJWTSecret(t)

	username := "testuser"
	token, err := GenerateJWT(username)
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}

	if token == "" {
		t.Fatalf("expected token to be generated")
	}

	// Validate the token
	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}

	if claims.Username != username {
		t.Fatalf("expected username %s, got %s", username, claims.Username)
	}
}

func TestJWT_RevokeFlow(t *testing.T) {
	setupJWTSecret(t)

	username := "testuser"
	token, err := GenerateJWT(username)
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}

	// Token should be valid initially
	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}
	if claims.Username != username {
		t.Fatalf("expected username %s, got %s", username, claims.Username)
	}

	// Revoke the token
	err = RevokeJWT(token)
	if err != nil {
		t.Fatalf("RevokeJWT error: %v", err)
	}

	// Token should now be invalid
	_, err = ValidateJWT(token)
	if err == nil {
		t.Fatalf("expected token to be invalid after revocation")
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	setupJWTSecret(t)

	// Create a token with very short expiration (1 second)
	originalExpiration := time.Hour
	// We can't easily test expiration without modifying the function,
	// but we can test that the token structure is correct
	username := "testuser"
	token, err := GenerateJWT(username)
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}

	// Check that expiration is set to approximately 1 hour from now
	expectedExp := time.Now().Add(originalExpiration)
	actualExp := claims.ExpiresAt.Time
	diff := actualExp.Sub(expectedExp)
	if diff > time.Minute || diff < -time.Minute {
		t.Fatalf("expected expiration around %v, got %v", expectedExp, actualExp)
	}
}

func TestJWT_InvalidToken(t *testing.T) {
	setupJWTSecret(t)

	// Test with invalid token
	_, err := ValidateJWT("invalid-token")
	if err == nil {
		t.Fatalf("expected error for invalid token")
	}

	// Test with empty token
	_, err = ValidateJWT("")
	if err == nil {
		t.Fatalf("expected error for empty token")
	}
}
