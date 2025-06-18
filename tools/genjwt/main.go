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
	config.JWT.SecretKey = "abc123"

	// Create an auth instance
	authInstance, err := auth.New(ctx, config, logger)
	if err != nil {
		logger.Fatal("Failed to create auth instance", zap.Error(err))
	}

	// Generate admin token
	adminToken, err := authInstance.GenerateToken(ctx, "admin", []string{"ADMIN"})
	if err != nil {
		logger.Fatal("Failed to generate admin token", zap.Error(err))
		return
	}

	fmt.Printf("\nAdmin Token: %s\n", adminToken)

	// Generate authuser token
	authuserToken, err := authInstance.GenerateToken(ctx, "user", []string{"AUTHUSER"})
	if err != nil {
		logger.Fatal("Failed to generate authuser token", zap.Error(err))
		return
	}

	fmt.Printf("\nAuthuser Token: %s\n", authuserToken)

	// Validate admin token
	adminClaims, err := authInstance.ValidateToken(ctx, adminToken)
	if err != nil {
		logger.Fatal("Failed to validate admin token", zap.Error(err))
		return
	}
	fmt.Printf("\nValid Admin Token, Claims: %+v\n", adminClaims)

	// Validate authuser token
	authuserClaims, err := authInstance.ValidateToken(ctx, authuserToken)
	if err != nil {
		logger.Fatal("Failed to validate authuser token", zap.Error(err))
		return
	}
	fmt.Printf("\nValid Authuser Token, Claims: %+v\n", authuserClaims)
}
