package ports

import (
	"context"
	"github.com/abitofhelp/family-service/core/domain/entity"
)

// Repository is a generic repository interface for entity persistence operations
// This interface represents a port in the Hexagonal Architecture pattern
// It's defined in the domain layer but implemented in the infrastructure layer
type Repository[T any] interface {
	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id string) (T, error)

	// GetAll retrieves all entities
	GetAll(ctx context.Context) ([]T, error)

	// Save persists an entity
	Save(ctx context.Context, entity T) error
}

// FamilyRepository defines the interface for family persistence operations
// This interface represents a port in the Hexagonal Architecture pattern
// It's defined in the domain layer but implemented in the infrastructure layer
type FamilyRepository interface {
	// Embed the generic Repository interface with Family entity
	Repository[*entity.Family]

	// FindByParentID finds families that contain a specific parent
	FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error)

	// FindByChildID finds the family that contains a specific child
	FindByChildID(ctx context.Context, childID string) (*entity.Family, error)
}
