package payload_test

import (
	"strings"
	"testing"

	"github.com/replay/platform/apps/ingestion/internal/payload"
)

func TestProcessPayloadTruncates(t *testing.T) {
	raw := []byte(strings.Repeat("x", payload.TruncateAt256KB+1))
	data, hash, truncated := payload.ProcessPayload(raw)
	if !truncated || len(data) != payload.TruncateAt256KB {
		t.Fatalf("len=%d truncated=%v", len(data), truncated)
	}
	if hash == "" {
		t.Fatal("expected hash")
	}
}
