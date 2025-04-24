package consumer

import "context"

// EchoConsumer validates replayed messages in the MVP demo path.
type EchoConsumer struct {
	Seen []string
}

// Validate records a message key for later comparison.
func (c *EchoConsumer) Validate(ctx context.Context, topic string, payload []byte) error {
	_ = ctx
	c.Seen = append(c.Seen, topic+":"+string(payload))
	return nil
}
