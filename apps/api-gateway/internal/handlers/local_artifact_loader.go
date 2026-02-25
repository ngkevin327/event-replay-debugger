package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// FileArtifactLoader serves timeline/graph JSON from a fixtures directory for local MVP demos.
type FileArtifactLoader struct {
	BaseDir string
}

// LoadTimeline reads fixtures/local/timeline.json (or BaseDir/timeline.json).
func (f FileArtifactLoader) LoadTimeline(ctx context.Context, incidentID string) ([]byte, int, error) {
	_ = ctx
	_ = incidentID
	path := filepath.Join(f.BaseDir, "timeline.json")
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("timeline not ready: %w", err)
	}
	return body, 1, nil
}

// LoadGraph reads fixtures/local/graph.json.
func (f FileArtifactLoader) LoadGraph(ctx context.Context, incidentID string) ([]byte, error) {
	_ = ctx
	_ = incidentID
	path := filepath.Join(f.BaseDir, "graph.json")
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("graph not ready: %w", err)
	}
	return body, nil
}

// ResolveArtifactLoader picks timeline loader from environment.
func ResolveTimelineLoader() TimelineLoader {
	if os.Getenv("LOCAL_DEMO_MODE") != "1" {
		return StubTimelineLoader{}
	}
	dir := os.Getenv("LOCAL_DEMO_FIXTURES_DIR")
	if dir == "" {
		dir = "fixtures/local"
	}
	return FileArtifactLoader{BaseDir: dir}
}

// ResolveGraphLoader picks graph loader from environment.
func ResolveGraphLoader() GraphLoader {
	if os.Getenv("LOCAL_DEMO_MODE") != "1" {
		return StubGraphLoader{}
	}
	dir := os.Getenv("LOCAL_DEMO_FIXTURES_DIR")
	if dir == "" {
		dir = "fixtures/local"
	}
	return FileArtifactLoader{BaseDir: dir}
}
