package report

import "testing"

func TestPersistReport(t *testing.T) {
	_, key, err := PersistReport("rep-1", Report{FirstMismatchIndex: 2})
	if err != nil || key == "" {
		t.Fatal(err)
	}
}
