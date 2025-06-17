package auth

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// contextKey is a private type for context keys
type contextKey int

const (
	userIDKey contextKey = iota
	userRolesKey
)

// AuthorizationService implements the ports.AuthorizationService interface
type AuthorizationService struct {
	logger *zap.Logger
	tracer trace.Tracer
}

// NewAuthorizationService creates a new authorization service
func NewAuthorizationService(logger *zap.Logger) *AuthorizationService {
	return &AuthorizationService{
		logger: logger,
		tracer: otel.Tracer("infrastructure.auth.service"),
	}
}

// IsAuthorized checks if the user is authorized to perform the operation
func (s *AuthorizationService) IsAuthorized(ctx context.Context, operation string) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "AuthorizationService.IsAuthorized")
	defer span.End()

	span.SetAttributes(attribute.String("operation", operation))

	// Check if user is admin
	isAdmin, err := s.IsAdmin(ctx)
	if err != nil {
		return false, err
	}

	if isAdmin {
		// Admins can do anything
		return true, nil
	}

	// Non-admin users can only perform query operations
	if strings.HasPrefix(operation, "parent:read") ||
		strings.HasPrefix(operation, "parent:list") ||
		strings.HasPrefix(operation, "child:read") ||
		strings.HasPrefix(operation, "child:list") {
		return true, nil
	}

	// If we get here, the user is not authorized
	return false, nil
}

// IsAdmin checks if the user has admin role
func (s *AuthorizationService) IsAdmin(ctx context.Context) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "AuthorizationService.IsAdmin")
	defer span.End()

	roles, err := s.GetUserRoles(ctx)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role == "admin" {
			return true, nil
		}
	}

	return false, nil
}

// GetUserID retrieves the user ID from the context
func (s *AuthorizationService) GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetUserRoles retrieves the user roles from the context
func (s *AuthorizationService) GetUserRoles(ctx context.Context) ([]string, error) {
	roles, ok := ctx.Value(userRolesKey).([]string)
	if !ok {
		return []string{}, nil
	}
	return roles, nil
}

// WithUserID returns a new context with the user ID
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithUserRoles returns a new context with the user roles
func WithUserRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, userRolesKey, roles)
}
