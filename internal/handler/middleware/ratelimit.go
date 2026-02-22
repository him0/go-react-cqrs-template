package middleware

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitConfig holds the configuration for the rate limiter middleware.
type RateLimitConfig struct {
	// RequestsPerSecond is the rate at which tokens are added to the bucket.
	RequestsPerSecond float64
	// BurstSize is the maximum number of tokens in the bucket.
	BurstSize int
	// CleanupInterval is how often to clean up stale entries.
	CleanupInterval time.Duration
	// StaleEntryTTL is how long after last access before an entry is considered stale.
	StaleEntryTTL time.Duration
	// TrustXForwardedFor controls whether X-Forwarded-For and X-Real-IP headers
	// are trusted for client IP extraction. Set to true only when running behind
	// a trusted reverse proxy. When false (default), only RemoteAddr is used.
	TrustXForwardedFor bool
}

// DefaultRateLimitConfig returns a RateLimitConfig with sensible defaults.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerSecond:  10,
		BurstSize:          20,
		CleanupInterval:    5 * time.Minute,
		StaleEntryTTL:      10 * time.Minute,
		TrustXForwardedFor: false,
	}
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen atomic.Int64 // Unix nanoseconds
}

// RateLimiter provides IP-based rate limiting using the Token Bucket algorithm.
type RateLimiter struct {
	config   RateLimitConfig
	visitors sync.Map
	stopCh   chan struct{}
	stopOnce sync.Once
}

// NewRateLimiter creates a new RateLimiter with the given configuration.
// It starts a background goroutine to clean up stale entries.
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		config: config,
		stopCh: make(chan struct{}),
	}

	go rl.cleanupLoop()

	return rl
}

// Stop stops the background cleanup goroutine.
// It is safe to call Stop multiple times.
func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.stopCh)
	})
}

// getVisitor retrieves or creates a rate limiter for the given IP address.
// Uses LoadOrStore to avoid TOCTOU races with concurrent requests for the same IP.
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	now := time.Now().UnixNano()

	newVis := &visitor{
		limiter: rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize),
	}
	newVis.lastSeen.Store(now)

	actual, loaded := rl.visitors.LoadOrStore(ip, newVis)
	vis := actual.(*visitor)
	if loaded {
		vis.lastSeen.Store(now)
	}
	return vis.limiter
}

// cleanupLoop periodically removes stale visitor entries.
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCh:
			return
		}
	}
}

// cleanup removes visitor entries that have not been seen recently.
func (rl *RateLimiter) cleanup() {
	threshold := time.Now().Add(-rl.config.StaleEntryTTL).UnixNano()
	rl.visitors.Range(func(key, value any) bool {
		vis := value.(*visitor)
		if vis.lastSeen.Load() < threshold {
			rl.visitors.Delete(key)
		}
		return true
	})
}

// Handler returns an HTTP middleware that enforces rate limiting per IP address.
// Requests that exceed the rate limit receive a 429 Too Many Requests response
// with a Retry-After header.
func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r, rl.config.TrustXForwardedFor)
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			retryAfter := int(math.Ceil(1.0 / rl.config.RequestsPerSecond))
			if retryAfter < 1 {
				retryAfter = 1
			}
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// extractIP extracts the client IP address from the request.
// When trustXFF is true, it checks X-Forwarded-For and X-Real-IP headers first.
// Otherwise, it only uses RemoteAddr (safe default for untrusted networks).
func extractIP(r *http.Request, trustXFF bool) string {
	if trustXFF {
		// Check X-Forwarded-For header (first IP in comma-separated list)
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			first, _, _ := strings.Cut(xff, ",")
			first = strings.TrimSpace(first)
			if net.ParseIP(first) != nil {
				return first
			}
		}

		// Check X-Real-IP header
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			xri = strings.TrimSpace(xri)
			if net.ParseIP(xri) != nil {
				return xri
			}
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
