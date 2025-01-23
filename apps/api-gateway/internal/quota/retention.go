package quota

import "time"

const (
	starterRetentionDays = 7
	proRetentionDays     = 30
)

// ClampQueryWindow limits query bounds to the plan retention window.
func ClampQueryWindow(planTier string, windowStart, windowEnd time.Time, now time.Time) (time.Time, time.Time) {
	days := starterRetentionDays
	if planTier == "pro" || planTier == "enterprise" {
		days = proRetentionDays
	}
	earliest := now.Add(-time.Duration(days) * 24 * time.Hour)
	if windowStart.Before(earliest) {
		windowStart = earliest
	}
	if windowEnd.After(now) {
		windowEnd = now
	}
	if windowEnd.Before(windowStart) {
		windowEnd = windowStart
	}
	return windowStart, windowEnd
}

// RetentionDays returns configured retention for a plan tier.
func RetentionDays(planTier string) int {
	if planTier == "pro" || planTier == "enterprise" {
		return proRetentionDays
	}
	return starterRetentionDays
}
