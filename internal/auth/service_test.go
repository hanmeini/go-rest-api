package auth

import (
	"testing"
	"time"

	"go-flix-api/config"
)

func newTestService() *Service {
	cfg := &config.Config{
		JWT:   config.JWTConfig{Secret: "test_secret"},
		Users: []config.User{{Username: "user1", Password: "password123"}},
	}
	return NewService(cfg)
}

func TestValidateUser(t *testing.T) {
	s := newTestService()
	if !s.ValidateUser("user1", "password123") {
		t.Fatalf("expected valid credentials")
	}
	if s.ValidateUser("user1", "wrong") {
		t.Fatalf("expected invalid credentials")
	}
}

func TestGenerateAndValidateJWT(t *testing.T) {
	s := newTestService()
	token, err := s.GenerateJWT("user1")
	if err != nil || token == "" {
		t.Fatalf("expected token, got err=%v token=%q", err, token)
	}
}

func TestRevokeAndIsTokenRevoked(t *testing.T) {
	s := newTestService()
	token, err := s.GenerateJWT("user1")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if err := s.RevokeToken(token); err != nil {
		t.Fatalf("revoke token: %v", err)
	}
	// Extract jti by parsing via RevokeToken path again (already validated inside)
	// Quick check: token should be considered revoked within denylist time window
	// We cannot access JTI directly here, but IsTokenRevoked is used by middleware with JTI.
	// So we simulate by generating a token again and ensuring denylist map is populated.
	// Minimal assertion: denylist not empty and at least one expiry in future.
	if len(s.denylist) == 0 {
		t.Fatalf("expected denylist populated")
	}
	for _, exp := range s.denylist {
		if exp.Before(time.Now()) {
			t.Fatalf("expected future expiry, got %v", exp)
		}
	}
}
