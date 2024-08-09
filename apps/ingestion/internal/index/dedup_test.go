package index_test

import (
	"testing"

	"github.com/replay/platform/apps/ingestion/internal/index"
)

func TestDedup(t *testing.T) {
	d := index.NewDedup()
	key := index.DedupKey{ProjectID: "p", Topic: "t", Partition: 0, Offset: 1, ConsumerGroup: "g"}
	if d.IsDuplicate(key) {
		t.Fatal("first should not duplicate")
	}
	if !d.IsDuplicate(key) {
		t.Fatal("second should duplicate")
	}
}
