package handlers_test

import (
	"testing"
	"time"

	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func TestRegisterHeartbeatOffline(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	now := time.Now().Add(-20 * time.Minute)
	if store.AgentStatus(&now) != "offline" {
		t.Fatalf("expected offline")
	}
	healthy := time.Now()
	if store.AgentStatus(&healthy) != "healthy" {
		t.Fatalf("expected healthy")
	}
}

func TestAgentLifecycle(t *testing.T) {
	TestRegisterHeartbeatOffline(t)
}
