// Copyright (c) 2025 A Bit of Help, Inc.

// Package repositorywrapper provides a wrapper around servicelib/repository to ensure
// the domain layer doesn't directly depend on external libraries.
package repositorywrapper

import (
	"context"
	"github.com/abitofhelp/servicelib/repository"
)

// Repository is a generic interface for repository operations
// It wraps the servicelib/repository.Repository interface
type Repository[T any] interface {
	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id string) (T, error)

	// Save persists an entity
	Save(ctx context.Context, entity T) error

	// GetAll retrieves all entities
	GetAll(ctx context.Context) ([]T, error)
}

// RepositoryWrapper is a wrapper around servicelib/repository.Repository
type RepositoryWrapper[T any] struct {
	repo repository.Repository[T]
}

// NewRepositoryWrapper creates a new RepositoryWrapper
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