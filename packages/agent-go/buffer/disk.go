package buffer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/replay/platform/packages/shared-go/event"
)

// DiskBuffer persists captured events with drop-oldest overflow.
type DiskBuffer struct {
	dir     string
	max     int
	mu      sync.Mutex
	files   []string
}

// NewDiskBuffer opens a directory-backed buffer.
func NewDiskBuffer(dir string, max int) (*DiskBuffer, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	if max <= 0 {
		max = 1000
	}
	return &DiskBuffer{dir: dir, max: max}, nil
}

// Append stores an event, dropping oldest when full.
func (b *DiskBuffer) Append(ev event.CapturedEvent) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.files) >= b.max {
		if err := b.dropOldest(); err != nil {
			return err
		}
	}
	name := filepath.Join(b.dir, ev.EventID+".json")
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if err := os.WriteFile(name, data, 0o644); err != nil {
		return err
	}
	b.files = append(b.files, name)
	return nil
}

// DropOldest removes the oldest buffered file.
func (b *DiskBuffer) DropOldest() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.dropOldest()
}

func (b *DiskBuffer) dropOldest() error {
	if len(b.files) == 0 {
		return nil
	}
	oldest := b.files[0]
	b.files = b.files[1:]
	return os.Remove(oldest)
}

// Drain reads all buffered events and clears the buffer.
func (b *DiskBuffer) Drain() ([]event.CapturedEvent, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []event.CapturedEvent
	for _, f := range b.files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		var ev event.CapturedEvent
		if err := json.Unmarshal(data, &ev); err != nil {
			return nil, err
		}
		out = append(out, ev)
		_ = os.Remove(f)
	}
	b.files = nil
	return out, nil
}
