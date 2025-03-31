CREATE TYPE replay_status AS ENUM (
    'pending', 'running', 'succeeded', 'failed', 'diverged', 'cancelled'
);

CREATE TABLE replay_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    status replay_status NOT NULL DEFAULT 'pending',
    timing_mode TEXT NOT NULL DEFAULT 'strict',
    divergence_index INTEGER,
    report_uri TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_replay_runs_incident ON replay_runs(incident_id);
