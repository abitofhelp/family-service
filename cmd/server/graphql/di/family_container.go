// Copyright (c) 2025 A Bit of Help, Inc.
package di

import (
	"context"
	"fmt"

	"github.com/abitofhelp/family-service/core/application/ports"
	application "github.com/abitofhelp/family-service/core/application/services"
	domainports "github.com/abitofhelp/family-service/core/domain/ports"
	domainservices "github.com/abitofhelp/family-service/core/domain/services"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// FamilyRepositoryInitializer is a function type that initializes a family repository
type FamilyRepositoryInitializer[T domainports.FamilyRepository] func(ctx context.Context, connectionString string, logger *zap.Logger) (T, error)

// FamilyContainer is a dependency injection container for the family domain
type FamilyContainer[T domainports.FamilyRepository] struct {
	*Container
	familyRepo          T
	familyDomainService *domainservices.FamilyDomainService
	familyAppService    ports.FamilyApplicationService
}

// NewFamilyContainer creates a new family dependency injection container
func NewFamilyContainer[T domainports.FamilyRepository](
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
	initRepo FamilyRepositoryInitializer[T],
	connectionString string,
) (*FamilyContainer[T], error) {
	// Create base container
	baseContainer, err := NewContainer(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base container: %w", err)
	}

	// Create family container
	container := &FamilyContainer[T]{
		Container: baseContainer,
	}

	// Initialize repository
	container.familyRepo, err = initRepo(ctx, connectionString, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Initialize domain service
	contextLogger := logging.NewContextLogger(logger)
	container.familyDomainService = domainservices.NewFamilyDomainService(container.familyRepo, contextLogger)

	// Initialize application service
	container.familyAppService = application.NewFamilyApplicationService(
		container.familyDomainService,
		container.familyRepo,
		contextLogger,
	)

	return container, nil
}

// GetFamilyRepository returns the family repository
func (c *FamilyContainer[T]) GetFamilyRepository() domainports.FamilyRepository {
	return c.familyRepo
}

// GetFamilyDomainService returns the family domain service
func (c *FamilyContainer[T]) GetFamilyDomainService() *domainservices.FamilyDomainService {
	return c.familyDomainService
}

// GetFamilyApplicationService returns the family application service
func (c *FamilyContainer[T]) GetFamilyApplicationService() ports.FamilyApplicationService {
	return c.familyAppService
}

// GetRepositoryFactory returns the family repository
func (c *FamilyContainer[T]) GetRepositoryFactory() interface{} {
	return c.familyRepo
}

// For backward compatibility with the test file
func (c *FamilyContainer[T]) GetFamilyService() interface{} {
	return c.familyAppService
}

// For backward compatibility with the test file
func (c *FamilyContainer[T]) GetAuthorizationService() interface{} {
	return c.Container.GetAuthService()
}

// Close closes all resources
func (c *FamilyContainer[T]) Close() error {
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
