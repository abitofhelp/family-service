// Copyright (c) 2025 A Bit of Help, Inc.

package security

import (
	"context"
	"net/http"

	"github.com/abitofhelp/servicelib/logging"
)

// SecurityHeadersConfig defines the configuration for security headers
type SecurityHeadersConfig struct {
	// Content Security Policy
	ContentSecurityPolicy string

	// X-XSS-Protection
	XSSProtection string

	// X-Content-Type-Options
	ContentTypeOptions string

	// X-Frame-Options
	FrameOptions string

	// Referrer-Policy
	ReferrerPolicy string

	// Strict-Transport-Security
	StrictTransportSecurity string

	// Feature-Policy
	FeaturePolicy string

	// Permissions-Policy
	PermissionsPolicy string
}

// DefaultSecurityHeadersConfig returns a default configuration for security headers
func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		ContentSecurityPolicy:   "default-src 'self'; script-src 'self'; object-src 'none'; img-src 'self'; media-src 'self'; frame-src 'none'; font-src 'self'; connect-src 'self'",
		XSSProtection:           "1; mode=block",
		ContentTypeOptions:      "nosniff",
		FrameOptions:            "DENY",
		ReferrerPolicy:          "strict-origin-when-cross-origin",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains",
		FeaturePolicy:           "camera 'none'; microphone 'none'; geolocation 'none'",
		PermissionsPolicy:       "camera=(), microphone=(), geolocation=()",
	}
}

// SecurityHeadersMiddleware is a middleware that adds security headers to HTTP responses
type SecurityHeadersMiddleware struct {
	config SecurityHeadersConfig
	logger *logging.ContextLogger
}

// NewSecurityHeadersMiddleware creates a new SecurityHeadersMiddleware
func NewSecurityHeadersMiddleware(config SecurityHeadersConfig, logger *logging.ContextLogger) *SecurityHeadersMiddleware {
	if logger == nil {
		panic("logger cannot be nil")
	}

	return &SecurityHeadersMiddleware{
		config: config,
		logger: logger,
	}
}

// Middleware returns an http.Handler middleware function
func (m *SecurityHeadersMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m.logger.Debug(ctx, "Adding security headers to response")

		// Add security headers
		if m.config.ContentSecurityPolicy != "" {
			w.Header().Set("Content-Security-Policy", m.config.ContentSecurityPolicy)
		}

		if m.config.XSSProtection != "" {
			w.Header().Set("X-XSS-Protection", m.config.XSSProtection)
		}

		if m.config.ContentTypeOptions != "" {
			w.Header().Set("X-Content-Type-Options", m.config.ContentTypeOptions)
		}

		if m.config.FrameOptions != "" {
			w.Header().Set("X-Frame-Options", m.config.FrameOptions)
		}

		if m.config.ReferrerPolicy != "" {
			w.Header().Set("Referrer-Policy", m.config.ReferrerPolicy)
		}

		if m.config.StrictTransportSecurity != "" {
			w.Header().Set("Strict-Transport-Security", m.config.StrictTransportSecurity)
		}

		if m.config.FeaturePolicy != "" {
			w.Header().Set("Feature-Policy", m.config.FeaturePolicy)
		}

		if m.config.PermissionsPolicy != "" {
			w.Header().Set("Permissions-Policy", m.config.PermissionsPolicy)
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// AddSecurityHeaders adds security headers to an http.ResponseWriter
func AddSecurityHeaders(ctx context.Context, w http.ResponseWriter, config SecurityHeadersConfig, logger *logging.ContextLogger) {
	logger.Debug(ctx, "Adding security headers to response")

	// Add security headers
	if config.ContentSecurityPolicy != "" {
		w.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)
	}

	if config.XSSProtection != "" {
		w.Header().Set("X-XSS-Protection", config.XSSProtection)
	}

	if config.ContentTypeOptions != "" {
		w.Header().Set("X-Content-Type-Options", config.ContentTypeOptions)
	}

	if config.FrameOptions != "" {
		w.Header().Set("X-Frame-Options", config.FrameOptions)
	}

	if config.ReferrerPolicy != "" {
		w.Header().Set("Referrer-Policy", config.ReferrerPolicy)
	}

	if config.StrictTransportSecurity != "" {
		w.Header().Set("Strict-Transport-Security", config.StrictTransportSecurity)
	}

	if config.FeaturePolicy != "" {
		w.Header().Set("Feature-Policy", config.FeaturePolicy)
	}

	if config.PermissionsPolicy != "" {
		w.Header().Set("Permissions-Policy", config.PermissionsPolicy)
	}

	logger.Debug(ctx, "Security headers added to response")
}
