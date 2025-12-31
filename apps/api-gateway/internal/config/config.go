package config

import (
	"fmt"
	"os"
)

// Config holds api-gateway runtime settings.
type Config struct {
	HTTPAddr      string
	DatabaseURL   string
	LogLevel      string
	JWTSecret     string
	EnvProduction bool
	JWTKeyID      string
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:      envOr("HTTP_ADDR", ":8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		LogLevel:        envOr("LOG_LEVEL", "info"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		EnvProduction:   os.Getenv("ENV") == "production",
		JWTKeyID:        envOr("JWT_KEY_ID", "replay-v1"),
	}
	if cfg.HTTPAddr == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR required")
	}
	return cfg, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
