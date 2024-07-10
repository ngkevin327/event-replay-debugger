package auth_test

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/auth"
)

func TestAPIKey(t *testing.T) {
	plain, prefix, hash, err := auth.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	if prefix == "" || hash == "" {
		t.Fatal("expected prefix and hash")
	}
	got, err := auth.HashKey(plain)
	if err != nil || got != hash {
		t.Fatalf("hash mismatch got=%s want=%s", got, hash)
	}
	if !auth.ValidateScope([]string{"ingest", "read"}, "ingest") {
		t.Fatal("expected ingest scope")
	}
	if auth.ValidateScope([]string{"read"}, "ingest") {
		t.Fatal("expected missing scope")
	}
}
