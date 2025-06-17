package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewAuthorizationService(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()

	// Execute
	service := NewAuthorizationService(logger)

	// Verify
	assert.NotNil(t, service)
	assert.Equal(t, logger, service.logger)
	assert.NotNil(t, service.tracer)
}

func TestIsAuthorized_Admin(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context with admin role
	ctx := context.Background()
	ctx = WithUserID(ctx, "admin-user")
	ctx = WithUserRoles(ctx, []string{"admin", "user"})
	
	// Execute
	authorized, err := service.IsAuthorized(ctx, "parent:create")
	
	// Verify
	assert.NoError(t, err)
	assert.True(t, authorized)
}

func TestIsAuthorized_NonAdmin_ReadOperation(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context with non-admin role
	ctx := context.Background()
	ctx = WithUserID(ctx, "regular-user")
	ctx = WithUserRoles(ctx, []string{"user"})
	
	// Test cases for read operations
	readOperations := []string{
		"parent:read",
		"parent:list",
		"child:read",
		"child:list",
	}
	
	for _, operation := range readOperations {
		t.Run(operation, func(t *testing.T) {
			// Execute
			authorized, err := service.IsAuthorized(ctx, operation)
			
			// Verify
			assert.NoError(t, err)
			assert.True(t, authorized)
		})
	}
}

func TestIsAuthorized_NonAdmin_WriteOperation(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context with non-admin role
	ctx := context.Background()
	ctx = WithUserID(ctx, "regular-user")
	ctx = WithUserRoles(ctx, []string{"user"})
	
	// Test cases for write operations
	writeOperations := []string{
		"parent:create",
		"parent:update",
		"parent:delete",
		"child:create",
		"child:update",
		"child:delete",
	}
	
	for _, operation := range writeOperations {
		t.Run(operation, func(t *testing.T) {
			// Execute
			authorized, err := service.IsAuthorized(ctx, operation)
			
			// Verify
			assert.NoError(t, err)
			assert.False(t, authorized)
		})
	}
}

func TestIsAdmin_WithAdminRole(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context with admin role
	ctx := context.Background()
	ctx = WithUserRoles(ctx, []string{"admin", "user"})
	
	// Execute
	isAdmin, err := service.IsAdmin(ctx)
	
	// Verify
	assert.NoError(t, err)
	assert.True(t, isAdmin)
}

func TestIsAdmin_WithoutAdminRole(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context with non-admin roles
	ctx := context.Background()
	ctx = WithUserRoles(ctx, []string{"user", "editor"})
	
	// Execute
	isAdmin, err := service.IsAdmin(ctx)
	
	// Verify
	assert.NoError(t, err)
	assert.False(t, isAdmin)
}

func TestIsAdmin_NoRoles(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context with no roles
	ctx := context.Background()
	
	// Execute
	isAdmin, err := service.IsAdmin(ctx)
	
	// Verify
	assert.NoError(t, err)
	assert.False(t, isAdmin)
}

func TestGetUserID_Present(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context with user ID
	ctx := context.Background()
	userID := "test-user-id"
	ctx = WithUserID(ctx, userID)
	
	// Execute
	retrievedID, err := service.GetUserID(ctx)
	
	// Verify
	assert.NoError(t, err)
	assert.Equal(t, userID, retrievedID)
}

func TestGetUserID_Missing(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context without user ID
	ctx := context.Background()
	
	// Execute
	retrievedID, err := service.GetUserID(ctx)
	
	// Verify
	assert.Error(t, err)
	assert.Empty(t, retrievedID)
	assert.Contains(t, err.Error(), "user ID not found in context")
}

func TestGetUserRoles_Present(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context with user roles
	ctx := context.Background()
	roles := []string{"admin", "user"}
	ctx = WithUserRoles(ctx, roles)
	
	// Execute
	retrievedRoles, err := service.GetUserRoles(ctx)
	
	// Verify
	assert.NoError(t, err)
	assert.Equal(t, roles, retrievedRoles)
}

func TestGetUserRoles_Missing(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	service := NewAuthorizationService(logger)
	
	// Create a context without user roles
	ctx := context.Background()
	
	// Execute
	retrievedRoles, err := service.GetUserRoles(ctx)
	
	// Verify
	assert.NoError(t, err)
	assert.Empty(t, retrievedRoles)
}

func TestWithUserID(t *testing.T) {
	// Setup
	ctx := context.Background()
	userID := "test-user-id"
	
	// Execute
	newCtx := WithUserID(ctx, userID)
	
	// Verify
	retrievedID, ok := newCtx.Value(userIDKey).(string)
	assert.True(t, ok)
	assert.Equal(t, userID, retrievedID)
}

func TestWithUserRoles(t *testing.T) {
	// Setup
	ctx := context.Background()
	roles := []string{"admin", "user"}
	
	// Execute
	newCtx := WithUserRoles(ctx, roles)
	
	// Verify
	retrievedRoles, ok := newCtx.Value(userRolesKey).([]string)
	assert.True(t, ok)
	assert.Equal(t, roles, retrievedRoles)
}