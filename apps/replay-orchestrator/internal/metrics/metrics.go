package metrics

import "sync/atomic"

var (
	replayDurationMs atomic.Uint64
	divergedTotal    atomic.Uint64
)

// ReplayDuration records the latest replay duration in ms.
func ReplayDuration(ms uint64) {
	replayDurationMs.Store(ms)
}

// DivergedTotal increments diverged replay counter.
func DivergedTotal() {
	divergedTotal.Add(1)
}

// Snapshot exports metric values.
func Snapshot() (durationMs, diverged uint64) {
	return replayDurationMs.Load(), divergedTotal.Load()
}
