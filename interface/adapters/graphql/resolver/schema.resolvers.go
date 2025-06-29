package resolver

import (
	"context"
	"fmt"
	"time"

	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// ParentCount is the resolver for the parentCount field.
func (r *familyResolver) ParentCount(ctx context.Context, obj *model.Family) (int, error) {
	return len(obj.Parents), nil
}

// ChildrenCount is the resolver for the childrenCount field.
func (r *familyResolver) ChildrenCount(ctx context.Context, obj *model.Family) (int, error) {
	return len(obj.Children), nil
}

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

func (r *familyResolver) Parents(ctx context.Context, obj *model.Family) ([]*model.Parent, error) {
	return obj.Parents, nil
}

func (r *familyResolver) Children(ctx context.Context, obj *model.Family) ([]*model.Child, error) {
	return obj.Children, nil
}

// Helper function to check authorization
func checkAuthorization(ctx context.Context, roles []string, actions []string, resource string) error {
	// TODO: Implement proper authorization check
	// This is a placeholder that always returns nil
	// In a real implementation, this would check the user's roles and permissions
	return nil
}

type familyResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *Resolver) Family() generated.FamilyResolver {
	return &familyResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}
