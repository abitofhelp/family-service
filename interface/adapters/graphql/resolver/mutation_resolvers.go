// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// CreateFamily is the resolver for the createFamily field.
func (r *mutationResolver) CreateFamily(ctx context.Context, input model.FamilyInput) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR"}, []string{"CREATE"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Convert input to domain DTO
	familyDTO, err := r.mapper.ToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Call service
	resultDTO, err := r.familyService.CreateFamily(ctx, familyDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to create family: %w", err)
	}

	// Convert result back to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

// AddParent is the resolver for the addParent field.
func (r *mutationResolver) AddParent(ctx context.Context, familyID identification.ID, input model.ParentInput) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR"}, []string{"CREATE"}, "PARENT"); err != nil {
		return nil, err
	}

	// Convert input to domain DTO
	parentDTO, err := r.mapper.ToParentDTO(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Call service
	resultDTO, err := r.familyService.AddParent(ctx, familyID.String(), parentDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to add parent: %w", err)
	}

	// Convert result back to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

// AddChild is the resolver for the addChild field.
func (r *mutationResolver) AddChild(ctx context.Context, familyID identification.ID, input model.ChildInput) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR"}, []string{"UPDATE"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Convert input to domain DTO
	childDTO, err := r.mapper.ToChildDTO(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Call service
	resultDTO, err := r.familyService.AddChild(ctx, familyID.String(), childDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to add child: %w", err)
	}

	// Convert result back to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

// RemoveChild is the resolver for the removeChild field.
func (r *mutationResolver) RemoveChild(ctx context.Context, familyID identification.ID, childID identification.ID) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR"}, []string{"UPDATE"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Call service
	resultDTO, err := r.familyService.RemoveChild(ctx, familyID.String(), childID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to remove child: %w", err)
	}

	// Convert result back to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

// UpdateFamily is the resolver for the updateFamily field.
func (r *mutationResolver) UpdateFamily(ctx context.Context, input model.FamilyInput) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR"}, []string{"UPDATE"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Convert input to domain DTO
	familyDTO, err := r.mapper.ToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Call service
	resultDTO, err := r.familyService.UpdateFamily(ctx, familyDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to update family: %w", err)
	}

	// Convert result back to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

// DeleteFamily is the resolver for the deleteFamily field.
func (r *mutationResolver) DeleteFamily(ctx context.Context, id identification.ID) (bool, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN"}, []string{"DELETE"}, "FAMILY"); err != nil {
		return false, err
	}

	// Call service
	if err := r.familyService.DeleteFamily(ctx, id.String()); err != nil {
		return false, fmt.Errorf("failed to delete family: %w", err)
	}

	return true, nil
}

// MarkParentDeceased is the resolver for the markParentDeceased field.
func (r *mutationResolver) MarkParentDeceased(ctx context.Context, familyID identification.ID, parentID identification.ID, deathDate string) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR"}, []string{"UPDATE"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Parse death date using RFC3339 format as required by the project guidelines
	parsedDeathDate, err := time.Parse(time.RFC3339, deathDate)
	if err != nil {
		return nil, fmt.Errorf("invalid death date format (expected RFC3339): %w", err)
	}

	// Call service
	resultDTO, err := r.familyService.MarkParentDeceased(ctx, familyID.String(), parentID.String(), parsedDeathDate)
	if err != nil {
		return nil, fmt.Errorf("failed to mark parent as deceased: %w", err)
	}

	// Convert result back to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

// Divorce is the resolver for the divorce field.
func (r *mutationResolver) Divorce(ctx context.Context, familyID identification.ID, custodialParentID identification.ID) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR"}, []string{"UPDATE"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Call service
	resultDTO, err := r.familyService.Divorce(ctx, familyID.String(), custodialParentID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to divorce family: %w", err)
	}

	// Convert result back to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}
