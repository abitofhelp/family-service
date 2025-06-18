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

	// Generate tokens with different roles
	adminToken, err := authInstance.GenerateToken(ctx, "admin_user", []string{"admin"})
	if err != nil {
		logger.Fatal("Failed to generate admin token", zap.Error(err))
		return
	}

	authuserToken, err := authInstance.GenerateToken(ctx, "authuser_user", []string{"authuser"})
	if err != nil {
		logger.Fatal("Failed to generate authuser token", zap.Error(err))
		return
	}

	editorToken, err := authInstance.GenerateToken(ctx, "editor_user", []string{"family_editor"})
	if err != nil {
		logger.Fatal("Failed to generate editor token", zap.Error(err))
		return
	}

	viewerToken, err := authInstance.GenerateToken(ctx, "viewer_user", []string{"viewer"})
	if err != nil {
		logger.Fatal("Failed to generate viewer token", zap.Error(err))
		return
	}

	fmt.Printf("\nAdmin Token: %s\n", adminToken)
	fmt.Printf("\nAuthUser Token: %s\n", authuserToken)
	fmt.Printf("\nEditor Token: %s\n", editorToken)
	fmt.Printf("\nViewer Token: %s\n", viewerToken)

	// Validate tokens and test authorization
	testAuthorization(ctx, authInstance, adminToken, "admin_user", []string{"admin"})
	testAuthorization(ctx, authInstance, authuserToken, "authuser_user", []string{"authuser"})
	testAuthorization(ctx, authInstance, editorToken, "editor_user", []string{"family_editor"})
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
	fmt.Printf("\n--- Testing Mutations (only 'admin' role allowed) ---\n")
	testOperation(ctxWithClaims, "CreateFamily", []string{"admin"})
	testOperation(ctxWithClaims, "AddParent", []string{"admin"})
	testOperation(ctxWithClaims, "AddChild", []string{"admin"})
	testOperation(ctxWithClaims, "RemoveChild", []string{"admin"})
	testOperation(ctxWithClaims, "MarkParentDeceased", []string{"admin"})
	testOperation(ctxWithClaims, "Divorce", []string{"admin"})

	fmt.Printf("\n--- Testing Queries (only 'admin' and 'authuser' roles allowed) ---\n")
	testOperation(ctxWithClaims, "GetFamily", []string{"admin", "authuser"})
	testOperation(ctxWithClaims, "GetAllFamilies", []string{"admin", "authuser"})
	testOperation(ctxWithClaims, "FindFamiliesByParent", []string{"admin", "authuser"})
	testOperation(ctxWithClaims, "FindFamilyByChild", []string{"admin", "authuser"})
	testOperation(ctxWithClaims, "Parents", []string{"admin", "authuser"})
	testOperation(ctxWithClaims, "CountFamilies", []string{"admin", "authuser"})
	testOperation(ctxWithClaims, "CountParents", []string{"admin", "authuser"})
	testOperation(ctxWithClaims, "CountChildren", []string{"admin", "authuser"})
}

func testOperation(ctx context.Context, operation string, allowedRoles []string) {
	isAuthorized := middleware.IsAuthorized(ctx, allowedRoles)
	fmt.Printf("Operation: %s, Authorized: %v\n", operation, isAuthorized)
}
