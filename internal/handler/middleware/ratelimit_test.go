package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}

func TestRateLimiter_WithinBurst(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 10,
		BurstSize:         5,
		CleanupInterval:   10 * time.Minute,
		StaleEntryTTL:     10 * time.Minute,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	handler := rl.Handler(newTestHandler())

	// All requests within the burst size should succeed
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected status %d, got %d", i, http.StatusOK, w.Code)
		}
	}
}

func TestRateLimiter_ExceedBurst(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 1,
		BurstSize:         3,
		CleanupInterval:   10 * time.Minute,
		StaleEntryTTL:     10 * time.Minute,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	handler := rl.Handler(newTestHandler())

	// Send requests rapidly to exceed the burst
	got429 := false
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code == http.StatusTooManyRequests {
			got429 = true
			// Verify Retry-After header is set
			if retryAfter := w.Header().Get("Retry-After"); retryAfter == "" {
				t.Error("expected Retry-After header to be set on 429 response")
			}
			break
		}
	}

	if !got429 {
		t.Error("expected to receive 429 Too Many Requests after exceeding burst")
	}
}

func TestRateLimiter_DifferentIPsAreIndependent(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 1,
		BurstSize:         2,
		CleanupInterval:   10 * time.Minute,
		StaleEntryTTL:     10 * time.Minute,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	handler := rl.Handler(newTestHandler())

	// Exhaust the burst for IP 1
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	// IP 2 should still be able to make requests
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d for different IP, got %d", http.StatusOK, w.Code)
	}
}

func TestRateLimiter_RetryAfterHeader(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 1,
		BurstSize:         1,
		CleanupInterval:   10 * time.Minute,
		StaleEntryTTL:     10 * time.Minute,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	handler := rl.Handler(newTestHandler())

	// First request should pass
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:8080"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected first request to succeed, got %d", w.Code)
	}

	// Second request should be rate limited
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:8080"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}

	retryAfter := w.Header().Get("Retry-After")
	if retryAfter != "1" {
		t.Errorf("expected Retry-After header to be '1', got '%s'", retryAfter)
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 10,
		BurstSize:         5,
		CleanupInterval:   10 * time.Minute,
		StaleEntryTTL:     0, // Immediately stale
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	// Add a visitor
	rl.getVisitor("192.168.1.1")

	// Verify visitor exists
	if _, ok := rl.visitors.Load("192.168.1.1"); !ok {
		t.Fatal("expected visitor to exist before cleanup")
	}

	// Run cleanup - entry should be removed because StaleEntryTTL is 0
	time.Sleep(1 * time.Millisecond) // Ensure lastSeen is in the past
	rl.cleanup()

	// Verify visitor was removed
	if _, ok := rl.visitors.Load("192.168.1.1"); ok {
		t.Error("expected visitor to be removed after cleanup")
	}
}

func TestExtractIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	ip := extractIP(req)
	if ip != "192.168.1.1" {
		t.Errorf("expected IP '192.168.1.1', got '%s'", ip)
	}
}

func TestExtractIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.RemoteAddr = "192.168.1.1:12345"

	ip := extractIP(req)
	if ip != "10.0.0.1" {
		t.Errorf("expected IP '10.0.0.1', got '%s'", ip)
	}
}

func TestExtractIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.50")
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.RemoteAddr = "192.168.1.1:12345"

	ip := extractIP(req)
	if ip != "203.0.113.50" {
		t.Errorf("expected IP '203.0.113.50', got '%s'", ip)
	}
}
