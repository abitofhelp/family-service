// Copyright (c) 2025 A Bit of Help, Inc.

// Package application contains application services that orchestrate domain operations.
//
// In Clean Architecture and Hexagonal Architecture, the application layer sits between
// the domain layer (core business logic) and the interface/infrastructure layers.
// It coordinates the flow of data and operations between these layers, ensuring that
// domain logic remains pure and focused on business rules.
//
// Application services implement use cases by:
// 1. Receiving requests from the interface layer (e.g., GraphQL resolvers)
// 2. Coordinating with domain services to execute business logic
// 3. Interacting with repositories to persist or retrieve data
// 4. Handling cross-cutting concerns like caching, logging, and error handling
// 5. Returning responses back to the interface layer
//
// This layer helps maintain separation of concerns and keeps the domain layer
// isolated from external dependencies and implementation details.
package application

import (
	"context"
	"fmt"
	"github.com/abitofhelp/servicelib/di"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	domainports "github.com/abitofhelp/family-service/core/domain/ports"
	domainservices "github.com/abitofhelp/family-service/core/domain/services"
	"github.com/abitofhelp/family-service/infrastructure/adapters/cachewrapper"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// BaseApplicationService is a generic implementation of the ApplicationService interface.
//
// This struct uses Go generics to provide a reusable base for application services
// that work with different entity types. The type parameters allow for type-safe
// operations across different entity types while sharing common functionality.
//
// Type Parameters:
//   - T: The domain entity type (e.g., *entity.Family)
//   - D: The DTO (Data Transfer Object) type (e.g., *entity.FamilyDTO)
//
// This pattern reduces code duplication and ensures consistent behavior
// across different application services.
type BaseApplicationService[T any, D any] struct {
	// Common dependencies and methods for all application services
}

// FamilyApplicationService implements the application service for family-related use cases.
//
// This service acts as an orchestrator between the interface layer (e.g., GraphQL resolvers)
// and the domain layer. It coordinates operations on Family entities by:
// 1. Delegating business logic to the domain service
// 2. Persisting changes through the repository
// 3. Handling cross-cutting concerns like caching and logging
// 4. Converting between domain entities and DTOs
//
// The service implements the ports.FamilyApplicationService interface, which
// allows the interface layer to interact with it without knowing the implementation details.
// This is a key aspect of the Hexagonal Architecture (Ports and Adapters) pattern.
type FamilyApplicationService struct {
	BaseApplicationService[*entity.Family, *entity.FamilyDTO]
	familyService *domainservices.FamilyDomainService // Domain service for family-related business logic
	familyRepo    domainports.FamilyRepository        // Repository for persisting and retrieving families
	logger        *logging.ContextLogger              // Logger for recording operations and errors
	cache         *cache.Cache                        // Optional cache for improving performance
}

// Ensure FamilyApplicationService implements di.ApplicationService
var _ di.ApplicationService = (*FamilyApplicationService)(nil)

// NewFamilyApplicationService creates a new FamilyApplicationService with the required dependencies.
//
// This function follows the Dependency Injection pattern, requiring all dependencies
// to be provided rather than creating them internally. This approach:
// 1. Makes dependencies explicit and visible
// 2. Improves testability by allowing mock implementations
// 3. Gives control over the lifecycle of dependencies to the caller
//
// Parameters:
//   - familyService: Domain service that contains family-related business logic
//   - familyRepo: Repository for persisting and retrieving family entities
//   - logger: Logger for recording operations and errors
//   - cache: Optional cache for improving performance (can be nil)
//
// Returns:
//   - A new FamilyApplicationService instance
//
// Panics if any required dependency (except cache) is nil, as the service
// cannot function without these core dependencies.
func NewFamilyApplicationService(
	familyService *domainservices.FamilyDomainService,
	familyRepo domainports.FamilyRepository,
	logger *logging.ContextLogger,
	cache *cache.Cache,
) *FamilyApplicationService {
	if familyService == nil {
		panic("family service cannot be nil")
	}
	if familyRepo == nil {
		panic("family repository cannot be nil")
	}
	if logger == nil {
		panic("logger cannot be nil")
	}
	return &FamilyApplicationService{
		familyService: familyService,
		familyRepo:    familyRepo,
		logger:        logger,
		cache:         cache,
	}
}

// Create creates a new family entity based on the provided DTO.
//
// This method implements the ApplicationService.Create interface method and serves
// as the entry point for creating new families. It:
// 1. Logs the operation for observability
// 2. Delegates the business logic to the domain service
// 3. Handles errors and logs failures
// 4. Returns the created family as a DTO
//
// The method doesn't directly interact with the repository - that's the
// responsibility of the domain service. This separation ensures that
// business rules are consistently applied regardless of how the operation
// is initiated.
//
// Parameters:
//   - ctx: Context for the operation (used for cancellation, tracing, etc.)
//   - dto: Data Transfer Object containing the family data to create
//
// Returns:
//   - The created family as a DTO if successful
//   - An error if the operation fails (e.g., validation error, database error)
//
// Example usage (from a GraphQL resolver):
//
//	family, err := familyService.Create(ctx, &entity.FamilyDTO{
//	    Status: "SINGLE",
//	    Parents: []entity.ParentDTO{
//	        {FirstName: "John", LastName: "Doe", BirthDate: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)},
//	    },
//	})
func (s *FamilyApplicationService) Create(ctx context.Context, dto *entity.FamilyDTO) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Creating new family", zap.String("family_id", dto.ID), zap.String("status", dto.Status))

	// Delegate to domain service
	family, err := s.familyService.CreateFamily(ctx, *dto)
	if err != nil {
		s.logger.Error(ctx, "Failed to create family", zap.Error(err), zap.String("family_id", dto.ID))
		return nil, err
	}

	s.logger.Info(ctx, "Successfully created family", zap.String("family_id", family.ID), zap.Int("parent_count", family.ParentCount), zap.Int("children_count", family.ChildrenCount))
	return family, nil
}

