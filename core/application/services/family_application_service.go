// Copyright (c) 2025 A Bit of Help, Inc.

package application

import (
	"context"
	"github.com/abitofhelp/servicelib/di"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	domainports "github.com/abitofhelp/family-service/core/domain/ports"
	domainservices "github.com/abitofhelp/family-service/core/domain/services"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// This file contains application services inside the core.  This includes the implementation of
// use-cases, accessing the outside world's databases and services through ports, and executing domain
// logic.

// BaseApplicationService is a generic implementation of the ApplicationService interface
type BaseApplicationService[T any, D any] struct {
	// Common dependencies and methods for all application services
}

// FamilyApplicationService implements the application service for family-related use cases
// It implements the appports.FamilyApplicationService interface
type FamilyApplicationService struct {
	BaseApplicationService[*entity.Family, *entity.FamilyDTO]
	familyService *domainservices.FamilyDomainService
	familyRepo    domainports.FamilyRepository
	logger        *logging.ContextLogger
}

// Ensure FamilyApplicationService implements di.ApplicationService
var _ di.ApplicationService = (*FamilyApplicationService)(nil)

// NewFamilyApplicationService creates a new FamilyApplicationService
func NewFamilyApplicationService(
	familyService *domainservices.FamilyDomainService,
	familyRepo domainports.FamilyRepository,
	logger *logging.ContextLogger,
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
	}
}

// Create creates a new family (implements ApplicationService.Create)
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

// GetByID retrieves a family by ID (implements ApplicationService.GetByID)
func (s *FamilyApplicationService) GetByID(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Retrieving family by ID", zap.String("family_id", id))

	// Delegate to domain service
	family, err := s.familyService.GetFamily(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to retrieve family", zap.Error(err), zap.String("family_id", id))
		return nil, err
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

// GetID returns the service ID (implements di.ApplicationService)
func (s *FamilyApplicationService) GetID() string {
	return "family-application-service"
}
