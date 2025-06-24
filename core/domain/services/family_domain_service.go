// Copyright (c) 2025 A Bit of Help, Inc.

package services

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/metrics"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FamilyDomainService is a domain service that coordinates operations on the Family aggregate
type FamilyDomainService struct {
	repo   ports.FamilyRepository
	logger *logging.ContextLogger
	tracer trace.Tracer
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
		tracer: otel.Tracer("family-domain-service"),
	}
}

// CreateFamily creates a new family
func (s *FamilyDomainService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
	// Start a new span for this operation
	ctx, span := s.tracer.Start(ctx, "FamilyDomainService.CreateFamily")
	defer span.End()

	// Start timer for operation duration
	startTime := time.Now()

	// Log operation start
	s.logger.Info(ctx, "Creating new family in domain service", zap.String("family_id", dto.ID), zap.String("status", dto.Status))

	// Convert DTO to domain entity
	fam, err := entity.FamilyFromDTO(dto)
	if err != nil {
		// Record metrics for failure
		metrics.FamilyOperationsTotal.WithLabelValues("create_family", metrics.StatusFailure).Inc()

		s.logger.Error(ctx, "Invalid family data", zap.Error(err), zap.String("family_id", dto.ID))
		return nil, errors.NewValidationError("invalid family data", "family", err)
	}

	// Create a span for repository operation
	ctx, saveSpan := s.tracer.Start(ctx, "Repository.Save")

	// Save to repository
	if err := s.repo.Save(ctx, fam); err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusFailure).Inc()
		saveSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("create_family", metrics.StatusFailure).Inc()

		s.logger.Error(ctx, "Failed to save family to repository", zap.Error(err), zap.String("family_id", fam.ID()))
		return nil, errors.NewDatabaseError("failed to create family", "save", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("save").Observe(time.Since(startTime).Seconds())
	saveSpan.End()

	// Update family member counts
	metrics.FamilyMemberCounts.WithLabelValues("parents").Add(float64(len(fam.Parents())))
	metrics.FamilyMemberCounts.WithLabelValues("children").Add(float64(len(fam.Children())))

	// Update family status counts
	metrics.FamilyStatusCounts.WithLabelValues(string(fam.Status())).Inc()

	// Record metrics for operation success
	metrics.FamilyOperationsTotal.WithLabelValues("create_family", metrics.StatusSuccess).Inc()
	metrics.FamilyOperationsDuration.WithLabelValues("create_family").Observe(time.Since(startTime).Seconds())

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
	// Start a new span for this operation
	ctx, span := s.tracer.Start(ctx, "FamilyDomainService.GetFamily")
	defer span.End()

	// Start timer for operation duration
	startTime := time.Now()

	s.logger.Info(ctx, "Retrieving family by ID in domain service", zap.String("family_id", id))

	if id == "" {
		// Record metrics for failure
		metrics.FamilyOperationsTotal.WithLabelValues("get_family", metrics.StatusFailure).Inc()

		s.logger.Warn(ctx, "Family ID is required")
		return nil, errors.NewValidationError("family ID is required", "id", nil)
	}

	// Create a span for repository operation
	ctx, getSpan := s.tracer.Start(ctx, "Repository.GetByID")

	fam, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusFailure).Inc()
		getSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("get_family", metrics.StatusFailure).Inc()

		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found", zap.String("family_id", id))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family from repository", zap.Error(err), zap.String("family_id", id))
		return nil, errors.NewDatabaseError("failed to retrieve family", "query", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("get_by_id").Observe(time.Since(startTime).Seconds())
	getSpan.End()

	// Record metrics for operation success
	metrics.FamilyOperationsTotal.WithLabelValues("get_family", metrics.StatusSuccess).Inc()
	metrics.FamilyOperationsDuration.WithLabelValues("get_family").Observe(time.Since(startTime).Seconds())

	dto := fam.ToDTO()
	s.logger.Info(ctx, "Successfully retrieved family in domain service", 
		zap.String("family_id", dto.ID), 
		zap.String("status", dto.Status))
	return &dto, nil
}

