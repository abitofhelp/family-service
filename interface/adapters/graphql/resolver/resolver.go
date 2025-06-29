// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"github.com/abitofhelp/family-service/core/application/ports"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/dto"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/generated"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver serves as a dependency injection container for your resolvers.
type Resolver struct {
	familyService ports.FamilyApplicationService
	mapper       dto.FamilyMapper
}

// NewResolver creates a new resolver with the given dependencies.
func NewResolver(familyService ports.FamilyApplicationService, mapper dto.FamilyMapper) *Resolver {
	return &Resolver{
		familyService: familyService,
		mapper:       mapper,
	}
}

// Query returns the query resolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// Mutation returns the mutation resolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

// Family returns the family resolver implementation.
func (r *Resolver) Family() generated.FamilyResolver {
	return &familyResolver{r}
}

type (
	queryResolver struct {
		*Resolver
	}
	mutationResolver struct {
		*Resolver
	}
	familyResolver struct {
		*Resolver
	}
)
