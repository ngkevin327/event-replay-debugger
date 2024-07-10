package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const keyPrefix = "rk_live_"

var (
	ErrInvalidAPIKey = errors.New("invalid api key")
	ErrInvalidScope  = errors.New("invalid scope")
)

// GenerateKey returns a new plaintext API key and its storage prefix/hash.
func GenerateKey() (plain, prefix, hash string, err error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", "", "", fmt.Errorf("rand: %w", err)
	}
	plain = keyPrefix + hex.EncodeToString(buf)
	prefix = plain[:16]
	hash, err = HashKey(plain)
	if err != nil {
		return "", "", "", err
	}
	return plain, prefix, hash, nil
}

// HashKey hashes a plaintext API key for storage.
func HashKey(plain string) (string, error) {
	if !strings.HasPrefix(plain, keyPrefix) {
		return "", ErrInvalidAPIKey
	}
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:]), nil
}

// ValidateScope ensures the key scope set includes required scope.
func ValidateScope(scopes []string, required string) bool {
	for _, s := range scopes {
		if s == required {
			return true
		}
	}
	return false
}
