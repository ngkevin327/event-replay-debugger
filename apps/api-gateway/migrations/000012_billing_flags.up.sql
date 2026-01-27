ALTER TABLE projects ADD COLUMN IF NOT EXISTS plan_tier TEXT NOT NULL DEFAULT 'starter';

CREATE INDEX IF NOT EXISTS idx_projects_plan_tier ON projects(plan_tier);
