package kafka

import "github.com/replay/platform/packages/shared-go/event"

// RetryGeneration tracks handler retry attempts for a message.
type RetryGeneration int

// LabelOutcome sets outcome and retry generation on a captured event.
func LabelOutcome(ev *event.CapturedEvent, gen RetryGeneration, outcome event.Outcome, err error) {
	ev.RetryGeneration = int(gen)
	ev.Outcome = outcome
	if err != nil {
		ev.Error = &event.CaptureError{Type: "handler", Message: err.Error()}
		if outcome == "" {
			ev.Outcome = event.OutcomeError
		}
	}
	if outcome == event.OutcomeRetryScheduled {
		ev.Outcome = event.OutcomeRetryScheduled
	}
}

// NextRetryGeneration increments generation after a retry_scheduled outcome.
func NextRetryGeneration(current RetryGeneration) RetryGeneration {
	return current + 1
}
