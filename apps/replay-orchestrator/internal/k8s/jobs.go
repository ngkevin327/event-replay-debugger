package k8s

import "context"

// JobProvisioner manages ephemeral replay worker jobs.
type JobProvisioner struct{}

// ProvisionWorker schedules a worker job for a replay run.
func (p *JobProvisioner) ProvisionWorker(ctx context.Context, replayID string) (string, error) {
	_ = ctx
	return "job-" + replayID, nil
}

// TearDown removes a worker job after completion.
func (p *JobProvisioner) TearDown(ctx context.Context, jobID string) error {
	_ = ctx
	_ = jobID
	return nil
}
