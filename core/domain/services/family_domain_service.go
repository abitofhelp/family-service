// Copyright (c) 2025 A Bit of Help, Inc.

package services

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// FamilyDomainService is a domain service that coordinates operations on the Family aggregate
type FamilyDomainService struct {
	repo   ports.FamilyRepository
	logger *logging.ContextLogger
}

// NewFamilyDomainService creates a new FamilyDomainService
func NewFamilyDomainService(repo ports.FamilyRepository, logger *logging.ContextLogger) *FamilyDomainService {
	if repo == nil {
		panic("repository cannot be nil")
	}
	if logger == nil {
		panic("logger cannot be nil")
	}
	return &FamilyDomainService{
		repo:   repo,
		logger: logger,
	}
}

// CreateFamily creates a new family
func (s *FamilyDomainService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Creating new family in domain service", zap.String("family_id", dto.ID), zap.String("status", dto.Status))

	// Convert DTO to domain entity
	fam, err := entity.FamilyFromDTO(dto)
	if err != nil {
		s.logger.Error(ctx, "Invalid family data", zap.Error(err), zap.String("family_id", dto.ID))
		return nil, errors.NewApplicationError(err, "invalid family data", "INVALID_INPUT")
	}

	// Save to repository
	if err := s.repo.Save(ctx, fam); err != nil {
		s.logger.Error(ctx, "Failed to save family to repository", zap.Error(err), zap.String("family_id", fam.ID()))
		return nil, errors.NewApplicationError(err, "failed to create family", "REPOSITORY_ERROR")
	}

	// Return the created family as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully created family in domain service", 
		zap.String("family_id", resultDTO.ID), 
		zap.Int("parent_count", resultDTO.ParentCount), 
		zap.Int("children_count", resultDTO.ChildrenCount))
	return &resultDTO, nil
}

// GetFamily retrieves a family by ID
func (s *FamilyDomainService) GetFamily(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Retrieving family by ID in domain service", zap.String("family_id", id))

	if id == "" {
		s.logger.Warn(ctx, "Family ID is required")
		return nil, errors.NewValidationError("family ID is required")
	}

	fam, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found", zap.String("family_id", id))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family from repository", zap.Error(err), zap.String("family_id", id))
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	dto := fam.ToDTO()
	s.logger.Info(ctx, "Successfully retrieved family in domain service", 
		zap.String("family_id", dto.ID), 
		zap.String("status", dto.Status))
	return &dto, nil
}

