package quota

import "sync"

// Limiter tracks concurrent replay runs for starter plans.
type Limiter struct {
	mu    sync.Mutex
	limit int
	active map[string]int
}

// NewLimiter creates a per-project concurrency limiter.
func NewLimiter(limit int) *Limiter {
	return &Limiter{limit: limit, active: make(map[string]int)}
}

// ConcurrentReplayLimit returns whether a new replay may start.
func (l *Limiter) ConcurrentReplayLimit(projectID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.active[projectID] < l.limit
}

// Acquire increments active runs.
func (l *Limiter) Acquire(projectID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.active[projectID]++
}

// Release decrements active runs.
func (l *Limiter) Release(projectID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.active[projectID] > 0 {
		l.active[projectID]--
	}
}
