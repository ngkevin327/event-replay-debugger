package feeder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTimeline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "timeline.json")
	_ = os.WriteFile(path, []byte(`{"events":[{"topic":"payments","partition":0,"offset":1}]}`), 0o644)
	art, err := LoadTimeline(context.Background(), "file://"+path)
	if err != nil || len(art.Events) != 1 {
		t.Fatal(err)
	}
}