// AddParent adds a parent to a family
func (s *FamilyDomainService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Adding parent to family in domain service", 
		zap.String("family_id", familyID), 
		zap.String("parent_id", parentDTO.ID),
		zap.String("parent_first_name", parentDTO.FirstName),
		zap.String("parent_last_name", parentDTO.LastName))

	if familyID == "" {
		s.logger.Warn(ctx, "Family ID is required for AddParent")
		return nil, errors.NewValidationError("family ID is required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for AddParent", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for AddParent", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Create parent entity from DTO
	p, err := entity.ParentFromDTO(parentDTO)
	if err != nil {
		s.logger.Error(ctx, "Invalid parent data for AddParent", 
			zap.Error(err), 
			zap.String("parent_id", parentDTO.ID))
		return nil, errors.NewApplicationError(err, "invalid parent data", "INVALID_INPUT")
	}

 	// Add parent to family
	if err := fam.AddParent(p); err != nil {
		s.logger.Error(ctx, "Failed to add parent to family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("parent_id", p.ID()))
		return nil, err
	}

	// Update family status if needed
	// If we're adding a second parent to a SINGLE family, change status to MARRIED
	if fam.Status() == entity.Single && len(fam.Parents()) == 2 {
		s.logger.Info(ctx, "Updating family status from SINGLE to MARRIED", zap.String("family_id", familyID))
		// Create a new family with updated status
		updatedFam, err := entity.NewFamily(fam.ID(), entity.Married, fam.Parents(), fam.Children())
		if err != nil {
			s.logger.Error(ctx, "Failed to update family status", 
				zap.Error(err), 
				zap.String("family_id", familyID))
			return nil, errors.NewApplicationError(err, "failed to update family status", "DOMAIN_ERROR")
		}
		fam = updatedFam
	}

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		s.logger.Error(ctx, "Failed to save family after adding parent", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to save family", "REPOSITORY_ERROR")
	}

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully added parent to family", 
		zap.String("family_id", resultDTO.ID), 
		zap.Int("parent_count", resultDTO.ParentCount),
		zap.String("status", resultDTO.Status))
	return &resultDTO, nil
}

// AddChild adds a child to a family
func (s *FamilyDomainService) AddChild(ctx context.Context, familyID string, childDTO entity.ChildDTO) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Adding child to family in domain service", 
		zap.String("family_id", familyID), 
		zap.String("child_id", childDTO.ID),
		zap.String("child_first_name", childDTO.FirstName),
		zap.String("child_last_name", childDTO.LastName))

	if familyID == "" {
		s.logger.Warn(ctx, "Family ID is required for AddChild")
		return nil, errors.NewValidationError("family ID is required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for AddChild", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for AddChild", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Create child entity from DTO
	c, err := entity.ChildFromDTO(childDTO)
	if err != nil {
		s.logger.Error(ctx, "Invalid child data for AddChild", 
			zap.Error(err), 
			zap.String("child_id", childDTO.ID))
		return nil, errors.NewApplicationError(err, "invalid child data", "INVALID_INPUT")
	}

	// Add child to family
	if err := fam.AddChild(c); err != nil {
		s.logger.Error(ctx, "Failed to add child to family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("child_id", c.ID()))
		return nil, err
	}

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		s.logger.Error(ctx, "Failed to save family after adding child", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to save family", "REPOSITORY_ERROR")
	}

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully added child to family", 
		zap.String("family_id", resultDTO.ID), 
		zap.Int("children_count", resultDTO.ChildrenCount))
	return &resultDTO, nil
}

// RemoveChild removes a child from a family
func (s *FamilyDomainService) RemoveChild(ctx context.Context, familyID string, childID string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Removing child from family in domain service", 
		zap.String("family_id", familyID), 
		zap.String("child_id", childID))

	if familyID == "" || childID == "" {
		s.logger.Warn(ctx, "Family ID and child ID are required for RemoveChild", 
			zap.String("family_id", familyID), 
			zap.String("child_id", childID))
		return nil, errors.NewValidationError("family ID and child ID are required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for RemoveChild", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for RemoveChild", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Remove child from family
	if err := fam.RemoveChild(childID); err != nil {
		s.logger.Error(ctx, "Failed to remove child from family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("child_id", childID))
		return nil, err
	}

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		s.logger.Error(ctx, "Failed to save family after removing child", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to save family", "REPOSITORY_ERROR")
	}

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully removed child from family", 
		zap.String("family_id", resultDTO.ID), 
		zap.Int("children_count", resultDTO.ChildrenCount))
	return &resultDTO, nil
}

// MarkParentDeceased marks a parent as deceased
func (s *FamilyDomainService) MarkParentDeceased(ctx context.Context, familyID string, parentID string, deathDate time.Time) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Marking parent as deceased in domain service", 
		zap.String("family_id", familyID), 
		zap.String("parent_id", parentID),
		zap.Time("death_date", deathDate))

	if familyID == "" || parentID == "" {
		s.logger.Warn(ctx, "Family ID and parent ID are required for MarkParentDeceased", 
			zap.String("family_id", familyID), 
			zap.String("parent_id", parentID))
		return nil, errors.NewValidationError("family ID and parent ID are required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for MarkParentDeceased", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for MarkParentDeceased", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Mark parent as deceased
	if err := fam.MarkParentDeceased(parentID, deathDate); err != nil {
		s.logger.Error(ctx, "Failed to mark parent as deceased", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("parent_id", parentID))
		return nil, err
	}

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		s.logger.Error(ctx, "Failed to save family after marking parent as deceased", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to save family", "REPOSITORY_ERROR")
	}

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully marked parent as deceased", 
		zap.String("family_id", resultDTO.ID), 
		zap.String("status", resultDTO.Status))
	return &resultDTO, nil
}

// Divorce handles the divorce process
func (s *FamilyDomainService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error) {
	s.logger.Info(ctx, "Processing divorce in domain service", 
		zap.String("family_id", familyID), 
		zap.String("custodial_parent_id", custodialParentID))

	if familyID == "" || custodialParentID == "" {
		s.logger.Warn(ctx, "Family ID and custodial parent ID are required for Divorce", 
			zap.String("family_id", familyID), 
			zap.String("custodial_parent_id", custodialParentID))
		return nil, errors.NewValidationError("family ID and custodial parent ID are required")
	}

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for Divorce", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for Divorce", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewApplicationError(err, "failed to retrieve family", "REPOSITORY_ERROR")
	}

	// Process divorce
	// Note: After our changes, fam.Divorce() now returns the new family with the remaining parent
	// The original family (fam) is modified in place to keep the custodial parent and children
	remainingFam, err := fam.Divorce(custodialParentID)
	if err != nil {
		s.logger.Error(ctx, "Failed to process divorce", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("custodial_parent_id", custodialParentID))
		return nil, err
	}

	s.logger.Info(ctx, "Divorce processed, saving family with custodial parent", 
		zap.String("family_id", fam.ID()), 
		zap.String("status", string(fam.Status())))

	// Save both families in a transaction if possible
	// For now, save them sequentially
	if err := s.repo.Save(ctx, fam); err != nil {
		s.logger.Error(ctx, "Failed to save family with custodial parent after divorce", 
			zap.Error(err), 
			zap.String("family_id", fam.ID()))
		return nil, errors.NewApplicationError(err, "failed to save family with custodial parent", "REPOSITORY_ERROR")
	}

	s.logger.Info(ctx, "Family with custodial parent saved, saving family with remaining parent", 
		zap.String("family_id", remainingFam.ID()), 
		zap.String("status", string(remainingFam.Status())))

	if err := s.repo.Save(ctx, remainingFam); err != nil {
		// This is a critical error - we've already updated the family with custodial parent
		// In a real system, we'd use transactions to ensure atomicity
		s.logger.Error(ctx, "Failed to save family with remaining parent after divorce - CRITICAL ERROR", 
			zap.Error(err), 
			zap.String("family_id", remainingFam.ID()),
			zap.String("custodial_parent_family_id", fam.ID()))
		return nil, errors.NewApplicationError(err, "failed to save family with remaining parent", "REPOSITORY_ERROR")
	}

	// Return the original family (now with custodial parent and children) as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully processed divorce", 
		zap.String("family_id", resultDTO.ID), 
		zap.String("status", resultDTO.Status),
		zap.Int("children_count", resultDTO.ChildrenCount))
	return &resultDTO, nil
}
