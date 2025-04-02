package run

import "fmt"

// Status is the replay run lifecycle state.
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusDiverged  Status = "diverged"
	StatusCancelled Status = "cancelled"
)

// StateMachine validates replay transitions.
type StateMachine struct {
	status Status
}

// NewStateMachine creates a machine in pending state.
func NewStateMachine() *StateMachine {
	return &StateMachine{status: StatusPending}
}

// Status returns current status.
func (m *StateMachine) Status() Status {
	return m.status
}

// Transition moves to a new status when allowed.
func (m *StateMachine) Transition(next Status) error {
	if !allowed(m.status, next) {
		return fmt.Errorf("invalid transition %s -> %s", m.status, next)
	}
	m.status = next
	return nil
}

func allowed(from, to Status) bool {
	switch from {
	case StatusPending:
		return to == StatusRunning || to == StatusCancelled
	case StatusRunning:
		return to == StatusSucceeded || to == StatusFailed || to == StatusDiverged || to == StatusCancelled
	default:
		return false
	}
}
