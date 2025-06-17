package ports

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
)

// ApplicationService is a generic interface for application services
// This interface represents a port in the Hexagonal Architecture pattern
// It's defined in the application layer but implemented in the application layer
// and used by the interface layer
type ApplicationService[T any, D any] interface {
	// Create creates a new entity
	Create(ctx context.Context, dto D) (D, error)

	// GetByID retrieves an entity by ID
	GetByID(ctx context.Context, id string) (D, error)

	// GetAll retrieves all entities
	GetAll(ctx context.Context) ([]D, error)
}

// FamilyApplicationService defines the interface for family application services
// This interface represents a port in the Hexagonal Architecture pattern
// It's defined in the application layer but implemented in the application layer
// and used by the interface layer
type FamilyApplicationService interface {
	// Embed the generic ApplicationService interface with Family entity and DTO
	ApplicationService[*entity.Family, *entity.FamilyDTO]

	// CreateFamily creates a new family (alias for Create)
	CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error)

	// GetFamily retrieves a family by ID (alias for GetByID)
	GetFamily(ctx context.Context, id string) (*entity.FamilyDTO, error)

	// GetAllFamilies retrieves all families (alias for GetAll)
	GetAllFamilies(ctx context.Context) ([]*entity.FamilyDTO, error)

	// AddParent adds a parent to a family
	AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error)

	// AddChild adds a child to a family
	AddChild(ctx context.Context, familyID string, childDTO entity.ChildDTO) (*entity.FamilyDTO, error)

	// RemoveChild removes a child from a family
	RemoveChild(ctx context.Context, familyID string, childID string) (*entity.FamilyDTO, error)

	// MarkParentDeceased marks a parent as deceased
	MarkParentDeceased(ctx context.Context, familyID string, parentID string, deathDate time.Time) (*entity.FamilyDTO, error)

	// Divorce handles the divorce process
	Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error)

	// FindFamiliesByParent finds families that contain a specific parent
	FindFamiliesByParent(ctx context.Context, parentID string) ([]*entity.FamilyDTO, error)

	// FindFamilyByChild finds the family that contains a specific child
	FindFamilyByChild(ctx context.Context, childID string) (*entity.FamilyDTO, error)
}