// GetByID retrieves a family by its unique identifier.
//
// This method implements the ApplicationService.GetByID interface method and
// demonstrates several important patterns:
// 1. Caching: It uses a cache-aside pattern to improve performance
// 2. Logging: It logs the operation for observability
// 3. Error handling: It properly handles and logs errors
//
// The caching implementation uses a function-based approach where the actual
// data retrieval is wrapped in a callback function. This pattern:
// - Centralizes cache logic to avoid duplication
// - Handles cache misses transparently
// - Ensures consistent error handling
//
// Parameters:
//   - ctx: Context for the operation (used for cancellation, tracing, etc.)
//   - id: Unique identifier of the family to retrieve
//
// Returns:
//   - The requested family as a DTO if found
//   - An error if the operation fails (e.g., not found, database error)
func (s *FamilyApplicationService) GetByID(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Retrieving family by ID", zap.String("family_id", id))

	// Create cache key
	cacheKey := fmt.Sprintf("family:%s", id)

	// Try to get from cache or call the domain service
	result, err := cache.WithContextCache(ctx, s.cache, cacheKey, func(ctx context.Context) (interface{}, error) {
		// Delegate to domain service
		return s.familyService.GetFamily(ctx, id)
	})
	if err != nil {
		s.logger.Error(ctx, "Failed to retrieve family", zap.Error(err), zap.String("family_id", id))
		return nil, err
	}

	// Type assertion
	family, ok := result.(*entity.FamilyDTO)
	if !ok {
		s.logger.Error(ctx, "Failed to cast cached result to FamilyDTO", zap.String("family_id", id))
		return nil, errors.NewApplicationError(errors.InternalErrorCode, "failed to cast cached result to FamilyDTO", nil)
	}

	s.logger.Info(ctx, "Successfully retrieved family", zap.String("family_id", family.ID), zap.String("status", family.Status))
	return family, nil
}

// GetAll retrieves all families (implements ApplicationService.GetAll)
func (s *FamilyApplicationService) GetAll(ctx context.Context) ([]*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Retrieving all families")

	// Use repository to get all families
	families, err := s.familyRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to retrieve all families", zap.Error(err))
		return nil, errors.NewApplicationError(errors.DatabaseErrorCode, "failed to get all families", err)
	}

	// Convert domain entities to DTOs
	dtos := make([]*entity.FamilyDTO, 0, len(families))
	for _, fam := range families {
		dto := fam.ToDTO()
		dtos = append(dtos, &dto)
	}

	s.logger.Info(ctx, "Successfully retrieved all families", zap.Int("count", len(dtos)))
	return dtos, nil
}

// AddParent adds a parent to a family
func (s *FamilyApplicationService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Adding parent to family", 
		zap.String("family_id", familyID), 
		zap.String("parent_id", parentDTO.ID),
		zap.String("parent_first_name", parentDTO.FirstName),
		zap.String("parent_last_name", parentDTO.LastName))

	// Delegate to domain service
	family, err := s.familyService.AddParent(ctx, familyID, parentDTO)
	if err != nil {
		s.logger.Error(ctx, "Failed to add parent to family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("parent_id", parentDTO.ID))
		return nil, err
	}

	s.logger.Info(ctx, "Successfully added parent to family", 
		zap.String("family_id", family.ID), 
		zap.Int("parent_count", family.ParentCount))
	return family, nil
}

