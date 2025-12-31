package server

import "net/http"

// ApplyTLSHeaders sets HSTS and secure transport headers for production.
func ApplyTLSHeaders(w http.ResponseWriter, production bool) {
	if !production {
		return
	}
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

// SecureCookieFlags returns cookie flags for session tokens in production.
func SecureCookieFlags(production bool) (secure, httpOnly bool, sameSite string) {
	if production {
		return true, true, "Strict"
	}
	return false, true, "Lax"
}
