package handlers

import (
	"context"
	"fmt"
)

// StubTimelineLoader returns not-ready until reconstruction artifacts exist.
type StubTimelineLoader struct{}

func (StubTimelineLoader) LoadTimeline(ctx context.Context, incidentID string) ([]byte, int, error) {
	_ = ctx
	_ = incidentID
	return nil, 0, fmt.Errorf("timeline not ready")
}

// StubGraphLoader returns not-ready until graph artifacts exist.
type StubGraphLoader struct{}

func (StubGraphLoader) LoadGraph(ctx context.Context, incidentID string) ([]byte, error) {
	_ = ctx
	_ = incidentID
	return nil, fmt.Errorf("graph not ready")
}
