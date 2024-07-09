package auth_test

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/auth"
)

func TestJWT(t *testing.T) {
	secret := "test-secret-at-least-32-bytes-long"
	pair, err := auth.IssueTokens(secret, "user-1", "org-1")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := auth.ParseAccessToken(secret, pair.Access)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "user-1" || claims.OrgID != "org-1" {
		t.Fatalf("claims %+v", claims)
	}
	rotated, err := auth.Refresh(secret, pair.Refresh, "org-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := auth.ParseAccessToken(secret, rotated.Access); err != nil {
		t.Fatal(err)
	}
}
