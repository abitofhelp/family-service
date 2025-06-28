// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"context"
	"fmt"
	"github.com/abitofhelp/servicelib/auth"
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
	config.JWT.SecretKey = "01234567890123456789012345678901"

	// Create an auth instance
	authInstance, err := auth.New(ctx, config, logger)
	if err != nil {
		logger.Fatal("Failed to create auth instance", zap.Error(err))
	}

	// Generate admin token with all scopes for all resources
	adminScopes := []string{"READ", "WRITE", "DELETE", "CREATE"}
	adminResources := []string{"FAMILY", "PARENT", "CHILD"}
	adminToken, err := authInstance.GenerateToken(ctx, "admin", []string{"ADMIN"}, adminScopes, adminResources)
	if err != nil {
		logger.Fatal("Failed to generate admin token", zap.Error(err))
		return
	}

	fmt.Printf("\nAdmin Token: %s\n", adminToken)

	// Generate editor token with all scopes for all resources
	editorScopes := []string{"READ", "WRITE", "DELETE", "CREATE"}
	editorResources := []string{"FAMILY", "PARENT", "CHILD"}
	editorToken, err := authInstance.GenerateToken(ctx, "editor", []string{"EDITOR"}, editorScopes, editorResources)
	if err != nil {
		logger.Fatal("Failed to generate editor token", zap.Error(err))
		return
	}

	fmt.Printf("\nEditor Token: %s\n", editorToken)

	// Generate viewer token with only READ scope for all resources
	viewerScopes := []string{"READ"}
	viewerResources := []string{"FAMILY", "PARENT", "CHILD"}
	viewerToken, err := authInstance.GenerateToken(ctx, "viewer", []string{"VIEWER"}, viewerScopes, viewerResources)
	if err != nil {
		logger.Fatal("Failed to generate viewer token", zap.Error(err))
		return
	}

	fmt.Printf("\nViewer Token: %s\n", viewerToken)

	// Validate admin token
	adminClaims, err := authInstance.ValidateToken(ctx, adminToken)
	if err != nil {
		logger.Fatal("Failed to validate admin token", zap.Error(err))
		return
	}
	fmt.Printf("\nValid Admin Token, Claims: %+v\n", adminClaims)

	// Validate editor token
	editorClaims, err := authInstance.ValidateToken(ctx, editorToken)
	if err != nil {
		logger.Fatal("Failed to validate editor token", zap.Error(err))
		return
	}
	fmt.Printf("\nValid Editor Token, Claims: %+v\n", editorClaims)

	// Validate viewer token
	viewerClaims, err := authInstance.ValidateToken(ctx, viewerToken)
	if err != nil {
		logger.Fatal("Failed to validate viewer token", zap.Error(err))
		return
	}
	fmt.Printf("\nValid Viewer Token, Claims: %+v\n", viewerClaims)
}
