package handlers

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileArtifactLoader(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "..", "fixtures", "local")
	if _, err := os.Stat(dir); err != nil {
		t.Skip("fixtures/local not found from module path")
	}
	loader := FileArtifactLoader{BaseDir: dir}
	body, ver, err := loader.LoadTimeline(context.Background(), "inc-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 || ver != 1 {
		t.Fatalf("unexpected timeline body/version")
	}
	g, err := loader.LoadGraph(context.Background(), "inc-1")
	if err != nil || len(g) == 0 {
		t.Fatalf("graph: %v", err)
	}
}
