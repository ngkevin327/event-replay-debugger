-- Membership roles are enforced at the API gateway (admin, member, viewer).
CREATE INDEX IF NOT EXISTS idx_memberships_org_user ON memberships(org_id, user_id);
