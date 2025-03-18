package snapshot

import "time"

// OffsetSnapshot records consumer positions at incident end.
type OffsetSnapshot struct {
	ConsumerGroup string    `json:"consumer_group"`
	Topic         string    `json:"topic"`
	Partition     uint32    `json:"partition"`
	Offset        uint64    `json:"offset"`
	CapturedAt    time.Time `json:"captured_at"`
}

// CaptureSnapshots builds snapshots from latest event offsets per group.
func CaptureSnapshots(events []struct {
	ConsumerGroup string
	Topic         string
	Partition     uint32
	Offset        uint64
}) []OffsetSnapshot {
	type key struct {
		group, topic string
		part         uint32
	}
	latest := map[key]uint64{}
	for _, ev := range events {
		k := key{ev.ConsumerGroup, ev.Topic, ev.Partition}
		if ev.Offset >= latest[k] {
			latest[k] = ev.Offset
		}
	}
	out := make([]OffsetSnapshot, 0, len(latest))
	now := time.Now().UTC()
	for k, off := range latest {
		out = append(out, OffsetSnapshot{
			ConsumerGroup: k.group,
			Topic:         k.topic,
			Partition:     k.part,
			Offset:        off,
			CapturedAt:    now,
		})
	}
	return out
}
