package middleware

import (
	"context"
	"net/http"
	"strings"
)

type shareScopeCtxKey string

const (
	ctxShareScope  shareScopeCtxKey = "share_scope"
	ctxShareIncID  shareScopeCtxKey = "share_incident_id"
	ctxShareOrgID  shareScopeCtxKey = "share_org_id"
	ShareScopeRead = "read_only"
)

// ShareScopeFromContext returns share token scope.
func ShareScopeFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxShareScope).(string)
	return v, ok && v != ""
}

// ShareIncidentIDFromContext returns incident id authorized by share token.
func ShareIncidentIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxShareIncID).(string)
	return v, ok && v != ""
}

// ShareOrgIDFromContext returns org id for the share token.
func ShareOrgIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxShareOrgID).(string)
	return v, ok && v != ""
}

// ShareTokenAuth validates read-only share tokens from header or query.
func ShareTokenAuth(lookup func(token string) (incidentID, orgID, scope string, expired bool, ok bool)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Share-Token")
			if token == "" {
				token = r.URL.Query().Get("token")
			}
			if token == "" {
				http.Error(w, `{"error":"unauthorized","message":"share token required"}`, http.StatusUnauthorized)
				return
			}
			incidentID, orgID, scope, expired, ok := lookup(token)
			if !ok {
				http.Error(w, `{"error":"unauthorized","message":"invalid share token"}`, http.StatusUnauthorized)
				return
			}
			if expired {
				http.Error(w, `{"error":"gone","message":"share token expired"}`, http.StatusGone)
				return
			}
			if strings.TrimSpace(scope) != ShareScopeRead {
				http.Error(w, `{"error":"forbidden","message":"invalid scope"}`, http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), ctxShareScope, scope)
			ctx = context.WithValue(ctx, ctxShareIncID, incidentID)
			ctx = context.WithValue(ctx, ctxShareOrgID, orgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireReadOnly rejects mutating methods for share-scoped requests.
func RequireReadOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := ShareScopeFromContext(r.Context()); !ok {
			next.ServeHTTP(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, `{"error":"forbidden","message":"share token is read-only"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
