//go:build integration

package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/auth"
	"github.com/replay/platform/apps/api-gateway/internal/config"
	"github.com/replay/platform/apps/api-gateway/internal/db"
	"github.com/replay/platform/apps/api-gateway/internal/server"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func TestCrossProjectDenied(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := context.Background()
	pool, err := db.OpenPostgres(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	st := store.NewStore(pool)

	secret := "integration-test-secret-32chars"
	orgA, _ := st.CreateOrganization(ctx, "org-a", "starter")
	orgB, _ := st.CreateOrganization(ctx, "org-b", "starter")
	userA, _ := st.CreateUser(ctx, "a@example.com", mustHash(t, "password-12chars"))
	_ = st.CreateMembership(ctx, orgA.ID, userA.ID, store.RoleAdmin)
	projectB, _ := st.CreateProject(ctx, orgB.ID, "project-b")

	token, err := auth.IssueTokens(secret, userA.ID, orgA.ID)
	if err != nil {
		t.Fatal(err)
	}

	srv := server.New(config.Config{HTTPAddr: ":0"}, server.RouteDeps{
		Store:     st,
		JWTSecret: secret,
	})
	req := httptest.NewRequest(http.MethodGet, "/v1/projects/"+projectB.ID, nil)
	req.Header.Set("Authorization", "Bearer "+token.Access)
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func mustHash(t *testing.T, pw string) string {
	t.Helper()
	h, err := auth.HashPassword(pw)
	if err != nil {
		t.Fatal(err)
	}
	return h
}
