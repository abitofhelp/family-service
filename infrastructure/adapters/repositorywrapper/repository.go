// Copyright (c) 2025 A Bit of Help, Inc.

// Package repositorywrapper provides a wrapper around servicelib/repository to ensure
// the domain layer doesn't directly depend on external libraries.
//
// This package implements the Adapter pattern from Hexagonal Architecture (Ports and Adapters).
// It adapts the external servicelib/repository interface to match the interface expected
// by the domain layer, allowing the domain to remain isolated from external dependencies.
//
// Key benefits of this approach:
// 1. Dependency Inversion: The domain layer depends on abstractions (interfaces),
//    not concrete implementations.
// 2. Testability: The domain can be tested with mock implementations of these interfaces.
// 3. Flexibility: The external library can be replaced without changing the domain.
// 4. Isolation: Changes to the external library don't directly impact the domain.
//
// This pattern is essential for maintaining a clean architecture where the
// domain layer remains at the center, free from external dependencies.
package repositorywrapper

import (
	"context"
	"github.com/abitofhelp/servicelib/repository"
)

// Repository is a generic interface for repository operations.
//
// This interface defines the contract that all repositories must fulfill,
// regardless of their implementation details. It uses Go generics to
// provide type-safe operations for different entity types.
//
// The interface is intentionally minimal, focusing on the core operations
// needed by most repositories. This follows the Interface Segregation Principle
// from SOLID, which suggests that clients should not be forced to depend on
// methods they don't use.
//
// Type Parameters:
//   - T: The entity type that this repository manages
//
// This interface serves as a port in the Hexagonal Architecture, allowing
// the domain layer to interact with repositories without knowing their
// implementation details.
type Repository[T any] interface {
	// GetByID retrieves an entity by its ID.
	//
	// Parameters:
	//   - ctx: Context for the operation (used for cancellation, tracing, etc.)
	//   - id: Unique identifier of the entity to retrieve
	//
	// Returns:
	//   - The requested entity if found
	//   - An error if the operation fails (e.g., not found, database error)
	GetByID(ctx context.Context, id string) (T, error)

	// Save persists an entity.
	//
	// This method either creates a new entity or updates an existing one,
	// depending on whether the entity already exists in the repository.
	//
	// Parameters:
	//   - ctx: Context for the operation (used for cancellation, tracing, etc.)
	//   - entity: The entity to persist
	//
	// Returns:
	//   - An error if the operation fails (e.g., validation error, database error)
	Save(ctx context.Context, entity T) error

	// GetAll retrieves all entities.
	//
	// This method returns all entities of type T in the repository.
	// For large collections, consider using pagination instead.
	//
	// Parameters:
	//   - ctx: Context for the operation (used for cancellation, tracing, etc.)
	//
	// Returns:
	//   - A slice containing all entities
	//   - An error if the operation fails (e.g., database error)
	GetAll(ctx context.Context) ([]T, error)
}

// RepositoryWrapper is a wrapper around servicelib/repository.Repository.
//
// This struct implements the Adapter pattern by wrapping the external
// servicelib/repository.Repository interface and adapting it to match
// the Repository interface expected by the domain layer.
//
// The wrapper uses Go generics to provide type-safe operations for
// different entity types, just like the interface it implements.
//
// Type Parameters:
//   - T: The entity type that this repository manages
type RepositoryWrapper[T any] struct {
	repo repository.Repository[T] // The wrapped external repository
}

// NewRepositoryWrapper creates a new RepositoryWrapper.
//
// This function creates a new adapter that wraps an external repository
// implementation and makes it conform to the Repository interface
// expected by the domain layer.
//
// Parameters:
//   - repo: The external repository implementation to wrap
//
// Returns:
//   - A new RepositoryWrapper that implements the Repository interface
//
// This function follows the Adapter pattern, allowing the domain layer
// to use external repositories without directly depending on them.
func NewRepositoryWrapper[T any](repo repository.Repository[T]) *RepositoryWrapper[T] {
	return &RepositoryWrapper[T]{
		repo: repo,
	}
}

// GetByID retrieves an entity by its ID
func (r *RepositoryWrapper[T]) GetByID(ctx context.Context, id string) (T, error) {
	return r.repo.GetByID(ctx, id)
}

// Save persists an entity
func (r *RepositoryWrapper[T]) Save(ctx context.Context, entity T) error {
	return r.repo.Save(ctx, entity)
}

// GetAll retrieves all entities
func (r *RepositoryWrapper[T]) GetAll(ctx context.Context) ([]T, error) {
	return r.repo.GetAll(ctx)
}