// AddParent adds a parent to a family
func (s *FamilyDomainService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error) {
	// Start a new span for this operation
	ctx, span := s.tracer.Start(ctx, "FamilyDomainService.AddParent")
	defer span.End()

	// Start timer for operation duration
	startTime := time.Now()

	s.logger.Info(ctx, "Adding parent to family in domain service", 
		zap.String("family_id", familyID), 
		zap.String("parent_id", parentDTO.ID),
		zap.String("parent_first_name", parentDTO.FirstName),
		zap.String("parent_last_name", parentDTO.LastName))

	if familyID == "" {
		// Record metrics for failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_parent", metrics.StatusFailure).Inc()

		s.logger.Warn(ctx, "Family ID is required for AddParent")
		return nil, errors.NewValidationError("family ID is required", "familyID", nil)
	}

	// Create a span for retrieving the family
	ctx, getSpan := s.tracer.Start(ctx, "Repository.GetByID.AddParent")

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusFailure).Inc()
		getSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_parent", metrics.StatusFailure).Inc()

		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for AddParent", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for AddParent", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to retrieve family", "query", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("get_by_id").Observe(time.Since(startTime).Seconds())
	getSpan.End()

	// Create a span for parent creation and validation
	ctx, parentSpan := s.tracer.Start(ctx, "Domain.CreateParent")

	// Create parent entity from DTO
	p, err := entity.ParentFromDTO(parentDTO)
	if err != nil {
		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_parent", metrics.StatusFailure).Inc()
		parentSpan.End()

		s.logger.Error(ctx, "Invalid parent data for AddParent", 
			zap.Error(err), 
			zap.String("parent_id", parentDTO.ID))
		return nil, errors.NewValidationError("invalid parent data", "parent", err)
	}

	parentSpan.End()

	// Create a span for adding parent to family
	ctx, addParentSpan := s.tracer.Start(ctx, "Domain.AddParentToFamily")

 	// Add parent to family
	if err := fam.AddParent(p); err != nil {
		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_parent", metrics.StatusFailure).Inc()
		addParentSpan.End()

		s.logger.Error(ctx, "Failed to add parent to family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("parent_id", p.ID()))
		return nil, err
	}

	addParentSpan.End()

	// Create a span for updating family status
	ctx, updateStatusSpan := s.tracer.Start(ctx, "Domain.UpdateFamilyStatus")

	// Update family status if needed
	// If we're adding a second parent to a SINGLE family, change status to MARRIED
	statusChanged := false
	if fam.Status() == entity.Single && len(fam.Parents()) == 2 {
		statusChanged = true
		s.logger.Info(ctx, "Updating family status from SINGLE to MARRIED", zap.String("family_id", familyID))
		// Create a new family with updated status
		updatedFam, err := entity.NewFamily(fam.ID(), entity.Married, fam.Parents(), fam.Children())
		if err != nil {
			// Record metrics for operation failure
			metrics.FamilyOperationsTotal.WithLabelValues("add_parent", metrics.StatusFailure).Inc()
			updateStatusSpan.End()

			s.logger.Error(ctx, "Failed to update family status", 
				zap.Error(err), 
				zap.String("family_id", familyID))
			return nil, errors.NewDomainError(errors.BusinessRuleViolationCode, "failed to update family status", err)
		}
		fam = updatedFam
	}

	updateStatusSpan.End()

	// Create a span for saving the family
	ctx, saveSpan := s.tracer.Start(ctx, "Repository.Save.AddParent")

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusFailure).Inc()
		saveSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_parent", metrics.StatusFailure).Inc()

		s.logger.Error(ctx, "Failed to save family after adding parent", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to save family", "save", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("save").Observe(time.Since(startTime).Seconds())
	saveSpan.End()

	// Update family member counts
	metrics.FamilyMemberCounts.WithLabelValues("parents").Inc()

	// Update family status counts if status changed
	if statusChanged {
		metrics.FamilyStatusCounts.WithLabelValues("single").Dec()
		metrics.FamilyStatusCounts.WithLabelValues("married").Inc()
	}

	// Record metrics for operation success
	metrics.FamilyOperationsTotal.WithLabelValues("add_parent", metrics.StatusSuccess).Inc()
	metrics.FamilyOperationsDuration.WithLabelValues("add_parent").Observe(time.Since(startTime).Seconds())

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
	// Start a new span for this operation
	ctx, span := s.tracer.Start(ctx, "FamilyDomainService.AddChild")
	defer span.End()

	// Start timer for operation duration
	startTime := time.Now()

	s.logger.Info(ctx, "Adding child to family in domain service", 
		zap.String("family_id", familyID), 
		zap.String("child_id", childDTO.ID),
		zap.String("child_first_name", childDTO.FirstName),
		zap.String("child_last_name", childDTO.LastName))

	if familyID == "" {
		// Record metrics for failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_child", metrics.StatusFailure).Inc()

		s.logger.Warn(ctx, "Family ID is required for AddChild")
		return nil, errors.NewValidationError("family ID is required", "familyID", nil)
	}

	// Create a span for retrieving the family
	ctx, getSpan := s.tracer.Start(ctx, "Repository.GetByID.AddChild")

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusFailure).Inc()
		getSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_child", metrics.StatusFailure).Inc()

		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for AddChild", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for AddChild", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to retrieve family", "query", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("get_by_id").Observe(time.Since(startTime).Seconds())
	getSpan.End()

	// Create a span for child creation and validation
	ctx, childSpan := s.tracer.Start(ctx, "Domain.CreateChild")

	// Create child entity from DTO
	c, err := entity.ChildFromDTO(childDTO)
	if err != nil {
		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_child", metrics.StatusFailure).Inc()
		childSpan.End()

		s.logger.Error(ctx, "Invalid child data for AddChild", 
			zap.Error(err), 
			zap.String("child_id", childDTO.ID))
		return nil, errors.NewValidationError("invalid child data", "child", err)
	}

	childSpan.End()

	// Create a span for adding child to family
	ctx, addChildSpan := s.tracer.Start(ctx, "Domain.AddChildToFamily")

	// Add child to family
	if err := fam.AddChild(c); err != nil {
		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_child", metrics.StatusFailure).Inc()
		addChildSpan.End()

		s.logger.Error(ctx, "Failed to add child to family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("child_id", c.ID()))
		return nil, err
	}

	addChildSpan.End()

	// Create a span for saving the family
	ctx, saveSpan := s.tracer.Start(ctx, "Repository.Save.AddChild")

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusFailure).Inc()
		saveSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("add_child", metrics.StatusFailure).Inc()

		s.logger.Error(ctx, "Failed to save family after adding child", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to save family", "save", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("save").Observe(time.Since(startTime).Seconds())
	saveSpan.End()

	// Update family member counts
	metrics.FamilyMemberCounts.WithLabelValues("children").Inc()

	// Record metrics for operation success
	metrics.FamilyOperationsTotal.WithLabelValues("add_child", metrics.StatusSuccess).Inc()
	metrics.FamilyOperationsDuration.WithLabelValues("add_child").Observe(time.Since(startTime).Seconds())

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully added child to family", 
		zap.String("family_id", resultDTO.ID), 
		zap.Int("children_count", resultDTO.ChildrenCount))
	return &resultDTO, nil
}

