# Domain Services

## Overview

The Domain Services package contains services that coordinate operations on domain entities and aggregates. These services encapsulate complex business logic that doesn't naturally fit within a single entity or aggregate. The FamilyDomainService, in particular, orchestrates operations on the Family aggregate, ensuring that business rules are properly applied and that the aggregate remains in a consistent state.

## Architecture

The Domain Services package is part of the core domain layer in the Clean Architecture and Hexagonal Architecture patterns. It sits at the center of the application and has no dependencies on infrastructure or interface layers. The architecture follows these principles:

- **Domain-Driven Design (DDD)**: Services encapsulate domain logic that doesn't naturally fit within entities
- **Clean Architecture**: Services are independent of infrastructure and interface concerns
- **Hexagonal Architecture**: Services use ports (interfaces) to interact with infrastructure
- **Dependency Inversion**: Services depend on abstractions (interfaces) rather than concrete implementations

The package is organized into:

- **Family Domain Service**: Coordinates operations on the Family aggregate
- **Service Factory**: Creates and configures domain services with their dependencies
- **Service Interfaces**: Defines interfaces for domain services (in the domain ports package)

## Implementation Details

The Domain Services package implements the following design patterns:

1. **Service Pattern**: Encapsulates domain logic that doesn't naturally fit within entities
2. **Factory Pattern**: Factory methods create and configure domain services
3. **Dependency Injection**: Services receive their dependencies through constructor injection
4. **Command Pattern**: Each method represents a command that executes a specific domain operation
5. **Observer Pattern**: Services notify observers (metrics, logs) of important domain events

Key implementation details:

- **Repository Dependency**: Services depend on repository interfaces defined in the domain ports package
- **Context Propagation**: All methods accept a context.Context parameter for cancellation and value propagation
- **Error Handling**: Services use domain-specific error types for clear error handling
- **Logging**: All operations are logged with appropriate context for observability
- **Tracing**: Operations are traced using OpenTelemetry for distributed tracing
- **Metrics**: Important domain events are recorded as metrics for monitoring
- **Validation**: Services validate inputs and enforce business rules

## Features

- **Family Domain Service**: Coordinates operations on the Family aggregate
- **Transaction Management**: Ensures data consistency across complex operations
- **Business Rule Enforcement**: Applies domain-specific business rules
- **Distributed Tracing**: Integrates with OpenTelemetry for distributed tracing
- **Comprehensive Logging**: Provides detailed logging of all operations
- **Metrics Collection**: Collects metrics for monitoring application behavior
- **Error Handling**: Properly handles and propagates domain-specific errors

## API Documentation

### Core Types

#### FamilyDomainService

The FamilyDomainService coordinates operations on the Family aggregate, ensuring that business rules are properly applied.

```
// FamilyDomainService is a domain service that coordinates operations on the Family aggregate
type FamilyDomainService struct {
    repo   ports.FamilyRepository
    logger *loggingwrapper.ContextLogger
    tracer trace.Tracer
}
```

### Key Methods

#### CreateFamily

Creates a new family with validation.

```
// CreateFamily creates a new family
func (s *FamilyDomainService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error)
```

#### AddParent

Adds a parent to a family, updating the family status if necessary.

```
// AddParent adds a parent to a family
func (s *FamilyDomainService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error)
```

#### Divorce

Handles the divorce process, creating a new family for the non-custodial parent.

```
// Divorce handles the divorce process
func (s *FamilyDomainService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error)
```

## Examples

There may be additional examples in the /EXAMPLES directory.

Example of using the FamilyDomainService:

```
// Create a new domain service
repo := postgres.NewFamilyRepository(db)
logger := loggingwrapper.NewContextLogger(zapLogger)
tracer := telemetrywrapper.NewTracer("family-service")
service := services.NewFamilyDomainService(repo, logger, tracer)

// Create a new family
familyDTO := entity.FamilyDTO{
    Status: entity.StatusSingle,
    Parents: []*entity.ParentDTO{
        {
            FirstName: "John",
            LastName: "Doe",
            BirthDate: birthDate,
        },
    },
}
createdFamily, err := service.CreateFamily(ctx, familyDTO)
if err != nil {
    // Handle error
}

// Add a parent to the family
parentDTO := entity.ParentDTO{
    FirstName: "Jane",
    LastName: "Doe",
    BirthDate: birthDate,
}
updatedFamily, err := service.AddParent(ctx, createdFamily.ID, parentDTO)
if err != nil {
    // Handle error
}

// Divorce the family
newFamily, err := service.Divorce(ctx, updatedFamily.ID, updatedFamily.Parents[0].ID)
if err != nil {
    // Handle error
}
```

