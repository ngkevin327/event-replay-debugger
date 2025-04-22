package feeder

import "time"

// StrictTiming sleeps until the next event timestamp relative to replay start.
func StrictTiming(start, eventTime time.Time, replayStart time.Time) {
	delay := eventTime.Sub(start)
	if delay < 0 {
		delay = 0
	}
	time.Sleep(time.Until(replayStart.Add(delay)))
}

// CompressedTiming scales inter-event delay by factor (>1 speeds up).
func CompressedTiming(prev, cur time.Time, factor float64) {
	if factor <= 0 {
		factor = 1
	}
	delay := time.Duration(float64(cur.Sub(prev)) / factor)
	if delay > 0 {
		time.Sleep(delay)
	}
}
