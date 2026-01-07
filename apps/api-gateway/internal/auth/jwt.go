package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

var ErrInvalidToken = errors.New("invalid token")

// TokenPair holds access and refresh JWT strings.
type TokenPair struct {
	Access  string
	Refresh string
}

// Claims are embedded in access tokens.
type Claims struct {
	UserID string `json:"uid"`
	OrgID  string `json:"oid"`
	jwt.RegisteredClaims
}

// IssueTokens mints short-lived access and long-lived refresh tokens.
func IssueTokens(secret, userID, orgID string) (TokenPair, error) {
	if secret == "" {
		return TokenPair{}, fmt.Errorf("jwt secret required")
	}
	now := time.Now()
	keyID := os.Getenv("JWT_KEY_ID")
	if keyID == "" {
		keyID = "replay-v1"
	}
	accessClaims := Claims{
		UserID: userID,
		OrgID:  orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
			ID:        keyID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	if keyID != "" {
		token.Header["kid"] = keyID
	}
	access, err := token.SignedString([]byte(secret))
	if err != nil {
		return TokenPair{}, fmt.Errorf("sign access: %w", err)
	}

	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		Subject:   userID,
	}
	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		return TokenPair{}, fmt.Errorf("sign refresh: %w", err)
	}
	return TokenPair{Access: access, Refresh: refresh}, nil
}

// ParseAccessToken validates an access JWT and returns claims.
func ParseAccessToken(secret, token string) (Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return Claims{}, ErrInvalidToken
	}
	return *claims, nil
}

// Refresh re-issues tokens from a valid refresh token (rotation).
func Refresh(secret, refreshToken, orgID string) (TokenPair, error) {
	parsed, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return TokenPair{}, ErrInvalidToken
	}
	rc, ok := parsed.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsed.Valid || rc.Subject == "" {
		return TokenPair{}, ErrInvalidToken
	}
	return IssueTokens(secret, rc.Subject, orgID)
}
