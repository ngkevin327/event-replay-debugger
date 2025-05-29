package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/replay/platform/apps/api-gateway/internal/auth"
)

type sessionCtxKey string

const (
	ctxUserID sessionCtxKey = "user_id"
	ctxOrgID  sessionCtxKey = "org_id"
)

// UserIDFromContext returns the authenticated user id.
func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxUserID).(string)
	return v, ok && v != ""
}

// TestContextWithOrg attaches user/org for handler unit tests.
func TestContextWithOrg(ctx context.Context, orgID string) context.Context {
	ctx = context.WithValue(ctx, ctxUserID, "test-user")
	return context.WithValue(ctx, ctxOrgID, orgID)
}

// OrgIDFromContext returns the active organization id.
func OrgIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxOrgID).(string)
	return v, ok && v != ""
}

// RequireSession validates Bearer access tokens and attaches user/org to context.
func RequireSession(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"unauthorized","message":"missing bearer token"}`, http.StatusUnauthorized)
				return
			}
			token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
			claims, err := auth.ParseAccessToken(jwtSecret, token)
			if err != nil {
				http.Error(w, `{"error":"unauthorized","message":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
			ctx = context.WithValue(ctx, ctxOrgID, claims.OrgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
