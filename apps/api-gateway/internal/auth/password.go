package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLen = 12
	bcryptCost     = 12
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 12 characters")
	ErrPasswordMismatch = errors.New("password mismatch")
)

// HashPassword returns a bcrypt hash for storage.
func HashPassword(plain string) (string, error) {
	if len(plain) < minPasswordLen {
		return "", ErrPasswordTooShort
	}
	out, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(out), nil
}

// VerifyPassword compares plain text against a stored bcrypt hash.
func VerifyPassword(hash, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return ErrPasswordMismatch
	}
	return nil
}
