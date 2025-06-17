package services

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/servicelib/errors"
)

// FamilyDomainService is a domain service that coordinates operations on the Family aggregate
type FamilyDomainService struct {
	repo ports.FamilyRepository
	// Could add other dependencies like logging, metrics, etc.
}

// NewFamilyDomainService creates a new FamilyDomainService
func NewFamilyDomainService(repo ports.FamilyRepository) *FamilyDomainService {
	if repo == nil {
		panic("repository cannot be nil")
	}
	return &FamilyDomainService{repo: repo}
}

// CreateFamily creates a new family
func (s *FamilyDomainService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
	// Convert DTO to domain entity
	fam, err := entity.FamilyFromDTO(dto)
	if err != nil {
		return nil, errors.NewApplicationError(err, "invalid family data", "INVALID_INPUT")
	}

	// Save to repository
	if err := s.repo.Save(ctx, fam); err != nil {
		return nil, errors.NewApplicationError(err, "failed to create family", "REPOSITORY_ERROR")
	}

	// Return the created family as DTO
	resultDTO := fam.ToDTO()
	return &resultDTO, nil
}

// GetFamily retrieves a family by ID
func (s *FamilyDomainService) GetFamily(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	if id == "" {
		return nil, errors.NewValidationError("family ID is required")
	}

	fam, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			return nil, err // Pass through not found errors
		}
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	dto := fam.ToDTO()
	return &dto, nil
}

// AddParent adds a parent to a family
func (s *FamilyDomainService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error) {
	if familyID == "" {
		return nil, errors.NewValidationError("family ID is required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			return nil, err // Pass through not found errors
		}
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Create parent entity from DTO
	p, err := entity.ParentFromDTO(parentDTO)
	if err != nil {
		return nil, errors.NewApplicationError(err, "invalid parent data", "INVALID_INPUT")
	}

	// Add parent to family
	if err := fam.AddParent(p); err != nil {
		return nil, err
	}

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		return nil, errors.NewApplicationError(err, "failed to save family", "REPOSITORY_ERROR")
	}

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	return &resultDTO, nil
}

// AddChild adds a child to a family
func (s *FamilyDomainService) AddChild(ctx context.Context, familyID string, childDTO entity.ChildDTO) (*entity.FamilyDTO, error) {
	if familyID == "" {
		return nil, errors.NewValidationError("family ID is required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			return nil, err // Pass through not found errors
		}
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Create child entity from DTO
	c, err := entity.ChildFromDTO(childDTO)
	if err != nil {
		return nil, errors.NewApplicationError(err, "invalid child data", "INVALID_INPUT")
	}

	// Add child to family
	if err := fam.AddChild(c); err != nil {
		return nil, err
	}

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		return nil, errors.NewApplicationError(err, "failed to save family", "REPOSITORY_ERROR")
	}

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	return &resultDTO, nil
}

// RemoveChild removes a child from a family
func (s *FamilyDomainService) RemoveChild(ctx context.Context, familyID string, childID string) (*entity.FamilyDTO, error) {
	if familyID == "" || childID == "" {
		return nil, errors.NewValidationError("family ID and child ID are required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			return nil, err // Pass through not found errors
		}
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Remove child from family
	if err := fam.RemoveChild(childID); err != nil {
		return nil, err
	}

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		return nil, errors.NewApplicationError(err, "failed to save family", "REPOSITORY_ERROR")
	}

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	return &resultDTO, nil
}

// MarkParentDeceased marks a parent as deceased
func (s *FamilyDomainService) MarkParentDeceased(ctx context.Context, familyID string, parentID string, deathDate time.Time) (*entity.FamilyDTO, error) {
	if familyID == "" || parentID == "" {
		return nil, errors.NewValidationError("family ID and parent ID are required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			return nil, err // Pass through not found errors
		}
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Mark parent as deceased
	if err := fam.MarkParentDeceased(parentID, deathDate); err != nil {
		return nil, err
	}

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		return nil, errors.NewApplicationError(err, "failed to save family", "REPOSITORY_ERROR")
	}

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	return &resultDTO, nil
}

// Divorce handles the divorce process
func (s *FamilyDomainService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error) {
	if familyID == "" || custodialParentID == "" {
		return nil, errors.NewValidationError("family ID and custodial parent ID are required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			return nil, err // Pass through not found errors
		}
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Process divorce
	newFam, err := fam.Divorce(custodialParentID)
	if err != nil {
		return nil, err
	}

	// Save both families in a transaction if possible
	// For now, save them sequentially
	if err := s.repo.Save(ctx, fam); err != nil {
		return nil, errors.NewApplicationError(err, "failed to save original family", "REPOSITORY_ERROR")
	}

	if err := s.repo.Save(ctx, newFam); err != nil {
		// This is a critical error - we've already updated the original family
		// In a real system, we'd use transactions to ensure atomicity
		return nil, errors.NewApplicationError(err, "failed to save new family", "REPOSITORY_ERROR")
	}

	// Return the new family as DTO
	resultDTO := newFam.ToDTO()
	return &resultDTO, nil
}
