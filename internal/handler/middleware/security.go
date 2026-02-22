package middleware

import (
	"net/http"
	"os"
)

// defaultCSPPolicy is the default Content-Security-Policy header value.
const defaultCSPPolicy = "default-src 'self'"

// SecurityHeaders returns a middleware that sets security-related HTTP response headers.
// The Content-Security-Policy header can be overridden via the CSP_POLICY environment variable.
func SecurityHeaders(next http.Handler) http.Handler {
	cspPolicy := os.Getenv("CSP_POLICY")
	if cspPolicy == "" {
		cspPolicy = defaultCSPPolicy
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", cspPolicy)
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		next.ServeHTTP(w, r)
	})
}
