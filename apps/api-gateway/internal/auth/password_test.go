package auth_test

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/auth"
)

func TestPassword(t *testing.T) {
	hash, err := auth.HashPassword("valid-password-12")
	if err != nil {
		t.Fatal(err)
	}
	if err := auth.VerifyPassword(hash, "valid-password-12"); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := auth.VerifyPassword(hash, "wrong-password-12"); err == nil {
		t.Fatal("expected mismatch")
	}
}

func TestPasswordRejectsShort(t *testing.T) {
	if _, err := auth.HashPassword("short"); err != auth.ErrPasswordTooShort {
		t.Fatalf("got %v", err)
	}
}
