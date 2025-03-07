package cache

import (
	"sync"
	"time"
)

type entry struct {
	body    []byte
	version int
	expires time.Time
}

// TimelineCache stores hot timeline responses in memory (Redis in production).
type TimelineCache struct {
	mu    sync.RWMutex
	items map[string]entry
	ttl   time.Duration
}

// NewTimelineCache creates a cache with TTL.
func NewTimelineCache(ttl time.Duration) *TimelineCache {
	return &TimelineCache{items: make(map[string]entry), ttl: ttl}
}

// Get returns cached timeline bytes and version.
func (c *TimelineCache) Get(incidentID string) ([]byte, int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[incidentID]
	if !ok || time.Now().After(e.expires) {
		return nil, 0, false
	}
	return e.body, e.version, true
}

// Set stores timeline bytes.
func (c *TimelineCache) Set(incidentID string, body []byte, version int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[incidentID] = entry{
		body: body, version: version, expires: time.Now().Add(c.ttl),
	}
}
