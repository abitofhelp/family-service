package resolver

import (
	"context"
	"github.com/abitofhelp/family-service/core/application/ports"
	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/valueobject"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/abitofhelp/servicelib/graphql"
	"github.com/abitofhelp/servicelib/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Resolver is the resolver for GraphQL queries and mutations
type Resolver struct {
	familyService ports.FamilyApplicationService
	//FIXME: authService   ports.AuthorizationService
	logger *logging.ContextLogger
	tracer trace.Tracer
}

// NewResolver creates a new GraphQL resolver
func NewResolver(familyService ports.FamilyApplicationService, logger *logging.ContextLogger) *Resolver {
	// FIXME: func NewResolver(familyService ports.FamilyApplicationService, authService ports.AuthorizationService, logger *zap.Logger) *Resolver {
	if familyService == nil {
		panic("family application service cannot be nil")
	}
	if logger == nil {
		panic("logger cannot be nil")
	}
	return &Resolver{
		familyService: familyService,
		// FIXME: authService:   authService,
		logger: logger,
		tracer: otel.Tracer("graphql.resolver"),
	}
}

// dtoToModelFamily converts a domain FamilyDTO to a GraphQL Family model
func (r *Resolver) dtoToModelFamily(dto entity.FamilyDTO) *model.Family {
	// Convert domain status to GraphQL status
	status := model.FamilyStatus(dto.Status)

	// Convert parents from DTO to entity
	parents := make([]*entity.Parent, 0, len(dto.Parents))
	for _, parentDTO := range dto.Parents {
		parent, err := entity.ParentFromDTO(parentDTO)
		if err != nil {
			// Log error but continue
			continue
		}
		parents = append(parents, parent)
	}

	// Convert children from DTO to entity
	children := make([]*entity.Child, 0, len(dto.Children))
	for _, childDTO := range dto.Children {
		child, err := entity.ChildFromDTO(childDTO)
		if err != nil {
			// Log error but continue
			continue
		}
		children = append(children, child)
	}

	// Create and return the GraphQL model
	return &model.Family{
		ID:       valueobject.ID(dto.ID),
		Status:   status,
		Parents:  parents,
		Children: children,
	}
}

// GetParentCount returns the number of parents in a family
func (r *Resolver) GetParentCount(ctx context.Context, obj *model.Family) (int, error) {
	return len(obj.Parents), nil
}

// GetChildrenCount returns the number of children in a family
func (r *Resolver) GetChildrenCount(ctx context.Context, obj *model.Family) (int, error) {
	return len(obj.Children), nil
}

// HandleError processes an error and returns an appropriate GraphQL error
// It logs the error and converts it to a GraphQL error with appropriate extensions
func (r *Resolver) HandleError(ctx context.Context, err error, operation string) error {
	// Use the pkg/graphql.HandleError function
	return graphql.HandleError(ctx, err, operation, r.logger)
}
