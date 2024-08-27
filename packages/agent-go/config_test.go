package agentgo_test

import (
	"os"
	"testing"

	agentgo "github.com/replay/platform/packages/agent-go"
)

func TestConfigLoad(t *testing.T) {
	t.Setenv("REPLAY_PROJECT_ID", "proj-1")
	t.Setenv("REPLAY_API_KEY", "rk_live_test")
	cfg, err := agentgo.LoadFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ProjectID != "proj-1" || cfg.APIKey == "" {
		t.Fatalf("cfg %+v", cfg)
	}
}

func TestConfigLoadMissingRequired(t *testing.T) {
	os.Unsetenv("REPLAY_PROJECT_ID")
	os.Unsetenv("REPLAY_API_KEY")
	if _, err := agentgo.LoadFromEnv(); err == nil {
		t.Fatal("expected error")
	}
}
