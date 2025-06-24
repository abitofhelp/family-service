# Family Service GraphQL
[![codecov](https://codecov.io/gh/abitofhelp/family-service/graph/badge.svg?token=KQI02N7OS7)](https://codecov.io/gh/abitofhelp/family-service)

## 📖 Overview

Family-Service is a backend service designed to track dynamic family relationships using modern software engineering principles. It supports use cases such as marriage, divorce, remarriage, single-parent families, and parental death. Data integrity and family structure constraints are enforced through domain validation logic. The architecture is based on Domain-Driven Design (DDD), Clean Architecture, and Hexagonal Architecture, providing a scalable and testable foundation for managing families across SQLite (local development), MongoDB (document model), and PostgreSQL (relational model).

Key features include:

- Rich relationship modeling between families, parents, and children
- Transaction-safe CRUD operations for SQLite, MongoDB, and PostgreSQL
- Validation of unique parent identity within a family (name + birthdate)
- Cascade logic during divorce/remarriage events
- Logging with Zap, OS signal handling, and context propagation
  A backend service for managing complex family relationships (parents, children, families) using Domain-Driven Design, Clean Architecture, and Hexagonal Architecture.

## 🏗 Architecture

The system uses a strict separation of concerns:

- `core/domain`: Aggregates, value objects, and validation logic
- `core/services`: Domain services encapsulating business rules
- `core/ports`: Interfaces for repositories and services
- `infrastructure`: Layer containing adapters and infrastructure concerns for components DRIVEN by the application
  - `adapters`: Database implementations (MongoDB, PostgreSQL)
  - `pkg`: Cross-cutting concerns (errors, middleware, validation)
- `interface`: Layer containing adapters for components that DRIVE the application
  - `adapters/graphql`: API layer, built with gqlgen, acting as the application service layer

Design is driven by:

- **DDD**: Models families, parents, and children as aggregates with invariants
- **Clean Architecture**: Keeps business logic decoupled from frameworks and databases
- **Hexagonal Architecture**: Uses ports and adapters for flexibility and testability
- **Generics**: Leverages Go's generics for type-safe, reusable code patterns

MongoDB embeds family structure in documents for performance; PostgreSQL uses foreign keys and joins with transactions for consistency.

- Domain logic is isolated in `core/domain` and `core/services`
- Infrastructure concerns in `infrastructure/`:
  - Database adapters in `infrastructure/adapters/` (SQLite, MongoDB, PostgreSQL)
  - Cross-cutting concerns in `infrastructure/pkg/` (errors, middleware, validation)
- GraphQL transport in `interface/adapters/graphql/`
- Ports/interfaces in `core/ports/`
- Fully supports SQLite (local development), MongoDB (embedded), and PostgreSQL (normalized)
- Clean layering ensures separation of concerns and easy testability

## 🧩 Generics Implementation

Go 1.18 introduced generics, allowing for type-parameterized functions and types. This project leverages generics in several key areas to reduce code duplication and improve type safety:

### Repository Interfaces

A generic `Repository` interface provides common CRUD operations for any entity type:

```
// Repository is a generic repository interface for entity persistence operations
type Repository[T any] interface {
    // GetByID retrieves an entity by its ID
    GetByID(ctx context.Context, id string) (T, error)

    // GetAll retrieves all entities
    GetAll(ctx context.Context) ([]T, error)

    // Save persists an entity
    Save(ctx context.Context, entity T) error
}
```

The `FamilyRepository` interface embeds this generic interface with the specific entity type:

```
package ports

import (
    "context"
    "github.com/abitofhelp/family-service/core/domain/entity"
)

// FamilyRepository defines the interface for family persistence operations
type FamilyRepository interface {
    // Embed the generic Repository interface with Family entity
    Repository[*entity.Family]

    // FindByParentID finds families that contain a specific parent
    FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error)

    // FindByChildID finds the family that contains a specific child
    FindByChildID(ctx context.Context, childID string) (*entity.Family, error)
}
```

### Application Service Interfaces

A generic `ApplicationService` interface provides common operations for any entity type:

```
package ports

import (
    "context"
)

// ApplicationService is a generic interface for application services
type ApplicationService[T any, D any] interface {
    // Create creates a new entity
    Create(ctx context.Context, dto D) (D, error)

    // GetByID retrieves an entity by ID
    GetByID(ctx context.Context, id string) (D, error)

    // GetAll retrieves all entities
    GetAll(ctx context.Context) ([]D, error)
}
```

### Error Handling Utilities

A generic `AppError` type provides a common error structure for different error categories:

```
package errors

// AppError is a generic error type that can be used for different error categories
type AppError[T ~string] struct {
    Err     error
    Message string
    Code    string
    Type    T
}
```

### Benefits of Generics

- **Code Reuse**: Common patterns like CRUD operations are defined once and reused across different entity types
- **Type Safety**: The compiler ensures that the correct types are used with the generic interfaces
- **Maintainability**: Changes to common patterns only need to be made in one place
- **Extensibility**: New entity types can easily be added by implementing the generic interfaces

### Implementation Details

#### Repository Interface

The generic repository interface was implemented in `core/domain/ports/repository.go`. This interface provides a common set of CRUD operations that can be used with any entity type. The `FamilyRepository` interface embeds this generic interface with the specific `Family` entity type.

#### Application Service Interface

The generic application service interface was implemented in `core/application/ports/application.go`. This interface provides a common set of operations that can be used with any entity and DTO type. The `FamilyApplicationService` interface embeds this generic interface with the specific `Family` entity and `FamilyDTO` types.

#### Error Handling Utilities

The generic error type was implemented in `pkg/errors/errors.go`. This type provides a common structure for different error categories, allowing for type-safe error handling. The `DomainError`, `RepositoryError`, and `ApplicationError` types are all based on this generic error type.

### Future Improvements

Future improvements could include:

1. Implementing more generic repository implementations for different storage backends.
2. Creating generic validation utilities for common validation patterns.
3. Implementing generic query builders for more complex database operations.

## 🔌 ServiceLib Integration

This project leverages several packages from the `servicelib` library to enhance reliability, performance, and observability:

### Telemetry Integration

The `servicelib/telemetry` package is used for distributed tracing and metrics:

- **Distributed Tracing**: OpenTelemetry tracing is integrated throughout the application to provide end-to-end visibility into request flows. Traces are collected and can be exported to various backends.
- **Metrics**: Prometheus metrics are exposed at the `/metrics` endpoint, providing insights into application performance and behavior.

```
// Example of tracing integration in domain services
func (s *FamilyDomainService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
    // Start a new span for this operation
    ctx, span := s.tracer.Start(ctx, "FamilyDomainService.CreateFamily")
    defer span.End()

    // Operation logic...
}
```

### Cache Integration

The `servicelib/cache` package is used for caching frequently accessed data:

- **In-Memory Cache**: A configurable in-memory cache with TTL and automatic cleanup is implemented.
- **Cache Middleware**: Cache middleware functions make it easy to add caching to any operation.

```
// Example of cache integration in application services
func (s *FamilyApplicationService) GetByID(ctx context.Context, id string) (*entity.FamilyDTO, error) {
    // Create cache key
    cacheKey := fmt.Sprintf("family:%s", id)

    // Try to get from cache or call the domain service
    result, err := cache.WithContextCache(ctx, s.cache, cacheKey, func(ctx context.Context) (interface{}, error) {
        // Delegate to domain service
        return s.familyService.GetFamily(ctx, id)
    })

    // Process result...
}
```

### Circuit Breaker Integration

The `servicelib/circuit` package is used for circuit breaking on external dependencies:

- **Circuit Breaker Pattern**: Prevents cascading failures when external dependencies are unavailable.
- **Configurable Thresholds**: Error thresholds, volume thresholds, and sleep windows are all configurable.
- **Fallback Support**: Support for fallback functions when the circuit is open.

```
// Example of circuit breaker integration
func (repo *MongoRepository) GetByID(ctx context.Context, id string) (*entity.Family, error) {
    var family *entity.Family

    err := repo.circuitBreaker.Execute(ctx, "GetByID", func(ctx context.Context) error {
        // Database operation...
        return nil
    })

    return family, err
}
```

### Rate Limiter Integration

The `servicelib/rate` package is used for rate limiting to protect resources:

- **Token Bucket Algorithm**: Implements the token bucket algorithm for rate limiting.
- **Configurable Rates**: Requests per second and burst size are configurable.
- **Wait or Fail Options**: Support for both immediate failure and waiting for a token.

```
// Example of rate limiter integration
func (s *FamilyApplicationService) Create(ctx context.Context, dto *entity.FamilyDTO) (*entity.FamilyDTO, error) {
    var family *entity.FamilyDTO

    err := s.rateLimiter.Execute(ctx, "Create", func(ctx context.Context) error {
        // Create operation...
        return nil
    })

    return family, err
}
```

### Configuration

All servicelib integrations are configurable through the application's configuration system:

```yaml
# Example configuration
app:
  cache:
    enabled: true
    ttl: 5m
    max_size: 1000
    purge_interval: 10m

  circuit:
    enabled: true
    timeout: 5s
    max_concurrent: 100
    error_threshold: 0.5
    volume_threshold: 20
    sleep_window: 10s

  rate:
    enabled: true
    requests_per_second: 100
    burst_size: 50

  telemetry:
    tracing:
      enabled: true
      otlp:
        endpoint: localhost:4317
        insecure: true
        timeout: 5s
```

## 📊 UML Diagrams

The following UML diagrams provide visual representations of the system architecture, design, and processes:

### Domain and Architecture Diagrams
- [SRS Use Case Diagram](./docs/diagrams/SRS%20Use%20Case%20Diagram.svg) - Shows the main operations that can be performed by API clients
- [SRS Data Model Diagram](./docs/diagrams/SRS%20Data%20Model%20Diagram.svg) - Illustrates the key entities and their relationships
- [SDD Class Diagram](./docs/diagrams/SDD%20Class%20Diagram.svg) - Detailed structure of the system's classes and their relationships
- [SDD Sequence Diagram - Divorce Operation](./docs/diagrams/SDD%20Sequence%20Diagram%20-%20Divorce%20Operation.svg) - Shows the interactions during a complex workflow

### Testing Diagrams
- [STP Test Process Activity Diagram](./docs/diagrams/STP%20Test%20Process%20Activity%20Diagram.svg) - Illustrates the test execution process
- [STP Test Coverage Component Diagram](./docs/diagrams/STP%20Test%20Coverage%20Component%20Diagram.svg) - Shows test coverage targets for each component

### Deployment Diagrams
- [Deployment Diagram](./docs/diagrams/Deployment%20Diagram.svg) - Shows the physical architecture of the deployed system
- [Deployment Container Diagram](./docs/diagrams/Deployment%20Container%20Diagram.svg) - Illustrates the Docker Compose configuration

## 📋 Developer Documentation

The following documents provide comprehensive detail for developers to fully understand, implement, and maintain this system. Each includes embedded UML diagrams and references implementation decisions based on real-world constraints and design tradeoffs.

Each document is crafted to guide a developer from understanding to implementation:

- **Software Requirements Specification (SRS)**: Defines core use cases, business rules, and validation requirements for families, parents, and children. Includes lifecycle flows for marriage, divorce, and death.
- **Software Design Document (SDD)**: Describes the layered architecture, key interfaces, technology choices, and database models (embedded for MongoDB, normalized for PostgreSQL). Includes class and sequence diagrams.
- **Software Test Plan (STP)**: Covers testing strategy including unit, repository, integration, and end-to-end tests. Defines coverage targets, edge cases (e.g. invalid dates, duplicate parents), and transaction rollback scenarios.

All diagrams and specifications are embedded or referenced in the spec PDF and available individually in the `/docs` folder.

- [Software Requirements Specification (SRS)](./docs/SRS_FamilyService.md)
- [Software Design Document (SDD)](./docs/SDD_FamilyService.md)
- [Software Test Plan (STP)](./docs/STP_FamilyService.md)
- [Deployment Document](./docs/Deployment_FamilyService.md)

## ⚙️ Running the Project

### Prerequisites

- Go 1.24+
- SQLite (default for local development)
- MongoDB / PostgreSQL (optional for advanced use cases)
- Docker (optional for local dev)
- Secrets folder setup (see [Secrets Setup Guide](docs/Secrets_Setup_Guide.md))

### Database Initialization Scripts

The project includes initialization scripts for all supported databases that create the necessary database structures and populate them with sample data:

- `data/dev/mongo/mongodb_init.js`: Creates the MongoDB database, collections, indexes, and inserts sample family data
- `data/dev/postgres/postgresql_init.sql`: Creates the PostgreSQL tables, indexes, triggers, and inserts sample family data
- `data/dev/sqlite/sqlite_init.go`: Creates the SQLite database file, tables, and inserts sample family data

Note: While the SQLite adapter automatically creates the database file and tables when the application starts, the initialization script is provided for convenience to populate the database with sample data.

#### Sample Data

The scripts populate the databases with the following sample data:
- A traditional family with two parents and two children
- A family with two parents and one child
- A divorced family with one parent and no children
- A family with one living parent, one deceased parent, and two children

#### Running the Scripts Manually

**Prerequisites:**
- For MongoDB: MongoDB server installed and running (version 4.0+)
- For PostgreSQL: PostgreSQL server installed and running (version 12+)
- For SQLite: Go installed (version 1.20+)

**MongoDB Script:**
```bash
# Connect to MongoDB and run the script
mongosh < data/dev/mongo/mongodb_init.js

# Alternative method with specific host/port
mongosh --host localhost --port 27017 data/dev/mongo/mongodb_init.js

# To verify the data was inserted correctly
mongosh
> use family_service
> db.families.find().pretty()
> db.families.count()  # Should return 4
```

**PostgreSQL Script:**
```bash
# Create the database first (if it doesn't exist)
createdb family_service

# Run the initialization script
psql -d family_service -f data/dev/postgres/postgresql_init.sql

# Alternative method with username/password
psql -U username -W -d family_service -f data/dev/postgres/postgresql_init.sql

# To verify the data was inserted correctly
psql -d family_service
family_service=# SELECT COUNT(*) FROM families;  # Should return 4
family_service=# SELECT id, status FROM families;
```

**SQLite Script:**
```bash
# Run the initialization script
go run data/dev/sqlite/sqlite_init.go

# Alternative method with custom database path
SQLITE_URI="file:/path/to/custom.db?cache=shared&mode=rwc" go run data/dev/sqlite/sqlite_init.go

# To verify the data was inserted correctly
sqlite3 data/dev/sqlite/family_service.db
sqlite> SELECT COUNT(*) FROM families;  # Should return 4
sqlite> SELECT id, status FROM families;
```

#### Troubleshooting

**Common MongoDB Issues:**
- **Connection refused**: Ensure MongoDB is running with `sudo systemctl status mongodb` or `brew services list`
- **Authentication errors**: If your MongoDB requires authentication, use `mongosh -u username -p password --authenticationDatabase admin data/dev/mongo/mongodb_init.js`
- **Script errors**: Check MongoDB logs with `tail -f /var/log/mongodb/mongodb.log` or equivalent

**Common PostgreSQL Issues:**
- **Database doesn't exist**: Run `createdb family_service` before running the script
- **Permission denied**: Ensure your user has the right permissions with `psql -c "ALTER USER your_username WITH SUPERUSER;"`
- **Port conflicts**: If using a non-standard port, specify with `psql -p 5433 -d family_service -f data/dev/postgres/postgresql_init.sql`

**Common SQLite Issues:**
- **Go package not found**: Ensure you're running the script from the project root directory
- **Permission denied**: Check file permissions with `ls -la data/dev/sqlite/family_service.db` and ensure you have write access to the directory
- **Database locked**: Ensure no other process is using the database file with `lsof | grep data/dev/sqlite/family_service.db`
- **SQLite not installed**: The script uses the Go SQLite driver, but you need the SQLite CLI tool to verify the data: `brew install sqlite` or `apt-get install sqlite3`

#### Customizing the Scripts

You can modify the initialization scripts to suit your specific needs:

**Modifying Sample Data:**
- Edit the `db.families.insertMany([...])` array in `mongodb_init.js` to add, remove, or modify family records
- Edit the `INSERT INTO families` statements in `postgresql_init.sql` to change the sample data
- Modify the `createFamily1()`, `createFamily2()`, etc. functions in `sqlite_init.go` to change the sample data

**Changing Database Structure:**
- Modify the index creation in `mongodb_init.js` to optimize for your query patterns
- Adjust the table schema, indexes, or triggers in `postgresql_init.sql` to match your requirements
- The SQLite table structure is defined in the `ensureTableExists()` method in the repository, but you can add custom indexes or triggers in the script if needed

**Adding More Collections/Tables:**
- For MongoDB: Add more `db.collection.drop()` and `db.collection.createIndex()` calls
- For PostgreSQL: Add more `CREATE TABLE` and `CREATE INDEX` statements
- For SQLite: Add more table creation logic and corresponding entity creation functions

Remember to update all scripts to maintain consistency between MongoDB, PostgreSQL, and SQLite implementations.

These scripts can also run automatically as part of the Docker deployment as shown in the deployment section below.

### Environment Variables (`.env`)

```env
DB_DRIVER=sqlite|postgres|mongo
SQLITE_URI=file:data/dev/sqlite/family_service.db?cache=shared&mode=rwc
DB_URI=postgres://user:pass@localhost:5432/familydb
MONGO_URI=mongodb://localhost:27017
```

### Retry Configuration

The Family Service includes configurable retry logic for database operations to handle transient errors. The retry configuration is defined in the application configuration files:

```yaml
retry:
  max_retries: 3              # Maximum number of retry attempts
  initial_backoff: 100ms      # Initial backoff duration before the first retry
  max_backoff: 1s             # Maximum backoff duration for any retry
```

These settings can be overridden using environment variables:

```env
APP_RETRY_MAX_RETRIES=3
APP_RETRY_INITIAL_BACKOFF=100ms
APP_RETRY_MAX_BACKOFF=1s
```

The retry logic is implemented in all repository adapters (MongoDB, PostgreSQL, SQLite) and uses an exponential backoff strategy with jitter to prevent thundering herd problems. The retry mechanism automatically handles transient errors such as network issues, timeouts, and temporary database unavailability.

Retries are only attempted for operations that are safe to retry (idempotent operations) and for specific error types that are likely to be transient. Permanent errors such as validation failures or not found errors are not retried.

### Make Commands

```bash
make build                       # Builds the service (automatically regenerates GraphQL code)
make test                        # Runs all unit/integration tests
make lint                        # Lints the codebase
make run                         # Starts the service
make db-init                     # Initializes the database based on DB_DRIVER environment variable
make graphql-gen                 # Regenerates GraphQL code from schema.graphql
make plantuml                    # Regenerates all SVG diagrams from PlantUML files
make plantuml-deployment-container # Regenerates only the Deployment Container Diagram SVG file
```

### GraphQL Code Generation

The project uses [gqlgen](https://github.com/99designs/gqlgen) to generate Go code from the GraphQL schema. The GraphQL schema is defined in `interface/adapters/graphql/schema.graphql`.

GraphQL code generation is now automatically performed in the following scenarios:

1. **During Development**: When using Air for live reloading, any changes to `.graphql` files will automatically trigger code regeneration before rebuilding the application.

2. **During Build**: When running `make build`, the GraphQL code is automatically regenerated before building the application.

3. **During Docker Build**: When building the Docker image, the GraphQL code is automatically regenerated as part of the build process.

You can also manually regenerate the GraphQL code at any time by running:

```bash
make graphql-gen
```

This ensures that the generated code is always in sync with the GraphQL schema, preventing inconsistencies and errors that could occur if the code generation step is forgotten.

## 🚀 Deployment Guide

### Secrets Setup (Required)

Before deploying the application, you must set up the required secrets files:

1. Create a `secrets` folder in the project root if it doesn't exist
2. Create all the required files as described in the [Secrets Setup Guide](docs/Secrets_Setup_Guide.md)
3. Ensure the secrets folder is not committed to the repository (it's already in .gitignore)

### Option 1: Local Deployment (Manual)

1. Start MongoDB or PostgreSQL (based on your `.env` configuration)
2. Initialize the database with sample data:

```bash
# For MongoDB
mongosh < data/dev/mongo/mongodb_init.js

# For PostgreSQL
# Note: Make sure the database exists first
createdb familydb  # Create the database if it doesn't exist
psql -d familydb -f data/dev/postgres/postgresql_init.sql
```

3. Verify the initialization was successful:

```bash
# For MongoDB
mongosh
> use family_service
> db.families.count()  # Should return 4

# For PostgreSQL
psql -d familydb
familydb=# SELECT COUNT(*) FROM families;  # Should return 4
```

4. Build and run the service:

```bash
make build
make run
```

The service will connect to the database specified in your `.env` file and use the initialized data.

### Option 2: Docker-based Deployment

#### Docker Compose (recommended for local/dev)

```yaml
version: '3.8'
services:
  family-service:
    build: .
    ports:
      - "8089:8089"
    env_file:
      - .env
    depends_on:
      - mongodb
      - postgres
    volumes:
      - ./data/dev/sqlite/family_service.db:/app/data/dev/sqlite/family_service.db  # Mount SQLite database file

  # Service to initialize SQLite database
  sqlite-init:
    build: .
    command: go run data/dev/sqlite/sqlite_init.go
    volumes:
      - ./data/dev/sqlite/family.db:/app/data/dev/sqlite/family.db  # Mount SQLite database file
    environment:
      - SQLITE_URI=file:/app/data/dev/sqlite/family.db?cache=shared&mode=rwc

  mongodb:
    image: mongo:6
    ports:
      - "27017:27017"
    volumes:
      - ./data/dev/mongo:/data/dev/mongo
    # This command sequence:
    # 1. Starts MongoDB in background
    # 2. Runs the initialization script
    # 3. Restarts MongoDB to ensure clean state
    command: bash -c "mongod --fork --logpath /var/log/mongodb.log && mongosh /data/dev/mongo/mongodb_init.js && mongod --shutdown && mongod"

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: familydb
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    ports:
      - "5432:5432"
    volumes:
      - ./data/dev/postgres:/data/dev/postgres
    # This command sequence:
    # 1. Starts PostgreSQL in background
    # 2. Waits for it to be ready
    # 3. Runs the initialization script
    # 4. Waits for the PostgreSQL process to complete
    command: bash -c "docker-entrypoint.sh postgres & sleep 10 && psql -U user -d familydb -f /data/dev/postgres/postgresql_init.sql && wait"
```

Then run:

```bash
docker-compose up --build
```

This configuration will automatically initialize the databases with sample data using the provided scripts.

## 🔐 Authentication and Authorization

The Family Service includes built-in support for authentication and authorization using JWT (JSON Web Tokens). All requests to the service must be authenticated and provide authorization roles in a JWT token.

### JWT Authentication

The service uses the `servicelib/auth` package to validate JWT tokens locally. The token validation is applied as middleware to all routes, ensuring that only authenticated users can access the service.

### Configuration

JWT authentication is configured in the application configuration:

```yaml
auth:
  jwt:
    secret_key: "your-secret-key-here"  # Should be overridden in production
    token_duration: "24h"               # Token validity period
    issuer: "family-service"            # Token issuer
```

These settings can be overridden using environment variables:

```env
APP_AUTH_JWT_SECRET_KEY=your-secret-key
APP_AUTH_JWT_TOKEN_DURATION=24h
APP_AUTH_JWT_ISSUER=family-service
```

### Role-Based Authorization

The service supports role-based authorization with three roles:

1. **admin**: Has full access to all operations
2. **editor**: Has access to create, update, and delete operations
3. **viewer**: Has read-only access to data
4. **non-user**: Has no access (unauthenticated)

### Future Improvements

In the future, the service will be configured to use a remote authorization server instead of local token validation for improved security and centralized management.

## 📊 Monitoring and Observability

The Family Service includes built-in support for monitoring and observability using Prometheus and Grafana. This allows you to collect and visualize metrics about the application's performance and behavior.

### Prometheus Integration

Prometheus is used to scrape and store metrics from the Family Service. The service exposes metrics at the `/metrics` endpoint, which Prometheus scrapes at regular intervals as defined in `prometheus.yml`.

Once the application is deployed, you can access the Prometheus UI at:
```
http://localhost:9090
```

From here, you can query metrics and create graphs:

![Prometheus Go Info Query](docs/assets/prometheus_go_info_query.png)

You can also check that Prometheus is successfully scraping metrics from the Family Service:

![Prometheus Target Query](docs/assets/prometheus_target_query.png)

### Grafana Integration

Grafana is used to create dashboards that visualize the metrics collected by Prometheus. The Family Service includes a custom Grafana dashboard that provides insights into the application's performance.

Once the application is deployed, you can access the Grafana UI at:
```
http://localhost:3000
```

The Family Service includes a custom Grafana dashboard that visualizes key metrics:

![Grafana Dashboard](docs/assets/grafana_dashboard_for_family_service.png)

The dashboard configuration is available at `docs/assets/grafana_dashboard_for_family_service.json` and can be imported into Grafana.

### Available Metrics

The Family Service exposes the following metrics:

- **Go Runtime Metrics**: Heap allocations, memory usage, goroutines, etc.
- **HTTP Metrics**: Request counts, durations, and in-flight requests
- **Database Metrics**: Operation counts, durations, and connection pools
- **Application Metrics**: Error counts and custom business metrics

For more detailed information about monitoring and observability, refer to the [Deployment Document](docs/Deployment_FamilyService.md).

### Health Check

Visit: `http://localhost:8089/healthz`

## 🔍 Using GraphiQL

GraphiQL is an in-browser IDE for exploring GraphQL APIs. It provides a user-friendly interface to write queries, mutations, and view schema documentation.

### Accessing GraphiQL

After starting the service, you can access GraphiQL in two ways:

1. **Direct GraphiQL Interface**: Visit `http://localhost:8089/graphiql`
2. **GraphQL Playground** (alternative interface): Visit `http://localhost:8089/playground`

Both interfaces connect to the same GraphQL endpoint at `http://localhost:8089/query`.

### Basic Usage

1. **Exploring the Schema**: 
   - Click on the "Docs" button in the top-right corner to open the Documentation Explorer
   - Browse through available queries, mutations, and types

2. **Writing Queries**:
   - Type your GraphQL query in the left panel
   - Press the "Play" button or use the keyboard shortcut (Ctrl+Enter) to execute
   - Results will appear in the right panel

3. **Using Variables**:
   - Open the variables panel at the bottom
   - Define variables in JSON format
   - Reference them in your query with `$variableName`

### Example Queries

Here are some example queries you can try in GraphiQL:

**Get All Families**:
```graphql
{
  getAllFamilies {
    id
    status
    parents {
      id
      firstName
      lastName
    }
    children {
      id
      firstName
      lastName
    }
  }
}
```

**Get a Specific Family**:
```graphql
{
  getFamily(id: "fam-123456789") {
    id
    status
    parents {
      firstName
      lastName
      birthDate
    }
    children {
      firstName
      lastName
      birthDate
    }
  }
}
```

**Create a New Family (Mutation)**:
```graphql
mutation {
  createFamily(input: {
    id: "fam-new-123"
    status: SINGLE
    parents: [{
      id: "par-new-123"
      firstName: "John"
      lastName: "Doe"
      birthDate: "1980-01-01T00:00:00Z"
    }]
    children: []
  }) {
    id
    status
    parents {
      id
      firstName
      lastName
    }
  }
}
```

### Troubleshooting

If you encounter issues with GraphiQL:

1. **"Failed to fetch" error**:
   - Ensure the server is running (`make run`)
   - Check that you're using the correct URL
   - Verify there are no CORS issues (try using the built-in GraphiQL at `/graphiql` instead of an external client)

2. **Authentication errors**:
   - This API currently doesn't require authentication

3. **Invalid query syntax**:
   - Use the schema documentation to verify field names and types
   - Check for typos in field names
   - Ensure all brackets and braces are properly closed

4. **Browser console errors**:
   - Open your browser's developer tools (F12) to check for detailed error messages
   - Look for network errors or JavaScript exceptions

For more complex queries and mutations, refer to the GraphQL schema documentation available within GraphiQL.

## 📄 Software Requirements Specification (SRS)

The SRS outlines the essential business rules and functional constraints for managing families:

- A family must contain at least one parent
- Children may not exist without belonging to a family
- A parent may belong to multiple families (due to divorce/remarriage)
- A child belongs to only one family at a time
- Family lifecycle states: `single`, `married`, `divorced`, `widowed`
- On divorce: a new family is created for the custodial parent and children
- Each person must have a first name, last name, birthdate; deathdate optional
- No duplicate parents in a family (based on name + birthdate)
- Validation must occur in adapters and domain layer

UML References: 
- [SRS Use Case Diagram](./docs/diagrams/SRS%20Use%20Case%20Diagram.svg)
- [SRS Data Model Diagram](./docs/diagrams/SRS%20Data%20Model%20Diagram.svg)

## 📘 Software Design Document (SDD)

The SDD describes the architecture and technologies used:

- DDD drives entity and aggregate modeling (Family, Parent, Child)
- Clean Architecture provides layered structure:
  - `domain`: pure logic
  - `services`: coordination logic
  - `infrastructure`: outermost layer with adapters and cross-cutting concerns
  - `adapters`: database implementations
- Hexagonal Architecture enables easy persistence switching via `ports`
- SQLite provides a local file-based database for development
- MongoDB stores parents and children embedded in a `family` document
- PostgreSQL uses normalized schema with FK constraints
- Zap is used for structured logging
- OS signals and context are used for graceful shutdown

UML References:
- [SDD Class Diagram](./docs/diagrams/SDD%20Class%20Diagram.svg)
- [SDD Sequence Diagram - Divorce Operation](./docs/diagrams/SDD%20Sequence%20Diagram%20-%20Divorce%20Operation.svg)

## 🧪 Software Test Plan (STP)

The STP defines testing strategy:

- **Unit Tests**:
  - `core/domain`: validation of invariants
  - `core/services`: orchestration of complex transitions (e.g. divorce)
- **Repository Tests**:
  - Mongo and Postgres persistence checks
- **Integration Tests**:
  - GraphQL layer against real DBs
  - Transaction rollback validation
- **End-to-End Tests**:
  - From mutation/query to DB write/read

Coverage Target: 90%+ statements

Edge Cases:
- Empty names, invalid birthdates
- Children without families
- Divorce with no children

UML References:
- [STP Test Process Activity Diagram](./docs/diagrams/STP%20Test%20Process%20Activity%20Diagram.svg)
- [STP Test Coverage Component Diagram](./docs/diagrams/STP%20Test%20Coverage%20Component%20Diagram.svg)

## 🚢 Deployment Document

The Deployment Document provides detailed instructions for deploying the system:

- Docker-based containerization for all components
- Environment configuration via `.env` files
- MongoDB and PostgreSQL database setup
- Health checks and monitoring
- Backup and restore procedures
- Security considerations

UML References:
- [Deployment Diagram](./docs/diagrams/Deployment%20Diagram.svg)
- [Deployment Container Diagram](./docs/diagrams/Deployment%20Container%20Diagram.svg)

## ✅ Testing

- Unit, repository, integration, and end-to-end tests
- Target ≥ 90% test coverage
- Run: `go test ./... -coverprofile=coverage.out`
- Coverage badge and logs in output

## 📝 Design Notes

- Children belong to only one family
- Parents can be in multiple families
- Family state transitions include divorce, death, remarriage
- Transactional support for both databases
- Context propagation and OS signal handling integrated
