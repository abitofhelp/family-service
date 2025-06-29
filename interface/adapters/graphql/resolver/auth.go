// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
)

// IsAuthorized is a directive middleware for role-based access control.
func (r *Resolver) IsAuthorized(ctx context.Context, obj any, next graphql.Resolver, allowedRoles []model.Role, requiredScopes []model.Scope, resource *model.Resource) (res any, err error) {
	// TODO: Implement proper authorization check
	// This is a placeholder implementation that allows all requests
	// In a real implementation, this would check the user's roles and permissions
	return next(ctx)
}

// checkAuthorization is a helper function for role-based access control.
// It checks if the user has the required roles and permissions for a specific resource.
func checkAuthorization(ctx context.Context, allowedRoles []string, requiredScopes []string, resource string) error {
	// TODO: Implement proper authorization check
	// This is a placeholder implementation that allows all requests
	// In a real implementation, this would check the user's roles and permissions
	return nil
}
