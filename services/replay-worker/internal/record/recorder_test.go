package record

import "testing"

func TestRecorder(t *testing.T) {
	r := &Recorder{}
	r.RecordOutcome(Outcome{Topic: "payments", Result: "ok"})
	if len(r.Checkpoints) != 1 {
		t.Fatal()
	}
}