// RemoveChild removes a child from a family
func (s *FamilyDomainService) RemoveChild(ctx context.Context, familyID string, childID string) (*entity.FamilyDTO, error) {
	// Start a new span for this operation
	ctx, span := s.tracer.Start(ctx, "FamilyDomainService.RemoveChild")
	defer span.End()

	// Start timer for operation duration
	startTime := time.Now()

	s.logger.Info(ctx, "Removing child from family in domain service", 
		zap.String("family_id", familyID), 
		zap.String("child_id", childID))

	if familyID == "" || childID == "" {
		// Record metrics for failure
		metrics.FamilyOperationsTotal.WithLabelValues("remove_child", metrics.StatusFailure).Inc()

		s.logger.Warn(ctx, "Family ID and child ID are required for RemoveChild", 
			zap.String("family_id", familyID), 
			zap.String("child_id", childID))
		return nil, errors.NewValidationError("family ID and child ID are required", "familyID/childID", nil)
	}

	// Create a span for retrieving the family
	ctx, getSpan := s.tracer.Start(ctx, "Repository.GetByID.RemoveChild")

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusFailure).Inc()
		getSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("remove_child", metrics.StatusFailure).Inc()

		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for RemoveChild", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for RemoveChild", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to retrieve family", "query", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("get_by_id").Observe(time.Since(startTime).Seconds())
	getSpan.End()

	// Create a span for removing child from family
	ctx, removeChildSpan := s.tracer.Start(ctx, "Domain.RemoveChildFromFamily")

	// Remove child from family
	if err := fam.RemoveChild(childID); err != nil {
		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("remove_child", metrics.StatusFailure).Inc()
		removeChildSpan.End()

		s.logger.Error(ctx, "Failed to remove child from family", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("child_id", childID))
		return nil, err
	}

	removeChildSpan.End()

	// Create a span for saving the family
	ctx, saveSpan := s.tracer.Start(ctx, "Repository.Save.RemoveChild")

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusFailure).Inc()
		saveSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("remove_child", metrics.StatusFailure).Inc()

		s.logger.Error(ctx, "Failed to save family after removing child", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to save family", "save", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("save").Observe(time.Since(startTime).Seconds())
	saveSpan.End()

	// Update family member counts
	metrics.FamilyMemberCounts.WithLabelValues("children").Dec()

	// Record metrics for operation success
	metrics.FamilyOperationsTotal.WithLabelValues("remove_child", metrics.StatusSuccess).Inc()
	metrics.FamilyOperationsDuration.WithLabelValues("remove_child").Observe(time.Since(startTime).Seconds())

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully removed child from family", 
		zap.String("family_id", resultDTO.ID), 
		zap.Int("children_count", resultDTO.ChildrenCount))
	return &resultDTO, nil
}

