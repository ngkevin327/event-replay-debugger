package agentgo

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds capture agent runtime settings.
type Config struct {
	ProjectID       string
	APIKey          string
	IngestURL       string
	ControlPlaneURL string
	KafkaBrokers    []string
	ConsumerGroup   string
	TopicAllowlist  []string
	BufferDir       string
	BatchSize       int
	FlushInterval   time.Duration
	CaptureTimeout  time.Duration
	RedactRules     []string
}

// LoadFromEnv reads configuration from environment variables.
func LoadFromEnv() (Config, error) {
	cfg := Config{
		ProjectID:       os.Getenv("REPLAY_PROJECT_ID"),
		APIKey:          os.Getenv("REPLAY_API_KEY"),
		IngestURL:       envOr("REPLAY_INGEST_URL", "http://localhost:8081/v1/ingest/batch"),
		ControlPlaneURL: envOr("REPLAY_CONTROL_PLANE_URL", "http://localhost:8080"),
		ConsumerGroup:   envOr("REPLAY_CONSUMER_GROUP", "replay-agent"),
		BufferDir:       envOr("REPLAY_BUFFER_DIR", os.TempDir()+"/replay-agent"),
		BatchSize:       envIntOr("REPLAY_BATCH_SIZE", 100),
		FlushInterval:   envDurationOr("REPLAY_FLUSH_INTERVAL", 5*time.Second),
		CaptureTimeout:  envDurationOr("REPLAY_CAPTURE_TIMEOUT", 10*time.Millisecond),
	}
	if brokers := os.Getenv("REPLAY_KAFKA_BROKERS"); brokers != "" {
		cfg.KafkaBrokers = strings.Split(brokers, ",")
	} else {
		cfg.KafkaBrokers = []string{"localhost:9092"}
	}
	if topics := os.Getenv("REPLAY_TOPIC_ALLOWLIST"); topics != "" {
		cfg.TopicAllowlist = strings.Split(topics, ",")
	}
	if rules := os.Getenv("REPLAY_REDACT_RULES"); rules != "" {
		cfg.RedactRules = strings.Split(rules, ",")
	}
	if cfg.ProjectID == "" || cfg.APIKey == "" {
		return Config{}, fmt.Errorf("REPLAY_PROJECT_ID and REPLAY_API_KEY required")
	}
	return cfg, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envIntOr(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return def
}

func envDurationOr(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
