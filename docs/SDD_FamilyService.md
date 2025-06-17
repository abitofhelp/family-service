# Software Design Document (SDD)

## Family Service GraphQL

### 1. Introduction

#### 1.1 Purpose
This document describes the software architecture and design for the Family Service GraphQL application, providing a comprehensive guide for developers to understand the system's structure and implementation details.

#### 1.2 Scope
This design document covers the architectural patterns, component interactions, data models, and technical decisions for the Family Service GraphQL application.

#### 1.3 References
- Software Requirements Specification (SRS)
- UML Class Diagram
- UML Sequence Diagram

### 2. Architectural Overview

#### 2.1 Architectural Approach
The Family Service GraphQL application follows a combination of three architectural patterns:

1. **Domain-Driven Design (DDD)**: Focuses on modeling the domain and business logic accurately
2. **Clean Architecture**: Organizes code in concentric layers with dependencies pointing inward
3. **Hexagonal Architecture (Ports and Adapters)**: Isolates the core application from external concerns

#### 2.2 High-Level Architecture Diagram

    ┌─────────────────────────────────────────────────────────────┐
    │                                                             │
    │  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐   │
    │  │             │     │             │     │             │   │
    │  │  GraphQL    │     │  Services   │     │  Domain     │   │
    │  │  API        │────▶│  Layer      │────▶│  Layer      │   │
    │  │  (Adapters) │     │             │     │             │   │
    │  │             │     │             │     │             │   │
    │  └─────────────┘     └─────────────┘     └─────────────┘   │
    │         │                   │                   │           │
    │         │                   │                   │           │
    │         ▼                   ▼                   ▼           │
    │  ┌─────────────────────────────────────────────────────┐   │
    │  │                                                     │   │
    │  │                    Ports                            │   │
    │  │                                                     │   │
    │  └─────────────────────────────────────────────────────┘   │
    │         │                                   │               │
    │         │                                   │               │
    │         ▼                                   ▼               │
    │  ┌─────────────┐                    ┌─────────────┐        │
    │  │             │                    │             │        │
    │  │  MongoDB    │                    │  PostgreSQL │        │
    │  │  Adapter    │                    │  Adapter    │        │
    │  │             │                    │             │        │
    │  └─────────────┘                    └─────────────┘        │
    │                                                             │
    └─────────────────────────────────────────────────────────────┘

#### 2.3 Design Principles
- **Separation of Concerns**: Each component has a single responsibility
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Interface Segregation**: Clients depend only on the interfaces they use
- **Dependency Injection**: Dependencies are provided to components rather than created internally
- **Immutability**: Domain objects are immutable where possible
- **Validation at Boundaries**: Input validation occurs at system boundaries

### 3. Component Design

#### 3.1 Core Domain Layer

##### 3.1.1 Aggregates
The domain layer is organized around aggregates, which are clusters of domain objects treated as a single unit:

###### Family Aggregate
- **Root Entity**: `Family`
- **Entities**: `Parent`, `Child`
- **Value Objects**: Dates, Names, IDs
- **Invariants**:
  - A family must have at least one parent
  - A family cannot have more than two parents
  - No duplicate parents in a family
  - Family status must be consistent with parent count

Example Family struct:

    // Family is the root aggregate
    type Family struct {
        id       string
        status   Status
        parents  []*parent.Parent
        children []*child.Child
    }

##### 3.1.2 Domain Services
Domain services implement business logic that doesn't naturally fit within a single entity:

Example Divorce method:

    // Divorce handles the divorce process, creating a new family for the custodial parent
    func (f *Family) Divorce(custodialParentID string) (*Family, error) {
        // Business logic for divorce process
    }

##### 3.1.3 Value Objects
Value objects are immutable and identified by their attributes rather than an identity:

Example Date value object:

    // Dates are treated as value objects
    type Date time.Time

#### 3.2 Application Services Layer

##### 3.2.1 Application Services
The application services layer includes generic and specific service implementations:

Example BaseApplicationService:

    // BaseApplicationService is a generic implementation of the ApplicationService interface
    type BaseApplicationService[T any, D any] struct {
        // Common dependencies and methods for all application services
    }

Example FamilyApplicationService:

    // FamilyApplicationService implements the application service for family-related use cases
    type FamilyApplicationService struct {
        BaseApplicationService[*entity.Family, *entity.FamilyDTO]
        familyService *domainservices.FamilyDomainService
        familyRepo    domainports.FamilyRepository
    }