// AddChild adds a child to a family
func (s *FamilyApplicationService) AddChild(ctx context.Context, familyID string, childDTO entity.ChildDTO) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Adding child to family", 
		zap.String("family_id", familyID), 
		zap.String("child_id", childDTO.ID),
		zap.String("child_first_name", childDTO.FirstName),
		zap.String("child_last_name", childDTO.LastName))

	// Delegate to domain service
	family, err := s.familyService.AddChild(ctx, familyID, childDTO)
	if err != nil {
		s.logger.Error(ctx, "Failed to add child to family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("child_id", childDTO.ID))
		return nil, err
	}

	s.logger.Info(ctx, "Successfully added child to family", 
		zap.String("family_id", family.ID), 
		zap.Int("children_count", family.ChildrenCount))
	return family, nil
}

// RemoveChild removes a child from a family
func (s *FamilyApplicationService) RemoveChild(ctx context.Context, familyID string, childID string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Removing child from family", 
		zap.String("family_id", familyID), 
		zap.String("child_id", childID))

	// Delegate to domain service
	family, err := s.familyService.RemoveChild(ctx, familyID, childID)
	if err != nil {
		s.logger.Error(ctx, "Failed to remove child from family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("child_id", childID))
		return nil, err
	}

	s.logger.Info(ctx, "Successfully removed child from family", 
		zap.String("family_id", family.ID), 
		zap.Int("children_count", family.ChildrenCount))
	return family, nil
}

// MarkParentDeceased marks a parent as deceased
func (s *FamilyApplicationService) MarkParentDeceased(ctx context.Context, familyID string, parentID string, deathDate time.Time) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Marking parent as deceased", 
		zap.String("family_id", familyID), 
		zap.String("parent_id", parentID),
		zap.Time("death_date", deathDate))

	// Delegate to domain service
	family, err := s.familyService.MarkParentDeceased(ctx, familyID, parentID, deathDate)
	if err != nil {
		s.logger.Error(ctx, "Failed to mark parent as deceased", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("parent_id", parentID))
		return nil, err
	}

	s.logger.Info(ctx, "Successfully marked parent as deceased", 
		zap.String("family_id", family.ID), 
		zap.String("status", family.Status))
	return family, nil
}

// Divorce handles the divorce process
func (s *FamilyApplicationService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Processing divorce", 
		zap.String("family_id", familyID), 
		zap.String("custodial_parent_id", custodialParentID))

	// Delegate to domain service
	family, err := s.familyService.Divorce(ctx, familyID, custodialParentID)
	if err != nil {
		s.logger.Error(ctx, "Failed to process divorce", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("custodial_parent_id", custodialParentID))
		return nil, err
	}

	s.logger.Info(ctx, "Successfully processed divorce", 
		zap.String("family_id", family.ID), 
		zap.String("status", family.Status))
	return family, nil
}

// FindFamiliesByParent finds families that contain a specific parent
func (s *FamilyApplicationService) FindFamiliesByParent(ctx context.Context, parentID string) ([]*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Finding families by parent ID", zap.String("parent_id", parentID))

	if parentID == "" {
		s.logger.Warn(ctx, "Parent ID is required for FindFamiliesByParent")
		return nil, errors.NewValidationError("parent ID is required", "parentID", nil)
	}

	// Use repository to find families by parent ID
	families, err := s.familyRepo.FindByParentID(ctx, parentID)
	if err != nil {
		s.logger.Error(ctx, "Failed to find families by parent ID", 
			zap.Error(err), 
			zap.String("parent_id", parentID))
		return nil, errors.NewApplicationError(errors.DatabaseErrorCode, "failed to find families by parent ID", err)
	}

	// Convert domain entities to DTOs
	dtos := make([]*entity.FamilyDTO, 0, len(families))
	for _, fam := range families {
		dto := fam.ToDTO()
		dtos = append(dtos, &dto)
	}

	s.logger.Info(ctx, "Successfully found families by parent ID", 
		zap.String("parent_id", parentID), 
		zap.Int("family_count", len(dtos)))
	return dtos, nil
}

