package k8s

import (
	"context"
	"testing"
)

func TestProvisionWorker(t *testing.T) {
	p := &JobProvisioner{}
	id, err := p.ProvisionWorker(context.Background(), "rep-1")
	if err != nil || id == "" {
		t.Fatal(err)
	}
}
