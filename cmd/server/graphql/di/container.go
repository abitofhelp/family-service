// Package di provides a dependency injection container for the GraphQL server.
package di

import (
	"context"
	"fmt"

	appports "github.com/abitofhelp/family-service/core/application/ports"
	application "github.com/abitofhelp/family-service/core/application/services"
	domainports "github.com/abitofhelp/family-service/core/domain/ports"
	domainservices "github.com/abitofhelp/family-service/core/domain/services"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	basedi "github.com/abitofhelp/servicelib/di"
	"go.uber.org/zap"
)

// Container is a dependency injection container for the GraphQL server
type Container struct {
	*basedi.Container
	familyRepo          domainports.FamilyRepository
	familyDomainService *domainservices.FamilyDomainService
	familyAppService    appports.FamilyApplicationService
	dbType              string
}

// NewContainer creates a new dependency injection container for the GraphQL server
func NewContainer(ctx context.Context, logger *zap.Logger, cfg *config.Config) (*Container, error) {
	// Create base container
	baseContainer, err := basedi.NewContainer(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base container: %w", err)
	}

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
		repo, err := basedi.InitMongoRepository(ctx, cfg.Database.MongoDB.URI, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize MongoDB repository: %w", err)
		}
		container.familyRepo = repo
	case "postgres":
		// Initialize PostgreSQL repository
		repo, err := basedi.InitPostgresRepository(ctx, cfg.Database.Postgres.DSN, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PostgreSQL repository: %w", err)
		}
		container.familyRepo = repo
	case "sqlite":
		// Initialize SQLite repository
		repo, err := basedi.InitSQLiteRepository(ctx, cfg.Database.SQLite.URI, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SQLite repository: %w", err)
		}
		container.familyRepo = repo
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Initialize domain service
	container.familyDomainService = domainservices.NewFamilyDomainService(container.familyRepo)

	// Initialize application service
	container.familyAppService = application.NewFamilyApplicationService(
		container.familyDomainService,
		container.familyRepo,
	)

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
	return nil
}

// Close closes all resources
func (c *Container) Close() error {
	var errs []error

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
