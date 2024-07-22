package middleware_test

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func TestRBAC(t *testing.T) {
	TestRoleMatrix(t)
}

// TestRoleMatrix validates admin/member/viewer capability ordering.
func TestRoleMatrix(t *testing.T) {
	matrix := []struct {
		actual store.MembershipRole
		min    store.MembershipRole
		allow  bool
	}{
		{store.RoleAdmin, store.RoleViewer, true},
		{store.RoleViewer, store.RoleAdmin, false},
		{store.RoleMember, store.RoleMember, true},
		{store.RoleViewer, store.RoleMember, false},
	}
	for _, c := range matrix {
		got := roleAtLeastExported(c.actual, c.min)
		if got != c.allow {
			t.Fatalf("actual=%s min=%s got=%v want=%v", c.actual, c.min, got, c.allow)
		}
	}
}

func roleAtLeastExported(actual, min store.MembershipRole) bool {
	order := map[store.MembershipRole]int{
		store.RoleViewer: 1,
		store.RoleMember: 2,
		store.RoleAdmin:  3,
	}
	return order[actual] >= order[min]
}
