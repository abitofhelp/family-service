// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"context"
	"fmt"
	"github.com/abitofhelp/servicelib/auth"
	"github.com/abitofhelp/servicelib/auth/middleware"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create a context
	ctx := context.Background()

	// Create a configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "abc123"

	// Create an auth instance
	authInstance, err := auth.New(ctx, config, logger)
	if err != nil {
		logger.Fatal("Failed to create auth instance", zap.Error(err))
	}

	// Define scopes and resources for different roles
	adminScopes := []string{"READ", "WRITE", "DELETE", "CREATE"}
	adminResources := []string{"FAMILY", "PARENT", "CHILD"}

	editorScopes := []string{"READ", "WRITE", "DELETE", "CREATE"}
	editorResources := []string{"FAMILY", "PARENT", "CHILD"}

	viewerScopes := []string{"READ"}
	viewerResources := []string{"FAMILY", "PARENT", "CHILD"}

	// Generate tokens with different roles
	adminToken, err := authInstance.GenerateToken(ctx, "admin_user", []string{"admin"}, adminScopes, adminResources)
	if err != nil {
		logger.Fatal("Failed to generate admin token", zap.Error(err))
		return
	}

	editorToken, err := authInstance.GenerateToken(ctx, "editor_user", []string{"editor"}, editorScopes, editorResources)
	if err != nil {
		logger.Fatal("Failed to generate editor token", zap.Error(err))
		return
	}

	viewerToken, err := authInstance.GenerateToken(ctx, "viewer_user", []string{"viewer"}, viewerScopes, viewerResources)
	if err != nil {
		logger.Fatal("Failed to generate viewer token", zap.Error(err))
		return
	}

	fmt.Printf("\nAdmin Token: %s\n", adminToken)
	fmt.Printf("\nEditor Token: %s\n", editorToken)
	fmt.Printf("\nViewer Token: %s\n", viewerToken)

	// Validate tokens and test authorization
	testAuthorization(ctx, authInstance, adminToken, "admin_user", []string{"admin"})
	testAuthorization(ctx, authInstance, editorToken, "editor_user", []string{"editor"})
	testAuthorization(ctx, authInstance, viewerToken, "viewer_user", []string{"viewer"})
}

func testAuthorization(ctx context.Context, authInstance *auth.Auth, token string, expectedUserID string, expectedRoles []string) {
	// Validate token
	claims, err := authInstance.ValidateToken(ctx, token)
	if err != nil {
		fmt.Printf("\nFailed to validate token for %s: %v\n", expectedUserID, err)
		return
	}

	fmt.Printf("\nValid Token for %s, Claims: %+v\n", expectedUserID, claims)

	// Create a context with the claims
	ctxWithClaims := middleware.WithUserID(ctx, claims.UserID)
	ctxWithClaims = middleware.WithUserRoles(ctxWithClaims, claims.Roles)

	// Test IsAuthorized for different operations
	fmt.Printf("\n--- Testing Mutations (only 'admin' and 'editor' roles allowed) ---\n")
	testOperation(ctxWithClaims, "CreateFamily", []string{"admin", "editor"})
	testOperation(ctxWithClaims, "AddParent", []string{"admin", "editor"})
	testOperation(ctxWithClaims, "AddChild", []string{"admin", "editor"})
	testOperation(ctxWithClaims, "RemoveChild", []string{"admin", "editor"})
	testOperation(ctxWithClaims, "MarkParentDeceased", []string{"admin", "editor"})
	testOperation(ctxWithClaims, "Divorce", []string{"admin", "editor"})

	fmt.Printf("\n--- Testing Queries (all roles allowed) ---\n")
	testOperation(ctxWithClaims, "GetFamily", []string{"admin", "editor", "viewer"})
	testOperation(ctxWithClaims, "GetAllFamilies", []string{"admin", "editor", "viewer"})
	testOperation(ctxWithClaims, "FindFamiliesByParent", []string{"admin", "editor", "viewer"})
	testOperation(ctxWithClaims, "FindFamilyByChild", []string{"admin", "editor", "viewer"})
	testOperation(ctxWithClaims, "Parents", []string{"admin", "editor", "viewer"})
	testOperation(ctxWithClaims, "CountFamilies", []string{"admin", "editor", "viewer"})
	testOperation(ctxWithClaims, "CountParents", []string{"admin", "editor", "viewer"})
	testOperation(ctxWithClaims, "CountChildren", []string{"admin", "editor", "viewer"})
}

func testOperation(ctx context.Context, operation string, allowedRoles []string) {
	isAuthorized := middleware.IsAuthorized(ctx, allowedRoles)
	fmt.Printf("Operation: %s, Authorized: %v\n", operation, isAuthorized)
}
