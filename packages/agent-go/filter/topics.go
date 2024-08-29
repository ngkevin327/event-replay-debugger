package filter

import "strings"

// Allowlist restricts capture to configured topics.
type Allowlist struct {
	topics map[string]struct{}
}

// NewAllowlist builds a topic allowlist from names.
func NewAllowlist(topics []string) *Allowlist {
	m := make(map[string]struct{}, len(topics))
	for _, t := range topics {
		t = strings.TrimSpace(t)
		if t != "" {
			m[t] = struct{}{}
		}
	}
	return &Allowlist{topics: m}
}

// ShouldCapture returns true when topic is allowed.
func (a *Allowlist) ShouldCapture(topic string) bool {
	if a == nil || len(a.topics) == 0 {
		return true
	}
	_, ok := a.topics[topic]
	return ok
}
