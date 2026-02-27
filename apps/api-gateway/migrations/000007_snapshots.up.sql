CREATE TABLE offset_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    consumer_group TEXT NOT NULL,
    topic TEXT NOT NULL,
    partition INTEGER NOT NULL,
    "offset" BIGINT NOT NULL,
    captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_offset_snapshots_incident ON offset_snapshots(incident_id);
