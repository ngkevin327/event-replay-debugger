package sharedschema_test

import (
	"os"
	"path/filepath"
	"testing"

	sharedschema "github.com/replay/platform/packages/shared-schema"
)

func TestValidateValidFixture(t *testing.T) {
	data := readFixture(t, "valid_event.json")
	if err := sharedschema.ValidateCapturedEvent(data); err != nil {
		t.Fatalf("expected valid event: %v", err)
	}
}

func TestValidateInvalidFixtures(t *testing.T) {
	cases := []string{
		"invalid_missing_project.json",
		"invalid_oversized_ref.json",
		"invalid_retry.json",
		"invalid_timestamp.json",
	}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			data := readFixture(t, name)
			if err := sharedschema.ValidateCapturedEvent(data); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestValidateRejectsEmptyObject(t *testing.T) {
	if err := sharedschema.ValidateCapturedEvent([]byte(`{}`)); err == nil {
		t.Fatal("expected validation error for empty object")
	}
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("tests", "fixtures", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return data
}
