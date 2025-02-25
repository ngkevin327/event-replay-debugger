package timeline

import "github.com/replay/platform/apps/reconstruction/internal/fetch"

// GapMarker describes a missing offset range in a partition.
type GapMarker struct {
	Topic     string `json:"topic"`
	Partition uint32 `json:"partition"`
	Start     uint64 `json:"start_offset"`
	End       uint64 `json:"end_offset"`
}

// DetectGaps finds offset discontinuities per topic/partition.
func DetectGaps(rows []fetch.EventRow) []GapMarker {
	type key struct {
		topic string
		part  uint32
	}
	latest := map[key]uint64{}
	var gaps []GapMarker
	ordered := OrderByPartitionOffset(rows)
	for _, row := range ordered {
		k := key{row.Topic, row.Partition}
		if prev, ok := latest[k]; ok && row.Offset > prev+1 {
			gaps = append(gaps, GapMarker{
				Topic: row.Topic, Partition: row.Partition,
				Start: prev + 1, End: row.Offset - 1,
			})
		}
		latest[k] = row.Offset
	}
	return gaps
}
