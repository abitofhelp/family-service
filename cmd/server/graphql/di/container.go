// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides a dependency injection container for the GraphQL server.
package di

import (
	"context"
	"fmt"

	appports "github.com/abitofhelp/family-service/core/application/ports"
	application "github.com/abitofhelp/family-service/core/application/services"
	domainports "github.com/abitofhelp/family-service/core/domain/ports"
	domainservices "github.com/abitofhelp/family-service/core/domain/services"
	"github.com/abitofhelp/family-service/infrastructure/adapters/cache"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	adaptdi "github.com/abitofhelp/family-service/infrastructure/adapters/di"
	"github.com/abitofhelp/family-service/infrastructure/adapters/mongo"
	"github.com/abitofhelp/family-service/infrastructure/adapters/postgres"
	"github.com/abitofhelp/family-service/infrastructure/adapters/sqlite"
	"github.com/abitofhelp/servicelib/auth"
	basedi "github.com/abitofhelp/servicelib/di"
	"go.uber.org/zap"
)

// Container is a dependency injection container for the GraphQL server
type Container struct {
	*basedi.Container
	familyRepo          domainports.FamilyRepository
	familyDomainService *domainservices.FamilyDomainService
	familyAppService    appports.FamilyApplicationService
	authService         *auth.Auth
	dbType              string
	cache               *cache.Cache
}

// NewContainer creates a new dependency injection container for the GraphQL server
func NewContainer(ctx context.Context, logger *zap.Logger, cfg *config.Config) (*Container, error) {
	// Create base container
	baseContainer, err := basedi.NewContainer(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base container: %w", err)
	}

	// Set global configuration for repositories
	mongo.SetGlobalConfig(cfg)
	postgres.SetGlobalConfig(cfg)
	sqlite.SetGlobalConfig(cfg)

	// Create GraphQL-specific container
	container := &Container{
		Container: baseContainer,
		dbType:    cfg.Database.Type,
	}

	// Initialize repository based on database type
	dbType := cfg.Database.Type
	switch dbType {
	case "mongodb":
		// Initialize MongoDB repository
		repo, err := adaptdi.InitMongoRepository(ctx, cfg.Database.MongoDB.URI, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize MongoDB repository: %w", err)
		}
		container.familyRepo = repo
	case "postgres":
		// Initialize PostgreSQL repository
		repo, err := adaptdi.InitPostgresRepository(ctx, cfg.Database.Postgres.DSN, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PostgreSQL repository: %w", err)
		}
		container.familyRepo = repo
	case "sqlite":
		// Initialize SQLite repository
		repo, err := adaptdi.InitSQLiteRepository(ctx, cfg.Database.SQLite.URI, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SQLite repository: %w", err)
		}
		container.familyRepo = repo
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Initialize cache
	cacheInstance, err := cache.NewCache(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}
	container.cache = cacheInstance

	// Initialize domain service
	container.familyDomainService = domainservices.NewFamilyDomainService(container.familyRepo, container.GetContextLogger())

	// Initialize application service
	container.familyAppService = application.NewFamilyApplicationService(
		container.familyDomainService,
		container.familyRepo,
		container.GetContextLogger(),
		container.cache,
	)

	// Initialize auth service
	// Note: In the future, this should be configured to use a remote authorization server
	// instead of local token validation for improved security and centralized management.
	authConfig := auth.DefaultConfig()
	authConfig.JWT.SecretKey = cfg.Auth.JWT.SecretKey
	authConfig.JWT.Issuer = cfg.Auth.JWT.Issuer
	authConfig.JWT.TokenDuration = cfg.Auth.JWT.TokenDuration
	authConfig.Middleware.SkipPaths = []string{"/health", "/metrics", "/playground", "/graphql/health"}

	authService, err := auth.New(ctx, authConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %w", err)
	}
	container.authService = authService

	return container, nil
}

// GetFamilyRepository returns the family repository
func (c *Container) GetFamilyRepository() domainports.FamilyRepository {
	return c.familyRepo
}

// GetFamilyDomainService returns the family domain service
func (c *Container) GetFamilyDomainService() *domainservices.FamilyDomainService {
	return c.familyDomainService
}

// GetFamilyApplicationService returns the family application service
func (c *Container) GetFamilyApplicationService() appports.FamilyApplicationService {
	return c.familyAppService
}

// GetRepositoryFactory returns the family repository
func (c *Container) GetRepositoryFactory() interface{} {
	return c.familyRepo
}

// For backward compatibility with the test file
func (c *Container) GetFamilyService() interface{} {
	return c.familyAppService
}

// For backward compatibility with the test file
func (c *Container) GetAuthorizationService() interface{} {
	return c.authService
}

// GetAuthService returns the auth service
func (c *Container) GetAuthService() *auth.Auth {
	return c.authService
}

// Close closes all resources
func (c *Container) Close() error {
	var errs []error

	// Shutdown cache if it exists
	if c.cache != nil {
		c.cache.Shutdown()
	}

	// Add resource cleanup here as needed
	// For example, close database connections if they implement a Close method

	// Close base container
	if err := c.Container.Close(); err != nil {
		errs = append(errs, err)
	}

	// Return a combined error if any occurred
	if len(errs) > 0 {
		errMsg := "failed to close one or more resources:"
		for _, err := range errs {
			errMsg += " " + err.Error()
		}
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}
