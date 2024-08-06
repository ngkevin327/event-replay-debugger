CREATE DATABASE IF NOT EXISTS replay;

-- Events index for ingest pipeline (Stage 2).
CREATE TABLE IF NOT EXISTS replay.events
(
    event_id UUID,
    project_id UUID,
    org_id UUID,
    captured_at DateTime64(3, 'UTC'),
    source LowCardinality(String),
    topic LowCardinality(String),
    partition UInt32,
    offset UInt64,
    timestamp DateTime64(3, 'UTC'),
    consumer_group String,
    retry_generation UInt16,
    outcome LowCardinality(String),
    payload_hash FixedString(64),
    payload_truncated UInt8,
    s3_uri String,
    correlation_id String,
    inserted_at DateTime64(3, 'UTC') DEFAULT now64(3)
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(captured_at)
ORDER BY (project_id, topic, partition, offset, captured_at);
