package agentgo

import (
	"context"
	"sync"
)

// Agent coordinates capture, buffering, and shipping.
type Agent struct {
	cfg    Config
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New creates an agent from configuration.
func New(cfg Config) *Agent {
	return &Agent{cfg: cfg}
}

// Start launches background flush and lifecycle loops.
func (a *Agent) Start(ctx context.Context) error {
	ctx, a.cancel = context.WithCancel(ctx)
	_ = ctx
	return nil
}

// Stop shuts down background work gracefully.
func (a *Agent) Stop() {
	if a.cancel != nil {
		a.cancel()
	}
	a.wg.Wait()
}

// Config returns the agent configuration.
func (a *Agent) Config() Config {
	return a.cfg
}
