package quota

import "testing"

func TestConcurrentReplayLimit(t *testing.T) {
	l := NewLimiter(1)
	l.Acquire("proj-1")
	if l.ConcurrentReplayLimit("proj-1") {
		t.Fatal("expected limit reached")
	}
	l.Release("proj-1")
	if !l.ConcurrentReplayLimit("proj-1") {
		t.Fatal()
	}
}
