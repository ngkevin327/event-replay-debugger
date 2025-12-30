package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter applies token-bucket limits per API key and client IP.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64
	capacity float64
}

type tokenBucket struct {
	tokens float64
	last   time.Time
}

// NewRateLimiter creates a limiter with requests-per-second rate and burst capacity.
func NewRateLimiter(rps float64, burst float64) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     rps,
		capacity: burst,
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.buckets[key]
	now := time.Now()
	if !ok {
		b = &tokenBucket{tokens: rl.capacity, last: now}
		rl.buckets[key] = b
	}
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.capacity {
		b.tokens = rl.capacity
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func clientKey(r *http.Request) string {
	if k := r.Header.Get("X-Replay-Key"); k != "" {
		return "key:" + k
	}
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if idx := strings.Index(ip, ","); idx > 0 {
		ip = strings.TrimSpace(ip[:idx])
	}
	return "ip:" + ip
}

// RateLimit returns middleware that responds 429 when over threshold.
func RateLimit(rl *RateLimiter) func(http.Handler) http.Handler {
	if rl == nil {
		rl = NewRateLimiter(50, 100)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.allow(clientKey(r)) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"rate_limited","message":"too many requests"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
