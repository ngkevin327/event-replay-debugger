package filter_test

import (
	"testing"

	"github.com/replay/platform/packages/agent-go/filter"
)

func TestShouldCapture(t *testing.T) {
	a := filter.NewAllowlist([]string{"orders", "payments"})
	if !a.ShouldCapture("orders") {
		t.Fatal("expected orders allowed")
	}
	if a.ShouldCapture("other") {
		t.Fatal("expected other denied")
	}
}
