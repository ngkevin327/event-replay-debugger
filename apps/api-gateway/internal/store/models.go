package store

import (
	"time"
)

// Organization is a tenant container.
type Organization struct {
	ID        string
	Name      string
	PlanTier  string
	CreatedAt time.Time
}

// Project belongs to an organization.
type Project struct {
	ID        string
	OrgID     string
	Name      string
	CreatedAt time.Time
}

// User is a platform login identity.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// MembershipRole mirrors Postgres enum.
type MembershipRole string

const (
	RoleAdmin  MembershipRole = "admin"
	RoleMember MembershipRole = "member"
	RoleViewer MembershipRole = "viewer"
)