// MarkParentDeceased marks a parent as deceased
func (s *FamilyDomainService) MarkParentDeceased(ctx context.Context, familyID string, parentID string, deathDate time.Time) (*entity.FamilyDTO, error) {
	// Start a new span for this operation
	ctx, span := s.tracer.Start(ctx, "FamilyDomainService.MarkParentDeceased")
	defer span.End()

	// Start timer for operation duration
	startTime := time.Now()

	s.logger.Info(ctx, "Marking parent as deceased in domain service", 
		zap.String("family_id", familyID), 
		zap.String("parent_id", parentID),
		zap.Time("death_date", deathDate))

	if familyID == "" || parentID == "" {
		// Record metrics for failure
		metrics.FamilyOperationsTotal.WithLabelValues("mark_parent_deceased", metrics.StatusFailure).Inc()

		s.logger.Warn(ctx, "Family ID and parent ID are required for MarkParentDeceased", 
			zap.String("family_id", familyID), 
			zap.String("parent_id", parentID))
		return nil, errors.NewValidationError("family ID and parent ID are required", "familyID/parentID", nil)
	}

	// Create a span for retrieving the family
	ctx, getSpan := s.tracer.Start(ctx, "Repository.GetByID.MarkParentDeceased")

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusFailure).Inc()
		getSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("mark_parent_deceased", metrics.StatusFailure).Inc()

		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for MarkParentDeceased", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for MarkParentDeceased", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to retrieve family", "query", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("get_by_id").Observe(time.Since(startTime).Seconds())
	getSpan.End()

	// Create a span for marking parent as deceased
	ctx, markDeceasedSpan := s.tracer.Start(ctx, "Domain.MarkParentDeceased")

	// Check if this will change the family status (for metrics)
	originalStatus := fam.Status()

	// Mark parent as deceased
	if err := fam.MarkParentDeceased(parentID, deathDate); err != nil {
		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("mark_parent_deceased", metrics.StatusFailure).Inc()
		markDeceasedSpan.End()

		s.logger.Error(ctx, "Failed to mark parent as deceased", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("parent_id", parentID))
		return nil, err
	}

	markDeceasedSpan.End()

	// Create a span for saving the family
	ctx, saveSpan := s.tracer.Start(ctx, "Repository.Save.MarkParentDeceased")

	// Save updated family
	if err := s.repo.Save(ctx, fam); err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusFailure).Inc()
		saveSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("mark_parent_deceased", metrics.StatusFailure).Inc()

		s.logger.Error(ctx, "Failed to save family after marking parent as deceased", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to save family", "save", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("save").Observe(time.Since(startTime).Seconds())
	saveSpan.End()

	// Update family status counts if status changed
	if originalStatus != fam.Status() {
		// If status changed to widowed, update metrics
		if fam.Status() == entity.Widowed {
			metrics.FamilyStatusCounts.WithLabelValues("widowed").Inc()
			if originalStatus == entity.Married {
				metrics.FamilyStatusCounts.WithLabelValues("married").Dec()
			} else if originalStatus == entity.Single {
				metrics.FamilyStatusCounts.WithLabelValues("single").Dec()
			}
		}
	}

	// Record metrics for operation success
	metrics.FamilyOperationsTotal.WithLabelValues("mark_parent_deceased", metrics.StatusSuccess).Inc()
	metrics.FamilyOperationsDuration.WithLabelValues("mark_parent_deceased").Observe(time.Since(startTime).Seconds())

	// Return updated family as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully marked parent as deceased", 
		zap.String("family_id", resultDTO.ID), 
		zap.String("status", resultDTO.Status))
	return &resultDTO, nil
}

