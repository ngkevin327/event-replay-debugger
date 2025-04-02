package run

import "testing"

func TestMachineTerminalTransitions(t *testing.T) {
	m := NewStateMachine()
	if err := m.Transition(StatusRunning); err != nil {
		t.Fatal(err)
	}
	if err := m.Transition(StatusDiverged); err != nil {
		t.Fatal(err)
	}
	if m.Status() != StatusDiverged {
		t.Fatalf("got %s", m.Status())
	}
}

func TestMachineCancelFromPending(t *testing.T) {
	m := NewStateMachine()
	if err := m.Transition(StatusCancelled); err != nil {
		t.Fatal(err)
	}
}
