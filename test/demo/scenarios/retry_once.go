package scenarios

import "errors"

// InjectRetryOnce fails the first handler invocation then succeeds.
type InjectRetryOnce struct {
	seen bool
}

// Run executes fn with a single injected retry.
func (s *InjectRetryOnce) Run(fn func() error) error {
	err := fn()
	if err == nil {
		return nil
	}
	if !s.seen {
		s.seen = true
		return errors.Join(err, ErrRetryScheduled)
	}
	return fn()
}

// ErrRetryScheduled labels a retry-once demo failure.
var ErrRetryScheduled = errors.New("retry_scheduled")
