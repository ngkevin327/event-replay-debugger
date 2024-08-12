CREATE TABLE plan_quotas (
    plan_tier TEXT PRIMARY KEY,
    daily_events BIGINT NOT NULL,
    max_payload_bytes INTEGER NOT NULL DEFAULT 262144,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO plan_quotas (plan_tier, daily_events, max_payload_bytes)
VALUES ('starter', 1000000, 262144)
ON CONFLICT (plan_tier) DO NOTHING;
