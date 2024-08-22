package index_test

import (
	"testing"

	"github.com/replay/platform/apps/ingestion/internal/index"
)

func TestBackpressure(t *testing.T) {
	w := index.NewWriter(index.NewCHClient("http://localhost:8123"))
	if w.Backpressure() {
		t.Fatal("expected no backpressure initially")
	}
}