## Configuration

The Domain Services package can be configured with the following options:

- **Logger Configuration**: Configure the logger used by the services
- **Tracer Configuration**: Configure the tracer used for distributed tracing
- **Metrics Configuration**: Configure the metrics collected by the services
- **Validation Configuration**: Configure the validation rules applied by the services

Example configuration:

```
// Configure the logger
logger := loggingwrapper.NewContextLogger(zapLogger)

// Configure the tracer
tracer := telemetrywrapper.NewTracer("family-service")

// Create the domain service with the configuration
service := services.NewFamilyDomainService(
    repo,
    logger,
    tracer,
)
```

## Testing

The Domain Services package is tested through:

1. **Unit Tests**: Each service method has comprehensive unit tests
2. **Integration Tests**: Tests that verify the services work correctly with real repositories
3. **Mocking**: Tests use mocks for repositories to isolate the service being tested

Key testing approaches:

- **Repository Mocking**: Tests use mock repositories to isolate the service being tested
- **Business Rule Testing**: Tests verify that business rules are properly applied
- **Error Handling**: Tests verify that errors are properly handled and propagated
- **Edge Case Testing**: Tests verify that edge cases are handled correctly
- **Scenario Testing**: Tests verify that complex scenarios work correctly

Example of a test case:

```
func TestFamilyDomainService_CreateFamily(t *testing.T) {
    // Create mocks
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mock.NewMockFamilyRepository(ctrl)

    // Set up expectations
    mockRepo.EXPECT().
        Save(gomock.Any(), gomock.Any()).
        Return(nil)

    // Create the service
    service := services.NewFamilyDomainService(
        mockRepo,
        loggingwrapper.NewContextLogger(zapLogger),
        telemetrywrapper.NewTracer("test"),
    )

    // Call the method
    result, err := service.CreateFamily(context.Background(), entity.FamilyDTO{
        Status: entity.StatusSingle,
        Parents: []*entity.ParentDTO{
            {
                FirstName: "John",
                LastName: "Doe",
                BirthDate: time.Now(),
            },
        },
    })

    // Verify the result
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, entity.StatusSingle, result.Status)
    assert.Len(t, result.Parents, 1)
}
```

## Design Notes

1. **Domain Logic Encapsulation**: Domain services encapsulate complex domain logic that doesn't naturally fit within entities
2. **Statelessness**: Domain services are stateless, with all state maintained in the entities
3. **Repository Dependency**: Domain services depend on repository interfaces for data access
4. **Validation**: Domain services validate inputs and enforce business rules
5. **Error Handling**: Domain services use domain-specific error types for clear error handling
6. **Logging and Tracing**: Domain services log and trace all operations for observability
7. **Metrics**: Domain services collect metrics for monitoring application behavior

## Best Practices

1. **Single Responsibility**: Each domain service should focus on a specific domain concept
2. **Statelessness**: Domain services should be stateless, with all state maintained in the entities
3. **Dependency Injection**: Use dependency injection to provide repositories and other dependencies
4. **Comprehensive Logging**: Log all operations with appropriate context for observability
5. **Proper Error Handling**: Handle errors appropriately and provide meaningful error messages

## References

- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Service Pattern](https://martinfowler.com/articles/injection.html)
- [Domain Events](https://martinfowler.com/eaaDev/DomainEvent.html)
- [Domain Entities](../entity/README.md) - The entities managed by these services
- [Domain Ports](../ports/README.md) - The interfaces used by these services for data access
- [Domain Errors](../errors/README.md) - Domain-specific error types used by these services
- [Application Services](../../application/services/README.md) - Higher-level services that use these domain services
