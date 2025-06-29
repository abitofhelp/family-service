# Application Services

## Overview

The Application Services package implements the application layer in the Clean Architecture pattern. It contains services that orchestrate the execution of use cases by coordinating domain services, repositories, and other infrastructure components. These services act as a bridge between the interface layer and the domain layer, ensuring that business rules are properly applied while handling cross-cutting concerns like logging, caching, and error handling.

## Architecture

The Application Services package is part of the application layer in the Clean Architecture and Hexagonal Architecture patterns. It sits between the interface layer and the domain layer, coordinating the flow of data and control between them. The architecture follows these principles:

- **Clean Architecture**: Application services depend on the domain layer but are independent of the interface and infrastructure layers
- **Hexagonal Architecture**: Application services implement ports defined in the application layer and use ports defined in the domain layer
- **Dependency Inversion**: Application services depend on abstractions (interfaces) rather than concrete implementations
- **Separation of Concerns**: Each application service is responsible for a specific set of related use cases

The package is organized into:

- **Base Application Service**: A generic base class that provides common functionality for all application services
- **Family Application Service**: Implements family-specific use cases
- **Service Factory**: Creates and configures application services with their dependencies

## Implementation Details

The Application Services package implements the following design patterns:

1. **Facade Pattern**: Application services provide a simplified interface to the domain layer
2. **Command Pattern**: Each method represents a command that executes a specific use case
3. **Decorator Pattern**: Cross-cutting concerns like logging and caching are implemented as decorators
4. **Template Method Pattern**: The base application service defines a template for common operations
5. **Factory Pattern**: Service factories create and configure application services

Key implementation details:

- **Dependency Injection**: Services receive their dependencies through constructor injection
- **Context Propagation**: All methods accept a context.Context parameter for cancellation and value propagation
- **Error Handling**: Services translate domain errors to application errors when appropriate
- **Caching**: Frequently accessed data is cached to improve performance
- **Logging**: All operations are logged with appropriate context for observability
- **Transaction Management**: Services ensure data consistency across operations

## Features

- **Family Application Service**: Implements operations for managing families, parents, and children
- **Caching Integration**: Uses caching to improve performance for frequently accessed data
- **Comprehensive Logging**: Detailed logging of all operations for observability
- **Error Handling**: Proper error handling and propagation
- **Transaction Management**: Ensures data consistency across operations
- **Clean Architecture Compliance**: Follows the principles of Clean Architecture

## API Documentation

### Core Types

#### FamilyApplicationService

The FamilyApplicationService implements the application service for family-related use cases. It provides methods for creating and managing families, adding and removing parents and children, handling divorces, and finding families by parent or child.

```
// FamilyApplicationService implements the application service for family-related use cases
type FamilyApplicationService struct {
    BaseApplicationService[*entity.Family, *entity.FamilyDTO]
    familyService *domainservices.FamilyDomainService
    familyRepo    domainports.FamilyRepository
    logger        *logging.ContextLogger
    cache         *cache.Cache
}
```

### Key Methods

#### Create

Creates a new family.

```
// Create creates a new family
func (s *FamilyApplicationService) Create(ctx context.Context, dto *entity.FamilyDTO) (*entity.FamilyDTO, error)
```

#### AddParent

Adds a parent to a family.

```
// AddParent adds a parent to a family
func (s *FamilyApplicationService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error)
```

#### Divorce

Handles the divorce process, creating a new family for the custodial parent and children.

```
// Divorce handles the divorce process
func (s *FamilyApplicationService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error)
```

## Examples

There may be additional examples in the /EXAMPLES directory.

Example of using the FamilyApplicationService:

```
// Create a new family
familyDTO, err := familyService.Create(ctx, &entity.FamilyDTO{
    Status: entity.StatusSingle,
    Parents: []*entity.ParentDTO{
        {
            FirstName: "John",
            LastName: "Doe",
            BirthDate: birthDate,
        },
    },
})
if err != nil {
    // Handle error
}

// Add a child to the family
familyDTO, err = familyService.AddChild(ctx, familyDTO.ID, entity.ChildDTO{
    FirstName: "Jane",
    LastName: "Doe",
    BirthDate: childBirthDate,
})
if err != nil {
    // Handle error
}
```

