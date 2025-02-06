package jobs

import "context"

// HandlerFunc processes a reconstruction job.
type HandlerFunc func(ctx context.Context, job Job) error

// Dispatcher routes jobs to registered handlers.
type Dispatcher struct {
	handlers map[JobType]HandlerFunc
}

// NewDispatcher creates an empty handler map.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: make(map[JobType]HandlerFunc)}
}

// Register adds a handler for a job type.
func (d *Dispatcher) Register(t JobType, fn HandlerFunc) {
	d.handlers[t] = fn
}

// Dispatch invokes the handler for a job type.
func (d *Dispatcher) Dispatch(ctx context.Context, job Job) error {
	fn := d.handlers[job.Type]
	if fn == nil {
		return nil
	}
	return fn(ctx, job)
}
