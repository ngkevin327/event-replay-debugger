package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/replay/platform/apps/api-gateway/internal/auth"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

const headerReplayKey = "X-Replay-Key"

type apiKeyCtx string

const (
	ctxProjectID apiKeyCtx = "project_id"
	ctxKeyScopes apiKeyCtx = "api_key_scopes"
)

// ProjectIDFromContext returns project id set by API key middleware.
func ProjectIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxProjectID).(string)
	return v, ok && v != ""
}

// RequireAPIKey validates X-Replay-Key and required scope.
func RequireAPIKey(st *store.Store, requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			plain := strings.TrimSpace(r.Header.Get(headerReplayKey))
			if plain == "" {
				http.Error(w, `{"error":"unauthorized","message":"missing api key"}`, http.StatusUnauthorized)
				return
			}
			if len(plain) < 16 {
				http.Error(w, `{"error":"unauthorized","message":"invalid api key"}`, http.StatusUnauthorized)
				return
			}
			prefix := plain[:16]
			row, err := st.GetAPIKeyByPrefix(r.Context(), prefix)
			if err != nil {
				http.Error(w, `{"error":"unauthorized","message":"invalid api key"}`, http.StatusUnauthorized)
				return
			}
			hash, err := auth.HashKey(plain)
			if err != nil || hash != row.Hash {
				http.Error(w, `{"error":"unauthorized","message":"invalid api key"}`, http.StatusUnauthorized)
				return
			}
			if !auth.ValidateScope(row.Scopes, requiredScope) {
				http.Error(w, `{"error":"forbidden","message":"insufficient scope"}`, http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), ctxProjectID, row.ProjectID)
			ctx = context.WithValue(ctx, ctxKeyScopes, row.Scopes)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