## Configuration

The Application Services package can be configured with the following options:

- **Cache Configuration**: Configure the cache size, TTL, and eviction policy
- **Logging Configuration**: Configure the log level, format, and output destination
- **Transaction Configuration**: Configure transaction timeout and retry policy
- **Concurrency Configuration**: Configure the maximum number of concurrent operations

Example configuration:

```
// Configure the cache
cacheConfig := cache.Config{
    Size: 1000,
    TTL: 5 * time.Minute,
    EvictionPolicy: cache.LRU,
}

// Configure the logger
loggerConfig := logging.Config{
    Level: logging.InfoLevel,
    Format: logging.JSONFormat,
    Output: os.Stdout,
}

// Create the application service with the configuration
familyService := services.NewFamilyApplicationService(
    familyRepo,
    familyDomainService,
    cache.New(cacheConfig),
    logging.NewLogger(loggerConfig),
)
```

## Testing

The Application Services package is tested through:

1. **Unit Tests**: Each service method has comprehensive unit tests
2. **Integration Tests**: Tests that verify the services work correctly with real repositories
3. **Mocking**: Tests use mocks for dependencies to isolate the service being tested

Key testing approaches:

- **Dependency Mocking**: Tests use mocks for repositories and domain services
- **Context Propagation**: Tests verify that context is properly propagated
- **Error Handling**: Tests verify that errors are properly handled and propagated
- **Caching**: Tests verify that caching works correctly
- **Transaction Management**: Tests verify that transactions are properly managed

Example of a test case:

```
func TestFamilyApplicationService_Create(t *testing.T) {
    // Create mocks
    mockRepo := mocks.NewMockFamilyRepository(ctrl)
    mockDomainService := mocks.NewMockFamilyDomainService(ctrl)

    // Set up expectations
    mockDomainService.EXPECT().
        CreateFamily(gomock.Any(), gomock.Any()).
        Return(&entity.Family{}, nil)
    mockRepo.EXPECT().
        Save(gomock.Any(), gomock.Any()).
        Return(nil)

    // Create the service
    service := services.NewFamilyApplicationService(
        mockRepo,
        mockDomainService,
        cache.New(cache.Config{}),
        logging.NewLogger(logging.Config{}),
    )

    // Call the method
    result, err := service.Create(context.Background(), &entity.FamilyDTO{})

    // Verify the result
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## Design Notes

1. **Application vs. Domain Logic**: Application services orchestrate use cases but delegate domain logic to domain services
2. **Caching Strategy**: Caching is implemented at the application service level to improve performance
3. **Error Translation**: Domain errors are translated to application errors when appropriate
4. **Transaction Boundaries**: Application services define transaction boundaries
5. **Context Propagation**: Context is propagated through all method calls for cancellation and value propagation
6. **Dependency Injection**: Services receive their dependencies through constructor injection for testability
7. **Interface Segregation**: Services implement interfaces defined in the application ports package

## Best Practices

1. **Separation of Concerns**: Application services should focus on orchestrating use cases, delegating domain logic to domain services
2. **Comprehensive Logging**: Log all operations with appropriate context for observability
3. **Proper Error Handling**: Handle errors appropriately and provide meaningful error messages
4. **Caching Strategy**: Use caching for frequently accessed data to improve performance
5. **Transaction Management**: Ensure data consistency across operations

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [CQRS](https://martinfowler.com/bliki/CQRS.html)
- [Application Ports](../ports/README.md) - Defines the interfaces implemented by these services
- [Domain Services](../../domain/services/README.md) - Provides the domain logic used by these services
- [Domain Entities](../../domain/entity/README.md) - Defines the entity types used by these services
- [Repositories](../../domain/ports/README.md) - Provides data access for these services
