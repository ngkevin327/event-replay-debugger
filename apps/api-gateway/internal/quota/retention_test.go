package quota

import (
	"testing"
	"time"
)

func TestRetentionClampStarter(t *testing.T) {
	now := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)
	start := now.Add(-30 * 24 * time.Hour)
	end := now
	clampedStart, clampedEnd := ClampQueryWindow("starter", start, end, now)
	if clampedStart.Before(now.Add(-8 * 24 * time.Hour)) {
		t.Fatalf("starter start too old: %v", clampedStart)
	}
	if clampedEnd.After(now) {
		t.Fatalf("end after now")
	}
}

func TestRetentionClampPro(t *testing.T) {
	now := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)
	start := now.Add(-60 * 24 * time.Hour)
	end := now
	clampedStart, _ := ClampQueryWindow("pro", start, end, now)
	earliest := now.Add(-31 * 24 * time.Hour)
	if clampedStart.Before(earliest) {
		t.Fatalf("pro start too old: %v", clampedStart)
	}
}

func TestRetentionDaysTiers(t *testing.T) {
	if RetentionDays("starter") != 7 {
		t.Fatal("starter")
	}
	if RetentionDays("pro") != 30 {
		t.Fatal("pro")
	}
}