Key responsibilities:
- Orchestrating domain operations
- Transaction management
- Input validation
- Error handling and mapping

##### 3.2.2 Data Transfer Objects (DTOs)
DTOs facilitate data exchange between layers:

Example FamilyDTO:

    // FamilyDTO is a data transfer object for the Family aggregate
    type FamilyDTO struct {
        ID       string
        Status   string
        Parents  []parent.ParentDTO
        Children []child.ChildDTO
    }

#### 3.3 Ports Layer

##### 3.3.1 Repository Interfaces
The ports layer defines interfaces for external dependencies:

Example generic Repository interface:

    // Repository is a generic repository interface for entity persistence operations
    type Repository[T any] interface {
        // GetByID retrieves an entity by its ID
        GetByID(ctx context.Context, id string) (T, error)

        // GetAll retrieves all entities
        GetAll(ctx context.Context) ([]T, error)

        // Save persists an entity
        Save(ctx context.Context, entity T) error
    }

Example FamilyRepository interface that embeds the generic Repository:

    // FamilyRepository defines the interface for family persistence operations
    type FamilyRepository interface {
        // Embed the generic Repository interface with Family entity
        Repository[*entity.Family]

        // FindByParentID finds families that contain a specific parent
        FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error)

        // FindByChildID finds the family that contains a specific child
        FindByChildID(ctx context.Context, childID string) (*entity.Family, error)
    }

##### 3.3.2 Application Service Interfaces
The ports layer also defines interfaces for application services:

Example generic ApplicationService interface:

    // ApplicationService is a generic interface for application services
    type ApplicationService[T any, D any] interface {
        // Create creates a new entity
        Create(ctx context.Context, dto D) (D, error)

        // GetByID retrieves an entity by ID
        GetByID(ctx context.Context, id string) (D, error)

        // GetAll retrieves all entities
        GetAll(ctx context.Context) ([]D, error)
    }

Example FamilyApplicationService interface that embeds the generic ApplicationService:

    // FamilyApplicationService defines the interface for family application services
    type FamilyApplicationService interface {
        // Embed the generic ApplicationService interface with Family entity and DTO
        ApplicationService[*entity.Family, *entity.FamilyDTO]

        // AddParent adds a parent to a family
        AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error)

        // AddChild adds a child to a family
        AddChild(ctx context.Context, familyID string, childDTO entity.ChildDTO) (*entity.FamilyDTO, error)

        // Other family-specific methods...
    }

#### 3.4 Adapters Layer

##### 3.4.1 GraphQL Adapter
The GraphQL adapter provides the API interface:

Example Resolver:

    // Resolver handles GraphQL queries and mutations
    type Resolver struct {
        FamilySvc *services.FamilyService
    }

##### 3.4.2 MongoDB Adapter
The MongoDB adapter implements the repository interface for MongoDB:

Example MongoFamilyRepository:

    // MongoFamilyRepository implements the FamilyRepository interface for MongoDB
    type MongoFamilyRepository struct {
        Collection *mongo.Collection
    }

##### 3.4.3 PostgreSQL Adapter
The PostgreSQL adapter implements the repository interface for PostgreSQL:

Example PostgresFamilyRepository:

    // PostgresFamilyRepository implements the FamilyRepository interface for PostgreSQL
    type PostgresFamilyRepository struct {
        DB *pgxpool.Pool
    }

#### 3.5 Infrastructure Components

##### 3.5.1 Error Handling
Custom error types for different layers, using generics for type-safe error categories:

Example generic AppError:

    // AppError is a generic error type that can be used for different error categories
    type AppError[T ~string] struct {
        Err     error
        Message string
        Code    string
        Type    T
    }

Example DomainError using the generic AppError:

    // DomainErrorType represents the type of domain error
    type DomainErrorType string

    // Domain error type constants
    const (
        DomainErrorGeneral DomainErrorType = "DOMAIN_ERROR"
    )

    // DomainError represents an error that occurred in the domain layer
    type DomainError = AppError[DomainErrorType]

##### 3.5.2 Validation
Validation utilities:

Example ValidationResult:

    // ValidationResult holds the result of a validation operation
    type ValidationResult struct {
        errors *errors.ValidationErrors
    }

##### 3.5.3 Middleware
Middleware for request handling:

Example RequestContext middleware:

    // WithRequestContext adds request context information to the request
    func WithRequestContext(next http.Handler) http.Handler {
        // Implementation
    }

