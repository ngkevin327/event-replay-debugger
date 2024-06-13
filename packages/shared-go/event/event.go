package event

import "time"

// Outcome describes handler or delivery result for a captured record.
type Outcome string

const (
	OutcomeSuccess        Outcome = "success"
	OutcomeError          Outcome = "error"
	OutcomeRetryScheduled Outcome = "retry_scheduled"
)

// Source identifies consumer vs producer capture.
type Source string

const (
	SourceConsumer Source = "consumer"
	SourceProducer Source = "producer"
)

// CapturedEvent is the canonical capture record aligned with shared-schema JSON.
type CapturedEvent struct {
	EventID             string            `json:"event_id"`
	ProjectID           string            `json:"project_id"`
	CapturedAt          time.Time         `json:"captured_at"`
	Source              Source            `json:"source"`
	Topic               string            `json:"topic"`
	Partition           int32             `json:"partition"`
	Offset              int64             `json:"offset"`
	Timestamp           time.Time         `json:"timestamp"`
	Key                 string            `json:"key,omitempty"`
	Payload             string            `json:"payload,omitempty"`
	PayloadHash         string            `json:"payload_hash,omitempty"`
	PayloadTruncated    bool              `json:"payload_truncated,omitempty"`
	S3URI               string            `json:"s3_uri,omitempty"`
	Headers             map[string]string `json:"headers,omitempty"`
	ConsumerGroup       string            `json:"consumer_group,omitempty"`
	RetryGeneration     int               `json:"retry_generation,omitempty"`
	ProcessingLatencyMs int64             `json:"processing_latency_ms,omitempty"`
	Outcome             Outcome           `json:"outcome,omitempty"`
	Error               *CaptureError     `json:"error,omitempty"`
	ServiceName         string            `json:"service_name,omitempty"`
}

// CaptureError holds minimal error metadata for failed outcomes.
type CaptureError struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}
