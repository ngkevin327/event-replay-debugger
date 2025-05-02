package divergence

import (
	"testing"

	"github.com/replay/platform/services/replay-worker/internal/record"
)

func TestCompareChain(t *testing.T) {
	exp := []record.Outcome{{Result: "ok"}, {Result: "ok"}}
	act := []record.Outcome{{Result: "ok"}, {Result: "fail"}}
	ok, mm := CompareChain(exp, act)
	if ok || mm.Index != 1 {
		t.Fatalf("got ok=%v mm=%+v", ok, mm)
	}
}
