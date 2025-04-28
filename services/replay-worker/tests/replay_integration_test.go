//go:build integration

package tests

import "testing"

func TestReplayIntegrationRedpanda(t *testing.T) {
	t.Skip("requires docker-compose redpanda")
}
