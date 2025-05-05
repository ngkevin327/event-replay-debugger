package report

import (
	"encoding/json"
	"fmt"
)

// Report is the divergence artifact persisted to object storage.
type Report struct {
	ReplayID            string `json:"replay_id"`
	FirstMismatchIndex  int    `json:"first_mismatch_index"`
	Expected            string `json:"expected"`
	Actual              string `json:"actual"`
}

// PersistReport serializes a report and returns an S3 key.
func PersistReport(replayID string, r Report) ([]byte, string, error) {
	r.ReplayID = replayID
	body, err := json.Marshal(r)
	if err != nil {
		return nil, "", err
	}
	key := fmt.Sprintf("reports/%s/divergence.json", replayID)
	return body, key, nil
}
