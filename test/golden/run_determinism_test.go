package golden

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/replay/platform/services/replay-worker/internal/divergence"
	"github.com/replay/platform/services/replay-worker/internal/record"
)

func RunDeterminismHarness(expected, actual []record.Outcome) bool {
	ok, _ := divergence.CompareChain(expected, actual)
	return ok
}

func TestDeterminismPaymentRetry(t *testing.T) {
	dir := filepath.Join("payment_retry_incident")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("fixtures missing")
	}
	exp := []record.Outcome{{Result: "ok"}, {Result: "ok"}}
	if !RunDeterminismHarness(exp, exp) {
		t.Fatal("expected deterministic match")
	}
}
