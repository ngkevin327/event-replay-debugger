package feeder

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

// TimelineEvent is a replayable event entry.
type TimelineEvent struct {
	Topic     string `json:"topic"`
	Partition uint32 `json:"partition"`
	Offset    uint64 `json:"offset"`
	Payload   []byte `json:"-"`
}

// TimelineArtifact is stored timeline JSON.
type TimelineArtifact struct {
	Events []TimelineEvent `json:"events"`
}

// LoadTimeline reads timeline JSON from an artifact URI or local path.
func LoadTimeline(ctx context.Context, uri string) (TimelineArtifact, error) {
	if uri == "" {
		return TimelineArtifact{}, fmt.Errorf("empty timeline uri")
	}
	var body []byte
	var err error
	if len(uri) > 7 && uri[:7] == "file://" {
		body, err = os.ReadFile(uri[7:])
		if err != nil {
			return TimelineArtifact{}, err
		}
	} else {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return TimelineArtifact{}, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return TimelineArtifact{}, err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return TimelineArtifact{}, err
		}
	}
	var art TimelineArtifact
	if err := json.Unmarshal(body, &art); err != nil {
		return TimelineArtifact{}, err
	}
	return art, nil
}

// PublishOrdered sorts and publishes events in partition order.
func PublishOrdered(ctx context.Context, art TimelineArtifact, publish func(ctx context.Context, ev TimelineEvent) error) error {
	events := append([]TimelineEvent(nil), art.Events...)
	sort.Slice(events, func(i, j int) bool {
		if events[i].Topic != events[j].Topic {
			return events[i].Topic < events[j].Topic
		}
		if events[i].Partition != events[j].Partition {
			return events[i].Partition < events[j].Partition
		}
		return events[i].Offset < events[j].Offset
	})
	for _, ev := range events {
		if err := publish(ctx, ev); err != nil {
			return err
		}
	}
	return nil
}
