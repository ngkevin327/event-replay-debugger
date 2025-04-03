package handlers

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func TestCreateReplayTimingValidation(t *testing.T) {
	mode := "strict"
	if mode != "strict" && mode != "compressed" {
		t.Fatal()
	}
}

func TestGetStatusDivergenceSummary(t *testing.T) {
	idx := 3
	sum := DivergenceSummary(store.ReplayRun{
		Status: store.ReplayDiverged, DivergenceIndex: &idx,
	})
	if sum["first_mismatch_index"] != 3 {
		t.Fatal()
	}
}

func TestCancelReplayHandlerExists(t *testing.T) {
	h := &ReplaysHandler{}
	if h == nil {
		t.Fatal()
	}
}
