package batch_test

import (
	"testing"

	"github.com/replay/platform/apps/ingestion/internal/batch"
)

func TestUploadBatchWorkersDefault(t *testing.T) {
	u := &batch.Uploader{}
	if u.Workers != 0 {
		t.Fatalf("expected zero default")
	}
}
