// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"context"
	"fmt"

	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// CountFamilies is the resolver for the countFamilies field.
func (r *queryResolver) CountFamilies(ctx context.Context) (int, error) {
	// Get all families
	families, err := r.familyService.GetAllFamilies(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get families: %w", err)
	}

	// Return the count of families
	return len(families), nil
}

// CountParents is the resolver for the countParents field.
func (r *queryResolver) CountParents(ctx context.Context) (int, error) {
	// Get all families
	families, err := r.familyService.GetAllFamilies(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get families: %w", err)
	}

	// Count unique parents
	parents := make(map[string]bool)
	for _, family := range families {
		for _, parent := range family.Parents {
			parents[parent.ID] = true
		}
	}

	// Return the count of unique parents
	return len(parents), nil
}

// FindFamiliesByParent is the resolver for the findFamiliesByParent field.
func (r *queryResolver) FindFamiliesByParent(ctx context.Context, parentID identification.ID) ([]*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR", "VIEWER"}, []string{"READ"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Call service
	resultDTOs, err := r.familyService.FindFamiliesByParent(ctx, parentID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find families by parent: %w", err)
	}

	// Convert results to GraphQL models
	var results []*model.Family
	for _, dto := range resultDTOs {
		result, err := r.mapper.ToGraphQL(*dto)
		if err != nil {
			return nil, fmt.Errorf("failed to convert result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// FindFamilyByChild is the resolver for the findFamilyByChild field.
func (r *queryResolver) FindFamilyByChild(ctx context.Context, childID identification.ID) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR", "VIEWER"}, []string{"READ"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Call service
	resultDTO, err := r.familyService.FindFamilyByChild(ctx, childID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find family by child: %w", err)
	}

	// Convert result to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

// Parents is the resolver for the parents field.
func (r *queryResolver) Parents(ctx context.Context) ([]*model.Parent, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR", "VIEWER"}, []string{"READ"}, "PARENT"); err != nil {
		return nil, err
	}

	// Get all families
	families, err := r.familyService.GetAllFamilies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get families: %w", err)
	}

	// Collect unique parents
	parentMap := make(map[string]*model.Parent)
	for _, family := range families {
		for _, parentDTO := range family.Parents {
			// Convert parent DTO to GraphQL model
			parent, err := r.mapper.ToParent(parentDTO)
			if err != nil {
				return nil, fmt.Errorf("failed to convert parent: %w", err)
			}

			// Add to map if not already present
			if _, exists := parentMap[parent.ID.String()]; !exists {
				parentMap[parent.ID.String()] = parent
			}
		}
	}

	// Convert map to slice
	var parents []*model.Parent
	for _, parent := range parentMap {
		parents = append(parents, parent)
	}

	return parents, nil
}

// CountChildren is the resolver for the countChildren field.
func (r *queryResolver) CountChildren(ctx context.Context) (int, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR", "VIEWER"}, []string{"READ"}, "CHILD"); err != nil {
		return 0, err
	}

	// Get all families
	families, err := r.familyService.GetAllFamilies(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get families: %w", err)
	}

	// Count unique children
	children := make(map[string]bool)
	for _, family := range families {
		for _, child := range family.Children {
			children[child.ID] = true
		}
	}

	return len(children), nil
}

// GetFamily is the resolver for the getFamily field.
func (r *queryResolver) GetFamily(ctx context.Context, id identification.ID) (*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR", "VIEWER"}, []string{"READ"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Call service
	resultDTO, err := r.familyService.GetFamily(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get family: %w", err)
	}

	// Convert result to GraphQL model
	result, err := r.mapper.ToGraphQL(*resultDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}

	return result, nil
}

// GetAllFamilies is the resolver for the getAllFamilies field.
func (r *queryResolver) GetAllFamilies(ctx context.Context) ([]*model.Family, error) {
	// Check authorization
	if err := checkAuthorization(ctx, []string{"ADMIN", "EDITOR", "VIEWER"}, []string{"READ"}, "FAMILY"); err != nil {
		return nil, err
	}

	// Call service
	resultDTOs, err := r.familyService.GetAllFamilies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get families: %w", err)
	}

	// Convert results to GraphQL models
	var results []*model.Family
	for _, dto := range resultDTOs {
		result, err := r.mapper.ToGraphQL(*dto)
		if err != nil {
			return nil, fmt.Errorf("failed to convert result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}
