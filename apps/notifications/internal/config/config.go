package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds notification worker settings.
type Config struct {
	WebhookTimeout time.Duration
	MaxRetries     int
	EventBusURL    string
}

// Load reads environment configuration with defaults.
func Load() Config {
	timeout := 10 * time.Second
	if v := os.Getenv("WEBHOOK_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			timeout = d
		}
	}
	retries := 3
	if v := os.Getenv("WEBHOOK_MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			retries = n
		}
	}
	bus := os.Getenv("EVENT_BUS_URL")
	if bus == "" {
		bus = "memory://events"
	}
	return Config{
		WebhookTimeout: timeout,
		MaxRetries:     retries,
		EventBusURL:    bus,
	}
}
