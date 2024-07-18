package audit

import (
	"context"
	"log/slog"

	"github.com/replay/platform/apps/api-gateway/internal/db"
)

// Logger records append-only audit events.
type Logger struct {
	pool *db.Pool
	ch   chan entry
}

type entry struct {
	actorID      string
	action       string
	resourceType string
	resourceID   string
	ip           string
}

// New creates an async audit logger.
func New(pool *db.Pool) *Logger {
	l := &Logger{pool: pool, ch: make(chan entry, 128)}
	go l.worker()
	return l
}

// LogAction enqueues an audit row insert.
func (l *Logger) LogAction(actorID, action, resourceType, resourceID, ip string) {
	l.ch <- entry{actorID: actorID, action: action, resourceType: resourceType, resourceID: resourceID, ip: ip}
}

func (l *Logger) worker() {
	for e := range l.ch {
		ctx := context.Background()
		_, err := l.pool.Exec(ctx,
			`INSERT INTO audit_logs (actor_id, action, resource_type, resource_id, ip)
			 VALUES ($1, $2, $3, $4, NULLIF($5, '')::inet)`,
			e.actorID, e.action, e.resourceType, e.resourceID, e.ip,
		)
		if err != nil {
			slog.Error("audit insert failed", "err", err)
		}
	}
}
