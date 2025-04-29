CREATE TABLE replay_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    replay_id UUID NOT NULL REFERENCES replay_runs(id) ON DELETE CASCADE,
    topic TEXT NOT NULL,
    partition INTEGER NOT NULL,
    offset BIGINT NOT NULL,
    outcome TEXT NOT NULL,
    checkpoint_index INTEGER NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_replay_events_replay ON replay_events(replay_id);
