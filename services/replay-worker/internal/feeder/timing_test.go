package feeder

import (
	"testing"
	"time"
)

func TestCompressedTiming(t *testing.T) {
	prev := time.Now()
	CompressedTiming(prev, prev.Add(10*time.Millisecond), 10)
}
