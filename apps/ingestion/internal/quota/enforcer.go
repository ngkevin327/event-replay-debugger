package quota

import (
	"context"
	"errors"
	"sync/atomic"
)

const (
	maxPayloadBytes = 256 * 1024
	defaultDailyCap = 1_000_000
)

var ErrQuotaExceeded = errors.New("daily event quota exceeded")

// Enforcer applies per-tenant ingest limits.
type Enforcer struct {
	dailyCap int64
	used     atomic.Int64
}

// NewEnforcer creates a quota enforcer with starter-plan defaults.
func NewEnforcer(dailyCap int64) *Enforcer {
	if dailyCap <= 0 {
		dailyCap = defaultDailyCap
	}
	return &Enforcer{dailyCap: dailyCap}
}

// EnforceQuota rejects batches that exceed remaining daily allowance.
func (e *Enforcer) EnforceQuota(_ context.Context, _ string, batchSize int) error {
	if int64(batchSize)+e.used.Load() > e.dailyCap {
		return ErrQuotaExceeded
	}
	e.used.Add(int64(batchSize))
	return nil
}

// PayloadSizeLimit returns max payload bytes per event.
func PayloadSizeLimit() int {
	return maxPayloadBytes
}