### 4. Data Design

#### 4.1 Data Models

##### 4.1.1 MongoDB Data Model
MongoDB uses an embedded document model:

    Family Document {
        _id: string,
        status: string,
        parents: [
            {
                id: string,
                firstName: string,
                lastName: string,
                birthDate: string,
                deathDate: string (optional)
            }
        ],
        children: [
            {
                id: string,
                firstName: string,
                lastName: string,
                birthDate: string,
                deathDate: string (optional)
            }
        ]
    }

##### 4.1.2 PostgreSQL Data Model
PostgreSQL uses a normalized model with JSON for parent and child data:

    CREATE TABLE families (
        id TEXT PRIMARY KEY,
        status TEXT NOT NULL,
        parents JSONB NOT NULL,
        children JSONB NOT NULL
    );

#### 4.2 Data Flow

##### 4.2.1 Create Family Sequence
1. GraphQL resolver receives createFamily mutation
2. Input is converted to domain DTO
3. FamilyService creates Family aggregate
4. Repository saves Family to database
5. Result is converted back to GraphQL type and returned

##### 4.2.2 Divorce Sequence
1. GraphQL resolver receives divorce mutation
2. FamilyService retrieves Family from repository
3. Family.Divorce() creates new Family for custodial parent
4. Repository saves both families (original and new)
5. New Family is returned to client

### 5. Interface Design

#### 5.1 GraphQL Schema
The GraphQL schema defines the API contract:

    type Family {
      id: ID!
      status: FamilyStatus!
      parents: [Parent!]!
      children: [Child!]!
    }

    type Mutation {
      createFamily(input: FamilyInput!): Family!
      addParent(familyId: ID!, input: ParentInput!): Family!
      addChild(familyId: ID!, input: ChildInput!): Family!
      removeChild(familyId: ID!, childId: ID!): Family!
      markParentDeceased(familyId: ID!, parentId: ID!, deathDate: String!): Family!
      divorce(familyId: ID!, custodialParentId: ID!): Family!
    }

    type Query {
      getFamily(id: ID!): Family
      findFamiliesByParent(parentId: ID!): [Family!]
      findFamilyByChild(childId: ID!): Family
    }

### 6. Error Handling Design

#### 6.1 Error Types
- **ValidationError**: For input validation failures
- **DomainError**: For business rule violations
- **NotFoundError**: For resource not found situations
- **RepositoryError**: For data access issues
- **ApplicationError**: For general application errors

#### 6.2 Error Propagation
Errors are propagated up the call stack and transformed as needed:
1. Domain layer generates domain-specific errors
2. Application services wrap and enrich errors with context
3. GraphQL resolver maps errors to GraphQL-friendly format

### 7. Performance Considerations

#### 7.1 Database Optimization
- MongoDB uses embedded documents for efficient retrieval
- PostgreSQL uses JSONB for flexible querying with indexes
- Both implementations support efficient lookups by ID

#### 7.2 Caching Strategy
- No caching implemented in the current version
- Future versions could add caching at the repository or service layer

### 8. Security Considerations

#### 8.1 Input Validation
All inputs are validated at multiple levels:
- GraphQL schema validation
- Application service validation
- Domain entity validation

#### 8.2 Error Information Exposure
Error messages are sanitized before being returned to clients to prevent information leakage.

### 9. Monitoring and Observability Design

#### 9.1 Telemetry Architecture

The Family Service implements a comprehensive monitoring and observability solution using OpenTelemetry for instrumentation and Prometheus/Grafana for metrics collection and visualization.

##### 9.1.1 Components

- **OpenTelemetry SDK**: Provides the foundation for metrics collection and distributed tracing
- **Prometheus**: Time-series database for storing metrics
- **Grafana**: Visualization platform for metrics dashboards

##### 9.1.2 Architecture Diagram

```
┌─────────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                     │     │                 │     │                 │
│  Family Service     │     │  Prometheus     │     │  Grafana        │
│  ┌───────────────┐  │     │                 │     │                 │
│  │ OpenTelemetry │  │     │  ┌───────────┐  │     │  ┌───────────┐  │
│  │ SDK           │──┼────▶│  │ Metrics   │──┼────▶│  │ Dashboard │  │
│  └───────────────┘  │     │  │ Database  │  │     │  │           │  │
│                     │     │  └───────────┘  │     │  └───────────┘  │
│  ┌───────────────┐  │     │                 │     │                 │
│  │ /metrics      │◀─┼─────│  Scraper        │     │                 │
│  │ Endpoint      │  │     │                 │     │                 │
│  └───────────────┘  │     │                 │     │                 │
└─────────────────────┘     └─────────────────┘     └─────────────────┘
```

