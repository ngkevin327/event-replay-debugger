package handlers

import (
	"encoding/json"
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/export"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func TestExportSummarySchema(t *testing.T) {
	coverage := 92.5
	inc := store.Incident{
		ID: "inc-1", ProjectID: "proj-1", Status: store.IncidentReady,
		EventCount: 100, CoveragePercent: &coverage,
	}
	summary := export.BuildSummary(inc)
	body, err := json.Marshal(summary)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) > export.MaxExportBytes {
		t.Fatal("export exceeds max size unexpectedly")
	}
	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["incident_id"] != "inc-1" {
		t.Fatalf("unexpected payload %v", decoded)
	}
}

func TestExportSizeLimitConstant(t *testing.T) {
	if export.MaxExportBytes < 512*1024 {
		t.Fatal("export limit too small for MVP")
	}
}
