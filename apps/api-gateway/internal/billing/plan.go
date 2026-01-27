package billing

// PlanTier identifies subscription level.
type PlanTier string

const (
	PlanStarter PlanTier = "starter"
	PlanPro     PlanTier = "pro"
)

// FeatureFlags describes entitlements for a plan.
type FeatureFlags struct {
	MaxEventsPerDay   int64
	MaxIncidents      int
	ReplayEnabled     bool
	WebhooksEnabled   bool
	RetentionDays     int
	ManualInvoicing   bool
}

// StarterFlags returns MVP starter plan limits.
func StarterFlags() FeatureFlags {
	return FeatureFlags{
		MaxEventsPerDay: 1_000_000,
		MaxIncidents:    50,
		ReplayEnabled:   true,
		WebhooksEnabled: true,
		RetentionDays:   30,
		ManualInvoicing: true,
	}
}

// FlagsForTier resolves feature flags by plan tier string.
func FlagsForTier(tier string) FeatureFlags {
	switch PlanTier(tier) {
	case PlanStarter:
		return StarterFlags()
	default:
		return StarterFlags()
	}
}

// InvoiceHook documents manual invoicing integration point (MVP stub).
func InvoiceHook(orgID, planTier string, usageEvents int64) error {
	_ = orgID
	_ = planTier
	_ = usageEvents
	return nil
}
