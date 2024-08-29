package redact

import (
	"regexp"
	"strings"
)

var defaultRules = []*regexp.Regexp{
	regexp.MustCompile(`"password"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`"ssn"\s*:\s*"[^"]*"`),
	regexp.MustCompile(`\b\d{4}-\d{4}-\d{4}-\d{4}\b`),
}

// Engine applies PII redaction rules to payloads.
type Engine struct {
	rules []*regexp.Regexp
}

// NewEngine builds a redaction engine from rule names.
func NewEngine(ruleNames []string) *Engine {
	if len(ruleNames) == 0 {
		return &Engine{rules: defaultRules}
	}
	return &Engine{rules: defaultRules}
}

// RedactPII masks sensitive fields in a payload string.
func (e *Engine) RedactPII(payload string) string {
	out := payload
	for _, r := range e.rules {
		out = r.ReplaceAllString(out, `"redacted"`)
	}
	return strings.TrimSpace(out)
}
