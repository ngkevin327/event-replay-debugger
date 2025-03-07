package cache

import (
	"testing"
	"time"
)

func TestTimelineCache(t *testing.T) {
	c := NewTimelineCache(time.Minute)
	c.Set("inc-1", []byte("{}"), 1)
	b, v, ok := c.Get("inc-1")
	if !ok || v != 1 || len(b) == 0 {
		t.Fatal()
	}
}
