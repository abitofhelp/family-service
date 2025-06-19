// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"context"
	"fmt"
	"time"
	"github.com/99designs/gqlgen/graphql"
	"github.com/abitofhelp/family-service/core/application/ports"
	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/valueobject"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	mygraphql "github.com/abitofhelp/servicelib/graphql"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	AuthorizationCheckDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "authorization_check_duration_seconds",
		Help:    "Duration of authorization checks in seconds",
		Buckets: prometheus.DefBuckets,
	})
	AuthorizationFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "authorization_failures_total",
		Help: "Total number of authorization failures",
	})
)

func init() {
	// Register metrics - using Register instead of MustRegister to avoid panic on duplicate registration
	if err := prometheus.Register(AuthorizationCheckDuration); err != nil {
		// If the error is because the metric is already registered, we can ignore it
		// Otherwise, log the error but don't panic
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			fmt.Printf("Error registering AuthorizationCheckDuration metric: %v\n", err)
		}
	}

	if err := prometheus.Register(AuthorizationFailures); err != nil {
		// If the error is because the metric is already registered, we can ignore it
		// Otherwise, log the error but don't panic
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			fmt.Printf("Error registering AuthorizationFailures metric: %v\n", err)
		}
	}
}

// Resolver is the resolver for GraphQL queries and mutations
type Resolver struct {
	familyService ports.FamilyApplicationService
	logger        *logging.ContextLogger
	tracer        trace.Tracer
}

// NewResolver creates a new GraphQL resolver
func NewResolver(familyService ports.FamilyApplicationService, logger *logging.ContextLogger) *Resolver {
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
	// We need a context for logging, but this method doesn't receive one
	// Create a background context for logging purposes
	ctx := context.Background()

	// Convert domain status to GraphQL status
	status := model.FamilyStatus(dto.Status)

	// Convert parents from DTO to model
	parents := make([]*model.Parent, 0, len(dto.Parents))
	for _, parentDTO := range dto.Parents {
		parent, err := entity.ParentFromDTO(parentDTO)
		if err != nil {
			// Log error but continue
			r.logger.Warn(ctx, "Failed to convert parent DTO to entity in resolver", 
				zap.Error(err), 
				zap.String("parent_id", parentDTO.ID),
				zap.String("family_id", dto.ID))
			continue
		}

		// Convert entity to model
		deathDate := parent.DeathDate()
		var deathDateStr *string
		if deathDate != nil {
			formatted := parent.DeathDate().Format(time.RFC3339)
			deathDateStr = &formatted
		}

		modelParent := &model.Parent{
			ID:        valueobject.ID(parent.ID()),
			FirstName: parent.FirstName(),
			LastName:  parent.LastName(),
			BirthDate: parent.BirthDate().Format(time.RFC3339),
			DeathDate: deathDateStr,
		}

		parents = append(parents, modelParent)
	}

	// Convert children from DTO to model
	children := make([]*model.Child, 0, len(dto.Children))
	for _, childDTO := range dto.Children {
		child, err := entity.ChildFromDTO(childDTO)
		if err != nil {
			// Log error but continue
			r.logger.Warn(ctx, "Failed to convert child DTO to entity in resolver", 
				zap.Error(err), 
				zap.String("child_id", childDTO.ID),
				zap.String("family_id", dto.ID))
			continue
		}

		// Convert entity to model
		deathDate := child.DeathDate()
		var deathDateStr *string
		if deathDate != nil {
			formatted := child.DeathDate().Format(time.RFC3339)
			deathDateStr = &formatted
		}

		modelChild := &model.Child{
			ID:        valueobject.ID(child.ID()),
			FirstName: child.FirstName(),
			LastName:  child.LastName(),
			BirthDate: child.BirthDate().Format(time.RFC3339),
			DeathDate: deathDateStr,
		}

		children = append(children, modelChild)
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
	return mygraphql.HandleError(ctx, err, operation, r.logger)
}

// CheckAuthorization checks if the user is authorized to perform the specified operation
// It returns an error if the user is not authorized
func (r *Resolver) CheckAuthorization(ctx context.Context, allowedRoles []string, requiredScopes []string, resource string, operation string) error {
	// Use the generic CheckAuthorization function from servicelib
	return mygraphql.CheckAuthorization(ctx, allowedRoles, requiredScopes, resource, operation, r.logger)
}

func (r *Resolver) IsAuthorized(ctx context.Context, obj any, next graphql.Resolver, allowedRoles []model.Role, requiredScopes []model.Scope, resource *model.Resource) (res any, err error) {
	// Validate roles
	for _, role := range allowedRoles {
		if !role.IsValid() {
			r.logger.Warn(ctx, "Invalid role specified in authorization check", zap.String("role", role.String()))
			return nil, fmt.Errorf("invalid role specified: %s", role)
		}
	}

	// Validate scopes
	for _, scope := range requiredScopes {
		if !scope.IsValid() {
			r.logger.Warn(ctx, "Invalid scope specified in authorization check", zap.String("scope", scope.String()))
			return nil, fmt.Errorf("invalid scope specified: %s", scope)
		}
	}

	// Validate resource
	if resource != nil && !resource.IsValid() {
		r.logger.Warn(ctx, "Invalid resource specified in authorization check", zap.String("resource", resource.String()))
		return nil, fmt.Errorf("invalid resource specified: %s", resource)
	}

	// Convert model.Role to a string array for middleware check
	strRoles := mygraphql.ConvertRolesToStrings(allowedRoles)

	// Convert model.Scope to a string array for middleware check
	strScopes := mygraphql.ConvertRolesToStrings(requiredScopes)

	// Convert model.Resource to a string for middleware check
	strResource := ""
	if resource != nil {
		strResource = resource.String()
	}

	// Use the generic directive implementation from servicelib
	return mygraphql.IsAuthorizedDirective(ctx, obj, next, strRoles, strScopes, strResource, r.logger)
}
