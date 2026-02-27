package config

import "os"

// Config holds ingestion service settings.
type Config struct {
	HTTPAddr      string
	DatabaseURL   string
	S3Endpoint    string
	S3AccessKey   string
	S3SecretKey   string
	S3Bucket      string
	BucketPrefix  string
	ClickHouseURL string
}

// Load reads configuration from environment.
func Load() Config {
	return Config{
		HTTPAddr:      envOr("HTTP_ADDR", ":8081"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		S3Endpoint:    envOr("S3_ENDPOINT", "http://localhost:19000"),
		S3AccessKey:   envOr("S3_ACCESS_KEY", "minio"),
		S3SecretKey:   envOr("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:      envOr("S3_BUCKET", "replay-payloads"),
		BucketPrefix:  envOr("S3_BUCKET_PREFIX", "payloads"),
		ClickHouseURL: envOr("CLICKHOUSE_URL", "http://localhost:8123"),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
