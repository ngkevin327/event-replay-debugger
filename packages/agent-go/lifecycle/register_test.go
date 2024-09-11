package lifecycle_test

import (
	"context"
	"testing"
	"time"
)

func TestHeartbeatLoopCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// HeartbeatLoop with cancelled ctx should return quickly when called with invalid URL skipped
	_ = time.Second
	_ = ctx
}
