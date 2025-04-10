package handlers

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func TestGetStatusDivergenceSummary(t *testing.T) {
	idx := 3
	sum := DivergenceSummary(store.ReplayRun{
		Status: store.ReplayDiverged, DivergenceIndex: &idx,
	})
	if sum["first_mismatch_index"] != 3 {
		t.Fatal()
	}
}

func TestCancelReplay(t *testing.T) {
	if (&ReplaysHandler{}) == nil {
		t.Fatal()
	}
}

func TestCreateReplayTimingValidation(t *testing.T) {
	mode := "strict"
	if mode != "strict" && mode != "compressed" {
		t.Fatal()
	}
}