// FindFamilyByChild finds the family that contains a specific child
func (s *FamilyApplicationService) FindFamilyByChild(ctx context.Context, childID string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Finding family by child ID", zap.String("child_id", childID))

	if childID == "" {
		s.logger.Warn(ctx, "Child ID is required for FindFamilyByChild")
		return nil, errors.NewValidationError("child ID is required", "childID", nil)
	}

	// Use repository to find family by child ID
	fam, err := s.familyRepo.FindByChildID(ctx, childID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "No family found for child ID", zap.String("child_id", childID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to find family by child ID", 
			zap.Error(err), 
			zap.String("child_id", childID))
		return nil, errors.NewApplicationError(errors.DatabaseErrorCode, "failed to find family by child ID", err)
	}

	// Convert domain entity to DTO
	dto := fam.ToDTO()
	s.logger.Info(ctx, "Successfully found family by child ID", 
		zap.String("child_id", childID), 
		zap.String("family_id", dto.ID))
	return &dto, nil
}

// CreateFamily creates a new family (alias for Create for backward compatibility)
func (s *FamilyApplicationService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "CreateFamily called (alias for Create)", zap.String("family_id", dto.ID))
	return s.Create(ctx, &dto)
}

// GetFamily retrieves a family by ID (alias for GetByID for backward compatibility)
func (s *FamilyApplicationService) GetFamily(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "GetFamily called (alias for GetByID)", zap.String("family_id", id))
	return s.GetByID(ctx, id)
}

// GetAllFamilies retrieves all families (alias for GetAll for backward compatibility)
func (s *FamilyApplicationService) GetAllFamilies(ctx context.Context) ([]*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "GetAllFamilies called (alias for GetAll)")
	return s.GetAll(ctx)
}

// UpdateFamily updates an existing family
func (s *FamilyApplicationService) UpdateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Updating family", zap.String("family_id", dto.ID))

	// Check if the family exists
	_, err := s.familyRepo.GetByID(ctx, dto.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get family for update", zap.Error(err), zap.String("family_id", dto.ID))
		return nil, errors.NewApplicationError(errors.DatabaseErrorCode, "failed to get family for update", err)
	}

	// Convert DTO to domain entity
	family, err := entity.FamilyFromDTO(dto)
	if err != nil {
		s.logger.Error(ctx, "Failed to convert DTO to domain entity", zap.Error(err), zap.String("family_id", dto.ID))
		return nil, errors.NewApplicationError(errors.ValidationErrorCode, "failed to convert DTO to domain entity", err)
	}

	// Save the updated family
	err = s.familyRepo.Save(ctx, family)
	if err != nil {
		s.logger.Error(ctx, "Failed to save updated family", zap.Error(err), zap.String("family_id", dto.ID))
		return nil, errors.NewApplicationError(errors.DatabaseErrorCode, "failed to save updated family", err)
	}

	// Clear cache if using caching
	if s.cache != nil {
		cacheKey := fmt.Sprintf("family:%s", dto.ID)
		s.cache.Delete(cacheKey)
	}

	// Convert the updated entity back to DTO
	resultDTO := family.ToDTO()
	s.logger.Info(ctx, "Successfully updated family", zap.String("family_id", dto.ID))
	return &resultDTO, nil
}

// DeleteFamily deletes a family by ID
func (s *FamilyApplicationService) DeleteFamily(ctx context.Context, id string) error {
	s.logger.Info(ctx, "Deleting family", zap.String("family_id", id))

	// Check if the family exists
	family, err := s.familyRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get family for deletion", zap.Error(err), zap.String("family_id", id))
		return errors.NewApplicationError(errors.DatabaseErrorCode, "failed to get family for deletion", err)
	}

	// Create a new family with the same data but with a "DELETED" status
	// We can't directly modify the status field because it's private
	parents := family.Parents()
	children := family.Children()

	// Create a new family with the same ID but with a "DELETED" status
	deletedFamily, err := entity.NewFamily(
		family.ID(),
		entity.Status("DELETED"), // Use a custom status for deleted families
		parents,
		children,
	)

	if err != nil {
		s.logger.Error(ctx, "Failed to create deleted family", zap.Error(err), zap.String("family_id", id))
		return errors.NewApplicationError(errors.InternalErrorCode, "failed to create deleted family", err)
	}

	// Save the updated family
	err = s.familyRepo.Save(ctx, deletedFamily)
	if err != nil {
		s.logger.Error(ctx, "Failed to save deleted family", zap.Error(err), zap.String("family_id", id))
		return errors.NewApplicationError(errors.DatabaseErrorCode, "failed to save deleted family", err)
	}

	// Clear cache if using caching
	if s.cache != nil {
		cacheKey := fmt.Sprintf("family:%s", id)
		s.cache.Delete(cacheKey)
	}

	s.logger.Info(ctx, "Successfully deleted family", zap.String("family_id", id))
	return nil
}

// GetID returns the service ID (implements di.ApplicationService)
func (s *FamilyApplicationService) GetID() string {
	return "family-application-service"
}
