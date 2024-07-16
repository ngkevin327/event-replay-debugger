package middleware

import (
	"net/http"

	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// RequireRole ensures the session user has at least the given org role.
func RequireRole(st *store.Store, min store.MembershipRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			orgID, ok2 := OrgIDFromContext(r.Context())
			if !ok || !ok2 {
				http.Error(w, `{"error":"unauthorized","message":"session required"}`, http.StatusUnauthorized)
				return
			}
			role, err := st.GetMembershipRole(r.Context(), orgID, userID)
			if err != nil {
				http.Error(w, `{"error":"forbidden","message":"membership required"}`, http.StatusForbidden)
				return
			}
			if !roleAtLeast(role, min) {
				http.Error(w, `{"error":"forbidden","message":"insufficient role"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func roleAtLeast(actual, min store.MembershipRole) bool {
	order := map[store.MembershipRole]int{
		store.RoleViewer: 1,
		store.RoleMember: 2,
		store.RoleAdmin:  3,
	}
	return order[actual] >= order[min]
}
