package storage_test

import (
	"testing"

	"github.com/replay/platform/apps/ingestion/internal/storage"
)

func TestPayloadKeyLayout(t *testing.T) {
	key := storage.PayloadKey("org-1", "proj-1", "abc123")
	if key != "payloads/org-1/proj-1/abc123" {
		t.Fatalf("key %s", key)
	}
}
