package audit_test

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/audit"
	"github.com/replay/platform/apps/api-gateway/internal/db"
)

func TestAuditLoggerConstruct(t *testing.T) {
	l := audit.New(&db.Pool{})
	if l == nil {
		t.Fatal("expected logger")
	}
}
