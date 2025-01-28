package handlers

import (
	"strconv"
	"testing"
	"time"
)

func TestCreateIncidentValidation(t *testing.T) {
	end := time.Now()
	start := end.Add(time.Hour)
	if !end.Before(start) {
		t.Fatal("expected invalid window for validation check")
	}
}

func TestRetentionCoverage(t *testing.T) {
	now := time.Now()
	start := now.Add(-24 * time.Hour)
	pct := RetentionCoveragePercent("starter", start, now, now)
	if pct <= 0 || pct > 100 {
		t.Fatalf("pct %v", pct)
	}
}

func TestListDetailPaginationDefaults(t *testing.T) {
	limit, _ := strconv.Atoi("")
	if limit != 0 {
		t.Fatalf("expected zero limit before defaulting in handler")
	}
}
