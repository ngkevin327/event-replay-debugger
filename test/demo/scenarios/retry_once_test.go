package scenarios_test

import (
	"errors"
	"testing"

	"github.com/replay/platform/test/demo/scenarios"
)

func TestInjectRetryOnce(t *testing.T) {
	s := &scenarios.InjectRetryOnce{}
	calls := 0
	err := s.Run(func() error {
		calls++
		if calls == 1 {
			return errors.New("fail")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("calls %d", calls)
	}
}
