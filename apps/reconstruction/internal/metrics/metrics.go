package metrics

import (
	"sync/atomic"
	"time"
)

var (
	jobDurationMs atomic.Uint64
	jobFailures   atomic.Uint64
)

// JobDuration records last job duration in milliseconds.
func JobDuration(d time.Duration) {
	jobDurationMs.Store(uint64(d.Milliseconds()))
}

// JobFailures increments failed reconstruction jobs.
func JobFailures() {
	jobFailures.Add(1)
}

// Snapshot returns metric values for scraping.
func Snapshot() (durationMs uint64, failures uint64) {
	return jobDurationMs.Load(), jobFailures.Load()
}