#### 9.2 Instrumentation Design

##### 9.2.1 Metrics Collection

The application uses the OpenTelemetry SDK to collect various metrics:

- **Runtime Metrics**: Memory usage, goroutines, GC statistics
- **HTTP Metrics**: Request counts, durations, error rates
- **Database Metrics**: Query counts, durations, connection pool stats
- **Application Metrics**: Business-specific metrics and error counts

Example metrics implementation:

    // Initialize metrics
    httpRequestsTotal, _ := meter.Int64Counter(
        "http_requests_total",
        metric.WithDescription("Total number of HTTP requests"),
    )

    // Record metrics
    httpRequestsTotal.Add(ctx, 1, 
        attribute.String("method", "GET"),
        attribute.String("path", "/api/families"),
        attribute.Int("status", 200),
    )

##### 9.2.2 Middleware Integration

HTTP middleware automatically collects metrics for all requests:

    // HTTP middleware for metrics collection
    func MetricsMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Call the next handler
            next.ServeHTTP(w, r)

            // Record metrics after the request is processed
            duration := time.Since(start)
            RecordHTTPRequest(r.Context(), r.Method, r.URL.Path, w.Status(), duration)
        })
    }

##### 9.2.3 Database Instrumentation

Database operations are instrumented to track performance:

    // Record database operation metrics
    func RecordDBOperation(ctx context.Context, operation, database string, duration time.Duration, err error) {
        dbOperationsTotal.Add(ctx, 1,
            attribute.String("operation", operation),
            attribute.String("database", database),
            attribute.Bool("success", err == nil),
        )

        dbOperationDuration.Record(ctx, duration.Seconds(),
            attribute.String("operation", operation),
            attribute.String("database", database),
        )
    }

#### 9.3 Metrics Exposure

The application exposes metrics through a `/metrics` endpoint in Prometheus format. This endpoint is automatically scraped by Prometheus at regular intervals.

#### 9.4 Visualization

A custom Grafana dashboard (`grafana_dashboard_for_family_service.json`) provides visualization of key metrics:

- Heap allocations over time
- HTTP request rates and durations
- Database operation performance
- Error rates

The dashboard is designed to provide insights into the application's performance and health, allowing for proactive monitoring and troubleshooting.

#### 9.5 Alerting

The monitoring system supports alerting based on metric thresholds:

- High error rates
- Elevated response times
- Memory usage spikes
- Database connection pool exhaustion

Alerts can be configured in Grafana to notify operators via email, Slack, or other channels when predefined conditions are met.

### 10. Deployment View

#### 10.1 Components
- Family Service GraphQL application
- MongoDB database
- PostgreSQL database

#### 10.2 Deployment Architecture
The application is deployed as Docker containers using docker-compose:
- Application container
- MongoDB container
- PostgreSQL container

#### 10.3 Secrets Management
The application requires a `secrets` folder containing credential files for various services:
- Grafana admin credentials
- PostgreSQL credentials
- MongoDB credentials
- Redis credentials

These secrets are mounted into the application container at runtime and are used to authenticate with the respective services. For details on setting up the secrets folder, see the [Secrets Setup Guide](Secrets_Setup_Guide.md).

### 11. Development Considerations

#### 11.1 Build and Test Process
- Go modules for dependency management
- Makefile for common development tasks
- Unit tests for domain logic
- Integration tests for repositories
- End-to-end tests for GraphQL API

#### 11.2 Development Environment Setup
- Docker and docker-compose for local development
- Environment variables for configuration
- GraphQL Playground for API exploration

### 12. Appendices

#### 12.1 UML Class Diagram
![Class Diagram](diagrams/SDD%20Class%20Diagram.svg)

This class diagram illustrates the detailed structure of the system, showing the classes in each layer (Domain, Service, Ports, Adapters) and the relationships between them. It demonstrates how the system follows the architectural patterns described in this document.

#### 12.2 UML Sequence Diagram
![Sequence Diagram](diagrams/SDD%20Sequence%20Diagram%20-%20Divorce%20Operation.svg)

This sequence diagram shows the interactions between components during the Divorce operation, one of the more complex workflows in the system. It illustrates how the different layers work together to process this operation, from the API request to the database updates.
