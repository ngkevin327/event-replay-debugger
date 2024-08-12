package quota_test

import (
	"context"
	"testing"

	"github.com/replay/platform/apps/ingestion/internal/quota"
)

func TestEnforceQuota(t *testing.T) {
	e := quota.NewEnforcer(10)
	if err := e.EnforceQuota(context.Background(), "org", 5); err != nil {
		t.Fatal(err)
	}
	if err := e.EnforceQuota(context.Background(), "org", 6); err == nil {
		t.Fatal("expected quota error")
	}
}
