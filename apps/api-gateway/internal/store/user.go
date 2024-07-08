package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// CreateUser inserts a new user with email and password hash.
func (s *Store) CreateUser(ctx context.Context, email, passwordHash string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2)
		 RETURNING id, email, password_hash, created_at`,
		email, passwordHash,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

// GetUserByEmail loads a user by unique email.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, err
	}
	return u, err
}

// CreateMembership links a user to an organization.
func (s *Store) CreateMembership(ctx context.Context, orgID, userID string, role MembershipRole) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO memberships (org_id, user_id, role) VALUES ($1, $2, $3)`,
		orgID, userID, role,
	)
	return err
}
