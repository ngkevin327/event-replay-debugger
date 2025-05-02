package divergence

import "github.com/replay/platform/services/replay-worker/internal/record"

// Mismatch describes the first divergent checkpoint.
type Mismatch struct {
	Index  int
	Expect string
	Actual string
}

// CompareChain compares expected and actual outcome chains.
func CompareChain(expected, actual []record.Outcome) (bool, *Mismatch) {
	for i := range expected {
		if i >= len(actual) {
			return false, &Mismatch{Index: i, Expect: expected[i].Result, Actual: ""}
		}
		if expected[i].Result != actual[i].Result {
			return false, &Mismatch{Index: i, Expect: expected[i].Result, Actual: actual[i].Result}
		}
	}
	return true, nil
}

// FirstMismatch returns the first mismatch or nil when chains match.
func FirstMismatch(expected, actual []record.Outcome) *Mismatch {
	ok, mm := CompareChain(expected, actual)
	if ok {
		return nil
	}
	return mm
}
