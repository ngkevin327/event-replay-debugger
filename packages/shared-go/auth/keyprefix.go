package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const livePrefix = "rk_live_"

var ErrInvalidKey = errors.New("invalid api key format")

// ParseKeyPrefix returns the 16-character storage prefix from a plaintext key.
func ParseKeyPrefix(plain string) (string, error) {
	if !strings.HasPrefix(plain, livePrefix) || len(plain) < 16 {
		return "", ErrInvalidKey
	}
	return plain[:16], nil
}

// HashPrefix returns the SHA-256 hex digest used for api_keys.key_hash lookup.
func HashPrefix(plain string) (string, error) {
	if !strings.HasPrefix(plain, livePrefix) {
		return "", ErrInvalidKey
	}
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:]), nil
}
