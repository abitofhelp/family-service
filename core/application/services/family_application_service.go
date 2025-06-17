// Copyright (c) 2025 A Bit of Help, Inc.

package application

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	domainports "github.com/abitofhelp/family-service/core/domain/ports"
	domainservices "github.com/abitofhelp/family-service/core/domain/services"
	"github.com/abitofhelp/servicelib/errors"
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
	// Could add other dependencies like logging, metrics, transaction manager, etc.
}

// NewFamilyApplicationService creates a new FamilyApplicationService
func NewFamilyApplicationService(
	familyService *domainservices.FamilyDomainService,
	familyRepo domainports.FamilyRepository,
) *FamilyApplicationService {
	if familyService == nil {
		panic("family service cannot be nil")
	}
	if familyRepo == nil {
		panic("family repository cannot be nil")
	}
	return &FamilyApplicationService{
		familyService: familyService,
		familyRepo:    familyRepo,
	}
}

// Create creates a new family (implements ApplicationService.Create)
func (s *FamilyApplicationService) Create(ctx context.Context, dto *entity.FamilyDTO) (*entity.FamilyDTO, error) {
	// Delegate to domain service
	return s.familyService.CreateFamily(ctx, *dto)
}

// GetByID retrieves a family by ID (implements ApplicationService.GetByID)
func (s *FamilyApplicationService) GetByID(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	// Delegate to domain service
	return s.familyService.GetFamily(ctx, id)
}

// GetAll retrieves all families (implements ApplicationService.GetAll)
func (s *FamilyApplicationService) GetAll(ctx context.Context) ([]*entity.FamilyDTO, error) {
	// Use repository to get all families
	families, err := s.familyRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.NewApplicationError(err, "failed to get all families", "REPOSITORY_ERROR")
	}

	// Convert domain entities to DTOs
	dtos := make([]*entity.FamilyDTO, 0, len(families))
	for _, fam := range families {
		dto := fam.ToDTO()
		dtos = append(dtos, &dto)
	}

	return dtos, nil
}

// AddParent adds a parent to a family
func (s *FamilyApplicationService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error) {
	// Delegate to domain service
	return s.familyService.AddParent(ctx, familyID, parentDTO)
}

// AddChild adds a child to a family
func (s *FamilyApplicationService) AddChild(ctx context.Context, familyID string, childDTO entity.ChildDTO) (*entity.FamilyDTO, error) {
	// Delegate to domain service
	return s.familyService.AddChild(ctx, familyID, childDTO)
}

// RemoveChild removes a child from a family
func (s *FamilyApplicationService) RemoveChild(ctx context.Context, familyID string, childID string) (*entity.FamilyDTO, error) {
	// Delegate to domain service
	return s.familyService.RemoveChild(ctx, familyID, childID)
}

// MarkParentDeceased marks a parent as deceased
func (s *FamilyApplicationService) MarkParentDeceased(ctx context.Context, familyID string, parentID string, deathDate time.Time) (*entity.FamilyDTO, error) {
	// Delegate to domain service
	return s.familyService.MarkParentDeceased(ctx, familyID, parentID, deathDate)
}

// Divorce handles the divorce process
func (s *FamilyApplicationService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error) {
	// Delegate to domain service
	return s.familyService.Divorce(ctx, familyID, custodialParentID)
}

// FindFamiliesByParent finds families that contain a specific parent
func (s *FamilyApplicationService) FindFamiliesByParent(ctx context.Context, parentID string) ([]*entity.FamilyDTO, error) {
	if parentID == "" {
		return nil, errors.NewValidationError("parent ID is required")
	}

	// Use repository to find families by parent ID
	families, err := s.familyRepo.FindByParentID(ctx, parentID)
	if err != nil {
		return nil, errors.NewApplicationError(err, "failed to find families by parent ID", "REPOSITORY_ERROR")
	}

	// Convert domain entities to DTOs
	dtos := make([]*entity.FamilyDTO, 0, len(families))
	for _, fam := range families {
		dto := fam.ToDTO()
		dtos = append(dtos, &dto)
	}

	return dtos, nil
}

// FindFamilyByChild finds the family that contains a specific child
func (s *FamilyApplicationService) FindFamilyByChild(ctx context.Context, childID string) (*entity.FamilyDTO, error) {
	if childID == "" {
		return nil, errors.NewValidationError("child ID is required")
	}

	// Use repository to find family by child ID
	fam, err := s.familyRepo.FindByChildID(ctx, childID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			return nil, err // Pass through not found errors
		}
		return nil, errors.NewApplicationError(err, "failed to find family by child ID", "REPOSITORY_ERROR")
	}

	// Convert domain entity to DTO
	dto := fam.ToDTO()
	return &dto, nil
}

// CreateFamily creates a new family (alias for Create for backward compatibility)
func (s *FamilyApplicationService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
	return s.Create(ctx, &dto)
}

// GetFamily retrieves a family by ID (alias for GetByID for backward compatibility)
func (s *FamilyApplicationService) GetFamily(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	return s.GetByID(ctx, id)
}

// GetAllFamilies retrieves all families (alias for GetAll for backward compatibility)
func (s *FamilyApplicationService) GetAllFamilies(ctx context.Context) ([]*entity.FamilyDTO, error) {
	return s.GetAll(ctx)
}
