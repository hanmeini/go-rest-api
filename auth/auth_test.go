package auth

import (
	"os"
	"path/filepath"
	"testing"
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

func TestToken_RevokeFlow(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	if IsTokenValid(token) {
		t.Fatalf("token should not be valid before storing")
	}

	StoreToken(token)
	if !IsTokenValid(token) {
		t.Fatalf("expected token to be valid after storing")
	}

	RevokeToken(token)
	if IsTokenValid(token) {
		t.Fatalf("expected token to be invalid after revocation")
	}
}
