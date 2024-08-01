package payload

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/replay/platform/packages/shared-go/event"
)

const TruncateAt256KB = 256 * 1024

// ProcessPayload hashes, optionally truncates, and returns processed bytes.
func ProcessPayload(raw []byte) (data []byte, hash string, truncated bool) {
	hash = sha256Hex(raw)
	if len(raw) > TruncateAt256KB {
		raw = raw[:TruncateAt256KB]
		truncated = true
	}
	return raw, hash, truncated
}

// ApplyToEvent updates event payload fields after processing.
func ApplyToEvent(ev *event.CapturedEvent, data []byte, hash string, truncated bool) {
	ev.Payload = string(data)
	ev.PayloadHash = hash
	ev.PayloadTruncated = truncated
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
