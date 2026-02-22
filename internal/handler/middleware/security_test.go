package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	// Create a simple handler that the middleware wraps
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := SecurityHeaders(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":       "DENY",
		"X-XSS-Protection":      "1; mode=block",
		"Referrer-Policy":       "strict-origin-when-cross-origin",
		"Content-Security-Policy": "default-src 'self'",
		"Permissions-Policy":     "camera=(), microphone=(), geolocation=()",
	}

	for header, expected := range expectedHeaders {
		actual := rec.Header().Get(header)
		if actual != expected {
			t.Errorf("header %s: expected %q, got %q", header, expected, actual)
		}
	}
}

func TestSecurityHeaders_CSPOverride(t *testing.T) {
	// Override CSP via environment variable
	t.Setenv("CSP_POLICY", "default-src 'self' 'unsafe-inline'")

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := SecurityHeaders(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := "default-src 'self' 'unsafe-inline'"
	actual := rec.Header().Get("Content-Security-Policy")
	if actual != expected {
		t.Errorf("Content-Security-Policy: expected %q, got %q", expected, actual)
	}
}

func TestSecurityHeaders_PassesRequestThrough(t *testing.T) {
	called := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handler := SecurityHeaders(inner)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("inner handler was not called")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
