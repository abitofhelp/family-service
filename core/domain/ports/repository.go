// Copyright (c) 2025 A Bit of Help, Inc.

package ports

import (
	"context"
	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/infrastructure/adapters/repositorywrapper"
)

// FamilyRepository defines the interface for family persistence operations
// This interface represents a port in the Hexagonal Architecture pattern
// It's defined in the domain layer but implemented in the infrastructure layer
type FamilyRepository interface {
	// Embed the generic Repository interface with Family entity
	repositorywrapper.Repository[*entity.Family]

	// FindByParentID finds families that contain a specific parent
	FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error)

	// FindByChildID finds the family that contains a specific child
	FindByChildID(ctx context.Context, childID string) (*entity.Family, error)
}
