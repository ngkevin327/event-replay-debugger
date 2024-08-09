package index

import "sync"

// DedupKey identifies a Kafka record for idempotent ingest.
type DedupKey struct {
	ProjectID     string
	Topic         string
	Partition     int32
	Offset        int64
	ConsumerGroup string
}

// Dedup tracks seen keys in-process for MVP idempotency.
type Dedup struct {
	mu   sync.Mutex
	seen map[DedupKey]struct{}
}

// NewDedup creates an empty deduplicator.
func NewDedup() *Dedup {
	return &Dedup{seen: make(map[DedupKey]struct{})}
}

// IsDuplicate reports whether the key was already accepted.
func (d *Dedup) IsDuplicate(key DedupKey) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.seen[key]; ok {
		return true
	}
	d.seen[key] = struct{}{}
	return false
}
