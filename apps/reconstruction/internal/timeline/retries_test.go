package timeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/replay/platform/apps/reconstruction/internal/fetch"
)

func TestGroupRetries(t *testing.T) {
	rows := []fetch.EventRow{
		{CorrelationID: "c1", PayloadHash: "h1", RetryGeneration: 0},
		{CorrelationID: "c1", PayloadHash: "h1", RetryGeneration: 1},
		{CorrelationID: "c2", PayloadHash: "h2", RetryGeneration: 0},
	}
	chains := GroupRetries(rows)
	if len(chains) != 2 {
		t.Fatalf("got %d", len(chains))
	}
}

func TestRetryChainGolden(t *testing.T) {
	fixture := filepath.Join("..", "..", "..", "test", "golden", "retry_chain_expected.json")
	b, err := os.ReadFile(fixture)
	if err != nil {
		t.Skip("fixture missing")
	}
	var expected struct {
		ChainCount int `json:"chain_count"`
	}
	if err := json.Unmarshal(b, &expected); err != nil {
		t.Fatal(err)
	}
	rows := []fetch.EventRow{
		{CorrelationID: "pay-1", PayloadHash: "abc", RetryGeneration: 0},
		{CorrelationID: "pay-1", PayloadHash: "abc", RetryGeneration: 1},
	}
	if len(GroupRetries(rows)) != expected.ChainCount {
		t.Fatal("chain count mismatch")
	}
}
