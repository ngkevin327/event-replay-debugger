CREATE TABLE notification_preferences (
    project_id UUID PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    webhook_url TEXT,
    webhook_secret TEXT,
    email_enabled BOOLEAN NOT NULL DEFAULT true,
    email_recipients TEXT[] NOT NULL DEFAULT '{}',
    notify_incident_ready BOOLEAN NOT NULL DEFAULT true,
    notify_replay_completed BOOLEAN NOT NULL DEFAULT true,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
