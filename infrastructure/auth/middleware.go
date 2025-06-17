package auth

import (
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// AuthMiddleware is a middleware for handling authentication and authorization
type AuthMiddleware struct {
	jwtService  *JWTService
	oidcService *OIDCService
	logger      *zap.Logger
	tracer      trace.Tracer
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtService *JWTService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
		tracer:     otel.Tracer("infrastructure.auth.middleware"),
	}
}

// NewAuthMiddlewareWithOIDC creates a new auth middleware with OIDC support
func NewAuthMiddlewareWithOIDC(jwtService *JWTService, oidcService *OIDCService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:  jwtService,
		oidcService: oidcService,
		logger:      logger,
		tracer:      otel.Tracer("infrastructure.auth.middleware"),
	}
}

// Middleware is the HTTP middleware function
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := m.tracer.Start(r.Context(), "AuthMiddleware.Middleware")
		defer span.End()

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No token provided, continue as unauthenticated
			m.logger.Debug("No Authorization header provided")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Debug("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		span.SetAttributes(attribute.String("token.length", fmt.Sprintf("%d", len(tokenString))))

		var claims *Claims
		var err error

		// Try OIDC validation first if available
		if m.oidcService != nil {
			claims, err = m.oidcService.ValidateToken(ctx, tokenString)
			if err != nil {
				m.logger.Debug("OIDC validation failed, trying JWT", zap.Error(err))
				// Fall back to JWT validation
				claims, err = m.jwtService.ValidateToken(tokenString)
				if err != nil {
					m.logger.Debug("Invalid token", zap.Error(err))
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
			}
		} else {
			// Use JWT validation
			claims, err = m.jwtService.ValidateToken(tokenString)
			if err != nil {
				m.logger.Debug("Invalid token", zap.Error(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		}

		// Add user info to context
		ctx = WithUserID(ctx, claims.UserID)
		ctx = WithUserRoles(ctx, claims.Roles)

		span.SetAttributes(attribute.String("user.id", claims.UserID))

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
