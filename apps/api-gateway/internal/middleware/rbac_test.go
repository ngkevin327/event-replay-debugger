package middleware_test

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func TestRBAC(t *testing.T) {
	cases := []struct {
		actual store.MembershipRole
		min    store.MembershipRole
		allow  bool
	}{
		{store.RoleAdmin, store.RoleViewer, true},
		{store.RoleViewer, store.RoleAdmin, false},
		{store.RoleMember, store.RoleMember, true},
	}
	for _, c := range cases {
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