// Divorce handles the divorce process
func (s *FamilyDomainService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error) {
	// Start a new span for this operation
	ctx, span := s.tracer.Start(ctx, "FamilyDomainService.Divorce")
	defer span.End()

	// Start timer for operation duration
	startTime := time.Now()

	s.logger.Info(ctx, "Processing divorce in domain service", 
		zap.String("family_id", familyID), 
		zap.String("custodial_parent_id", custodialParentID))

	if familyID == "" || custodialParentID == "" {
		// Record metrics for failure
		metrics.FamilyOperationsTotal.WithLabelValues("divorce", metrics.StatusFailure).Inc()

		s.logger.Warn(ctx, "Family ID and custodial parent ID are required for Divorce", 
			zap.String("family_id", familyID), 
			zap.String("custodial_parent_id", custodialParentID))
		return nil, errors.NewValidationError("family ID and custodial parent ID are required", "familyID/custodialParentID", nil)
	}

	// Create a span for retrieving the family
	ctx, getSpan := s.tracer.Start(ctx, "Repository.GetByID.Divorce")

	// Get the family
	fam, err := s.repo.GetByID(ctx, familyID)
	if err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusFailure).Inc()
		getSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("divorce", metrics.StatusFailure).Inc()

		if _, ok := err.(*errors.NotFoundError); ok {
			s.logger.Info(ctx, "Family not found for Divorce", zap.String("family_id", familyID))
			return nil, err // Pass through not found errors
		}
		s.logger.Error(ctx, "Failed to retrieve family for Divorce", 
			zap.Error(err), 
			zap.String("family_id", familyID))
		return nil, errors.NewDatabaseError("failed to retrieve family", "query", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("get_by_id", metrics.StatusSuccess).Inc()
	metrics.RepositoryOperationsDuration.WithLabelValues("get_by_id").Observe(time.Since(startTime).Seconds())
	getSpan.End()

	// Create a span for the domain logic of divorce
	ctx, divorceLogicSpan := s.tracer.Start(ctx, "Domain.DivorceLogic")

	// Process divorce
	// Note: After our changes, fam.Divorce() now returns the new family with the remaining parent
	// The original family (fam) is modified in place to keep the custodial parent and children
	remainingFam, err := fam.Divorce(custodialParentID)
	if err != nil {
		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("divorce", metrics.StatusFailure).Inc()
		divorceLogicSpan.End()

		s.logger.Error(ctx, "Failed to process divorce", 
			zap.Error(err), 
			zap.String("family_id", familyID), 
			zap.String("custodial_parent_id", custodialParentID))
		return nil, err
	}

	divorceLogicSpan.End()

	s.logger.Info(ctx, "Divorce processed, saving family with custodial parent", 
		zap.String("family_id", fam.ID()), 
		zap.String("status", string(fam.Status())))

	// Create a span for saving the custodial parent family
	ctx, saveCustodialSpan := s.tracer.Start(ctx, "Repository.Save.CustodialFamily")

	// Save both families in a transaction if possible
	// For now, save them sequentially
	if err := s.repo.Save(ctx, fam); err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusFailure).Inc()
		saveCustodialSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("divorce", metrics.StatusFailure).Inc()

		s.logger.Error(ctx, "Failed to save family with custodial parent after divorce", 
			zap.Error(err), 
			zap.String("family_id", fam.ID()))
		return nil, errors.NewDatabaseError("failed to save family with custodial parent", "save", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusSuccess).Inc()
	saveCustodialSpan.End()

	s.logger.Info(ctx, "Family with custodial parent saved, saving family with remaining parent", 
		zap.String("family_id", remainingFam.ID()), 
		zap.String("status", string(remainingFam.Status())))

	// Create a span for saving the remaining parent family
	ctx, saveRemainingSpan := s.tracer.Start(ctx, "Repository.Save.RemainingFamily")

	if err := s.repo.Save(ctx, remainingFam); err != nil {
		// Record metrics for repository operation failure
		metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusFailure).Inc()
		saveRemainingSpan.End()

		// Record metrics for operation failure
		metrics.FamilyOperationsTotal.WithLabelValues("divorce", metrics.StatusFailure).Inc()

		// This is a critical error - we've already updated the family with custodial parent
		// In a real system, we'd use transactions to ensure atomicity
		s.logger.Error(ctx, "Failed to save family with remaining parent after divorce - CRITICAL ERROR", 
			zap.Error(err), 
			zap.String("family_id", remainingFam.ID()),
			zap.String("custodial_parent_family_id", fam.ID()))
		return nil, errors.NewDatabaseError("failed to save family with remaining parent", "save", "families", err)
	}

	// Record metrics for repository operation success
	metrics.RepositoryOperationsTotal.WithLabelValues("save", metrics.StatusSuccess).Inc()
	saveRemainingSpan.End()

	// Update family status counts - one family became divorced, one became single
	metrics.FamilyStatusCounts.WithLabelValues("divorced").Inc()
	metrics.FamilyStatusCounts.WithLabelValues("single").Inc()
	// If the original family was married, decrement that count
	if fam.Status() == entity.Divorced {
		metrics.FamilyStatusCounts.WithLabelValues("married").Dec()
	}

	// Record metrics for operation success
	metrics.FamilyOperationsTotal.WithLabelValues("divorce", metrics.StatusSuccess).Inc()
	metrics.FamilyOperationsDuration.WithLabelValues("divorce").Observe(time.Since(startTime).Seconds())

	// Return the original family (now with custodial parent and children) as DTO
	resultDTO := fam.ToDTO()
	s.logger.Info(ctx, "Successfully processed divorce", 
		zap.String("family_id", resultDTO.ID), 
		zap.String("status", resultDTO.Status),
		zap.Int("children_count", resultDTO.ChildrenCount))
	return &resultDTO, nil
}
