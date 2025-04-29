package record

// Outcome describes replay validation result at a checkpoint.
type Outcome struct {
	Topic     string
	Partition int
	Offset    int64
	Result    string
}

// Recorder stores replay outcomes per offset checkpoint.
type Recorder struct {
	Checkpoints []Outcome
}

// Checkpoint appends an outcome row.
func (r *Recorder) Checkpoint(o Outcome) {
	r.Checkpoints = append(r.Checkpoints, o)
}

// RecordOutcome records a validation result.
func (r *Recorder) RecordOutcome(o Outcome) {
	r.Checkpoint(o)
}
