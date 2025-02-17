package timeline

import (
	"sort"

	"github.com/replay/platform/apps/reconstruction/internal/fetch"
)

// OrderByPartitionOffset sorts events by topic, partition, then offset.
func OrderByPartitionOffset(rows []fetch.EventRow) []fetch.EventRow {
	out := append([]fetch.EventRow(nil), rows...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Topic != out[j].Topic {
			return out[i].Topic < out[j].Topic
		}
		if out[i].Partition != out[j].Partition {
			return out[i].Partition < out[j].Partition
		}
		return out[i].Offset < out[j].Offset
	})
	return out
}

// ArrivalOrderIndex annotates capture arrival sequence across partitions.
func ArrivalOrderIndex(rows []fetch.EventRow) []fetch.EventRow {
	out := append([]fetch.EventRow(nil), rows...)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].CapturedAt.Before(out[j].CapturedAt)
	})
	for i := range out {
		out[i].ArrivalIndex = i
	}
	return out
}
