package billing

import "testing"

func TestStarterFeatureFlags(t *testing.T) {
	f := StarterFlags()
	if !f.ReplayEnabled || !f.ManualInvoicing {
		t.Fatal("starter flags")
	}
	if f.MaxEventsPerDay < 1_000_000 {
		t.Fatal("event limit")
	}
}

func TestFlagsForTier(t *testing.T) {
	if FlagsForTier("starter").MaxIncidents != 50 {
		t.Fatal("tier mapping")
	}
}
