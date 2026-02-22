package middleware

import (
	"net"
	"net/http"
	"sync"
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
}

// DefaultRateLimitConfig returns a RateLimitConfig with sensible defaults.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   5 * time.Minute,
		StaleEntryTTL:     10 * time.Minute,
	}
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter provides IP-based rate limiting using the Token Bucket algorithm.
type RateLimiter struct {
	config   RateLimitConfig
	visitors sync.Map
	stopCh   chan struct{}
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
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

// getVisitor retrieves or creates a rate limiter for the given IP address.
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	if v, ok := rl.visitors.Load(ip); ok {
		vis := v.(*visitor)
		vis.lastSeen = time.Now()
		return vis.limiter
	}

	limiter := rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize)
	rl.visitors.Store(ip, &visitor{
		limiter:  limiter,
		lastSeen: time.Now(),
	})
	return limiter
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
	threshold := time.Now().Add(-rl.config.StaleEntryTTL)
	rl.visitors.Range(func(key, value any) bool {
		vis := value.(*visitor)
		if vis.lastSeen.Before(threshold) {
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
		ip := extractIP(r)
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// extractIP extracts the client IP address from the request.
// It checks X-Forwarded-For and X-Real-IP headers first, then falls back
// to RemoteAddr.
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For may contain multiple IPs; use the first one
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		// It might not have a port
		if netIP := net.ParseIP(xff); netIP != nil {
			return xff
		}
		// Take the first IP from a comma-separated list
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				candidate := xff[:i]
				if netIP := net.ParseIP(candidate); netIP != nil {
					return candidate
				}
				break
			}
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
