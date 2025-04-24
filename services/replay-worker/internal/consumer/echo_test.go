package consumer

import (
	"context"
	"testing"
)

func TestEchoConsumer(t *testing.T) {
	c := &EchoConsumer{}
	if err := c.Validate(context.Background(), "payments", []byte("ok")); err != nil {
		t.Fatal(err)
	}
	if len(c.Seen) != 1 {
		t.Fatal()
	}
}
