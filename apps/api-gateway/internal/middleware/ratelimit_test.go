package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimitReturns429OverThreshold(t *testing.T) {
	rl := NewRateLimiter(1, 1)
	h := RateLimit(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Replay-Key", "test-key")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("first=%d", rec.Code)
	}
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second=%d want 429", rec2.Code)
	}
}
