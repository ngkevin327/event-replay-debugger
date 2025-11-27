package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const signatureHeader = "X-Replay-Signature"

// SignPayload returns HMAC-SHA256 hex for body using secret.
func SignPayload(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature checks header against body and secret.
func VerifySignature(secret string, body []byte, header string) error {
	expected := SignPayload(secret, body)
	if !hmac.Equal([]byte(expected), []byte(strings.TrimSpace(header))) {
		return fmt.Errorf("invalid webhook signature")
	}
	return nil
}
