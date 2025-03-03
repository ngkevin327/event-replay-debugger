package timeline

import (
	"testing"
)

func TestBuildTimelineArtifact(t *testing.T) {
	b, err := BuildTimelineArtifact("inc-1", []EventEntry{{EventID: "e1"}}, nil, nil)
	if err != nil || len(b) == 0 {
		t.Fatal(err)
	}
}
