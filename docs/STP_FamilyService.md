# Software Test Plan (STP)

## Family Service GraphQL

### 1. Introduction

#### 1.1 Purpose
This document outlines the testing strategy and approach for the Family Service GraphQL application. It provides a comprehensive plan for validating that the system meets its requirements and functions correctly.

#### 1.2 Scope
This test plan covers unit testing, integration testing, and end-to-end testing for the Family Service GraphQL application. It includes test scenarios, test data, and test environments.

#### 1.3 References
- Software Requirements Specification (SRS)
- Software Design Document (SDD)

### 2. Test Strategy

#### 2.1 Testing Levels

##### 2.1.1 Unit Testing
Unit tests focus on testing individual components in isolation:
- Domain entities and aggregates
- Domain services
- Application services
- Repository implementations (with mocks)
- GraphQL resolvers (with mocks)

##### 2.1.2 Integration Testing
Integration tests focus on testing the interaction between components:
- Repository implementations with actual databases
- GraphQL resolvers with actual services
- Service layer with actual repositories

##### 2.1.3 End-to-End Testing
End-to-end tests focus on testing the entire system from the API to the database:
- GraphQL API with actual databases
- Complete workflows (e.g., create family, add parent, divorce)

#### 2.2 Testing Approaches

##### 2.2.1 White Box Testing
- Code coverage analysis
- Path testing
- Boundary value analysis

##### 2.2.2 Black Box Testing
- Equivalence partitioning
- Error guessing
- Use case testing

#### 2.3 Test Environment

##### 2.3.1 Development Environment
- Local machine with Go 1.24+
- Docker for containerized databases
- MongoDB and PostgreSQL instances

##### 2.3.2 CI/CD Environment
- GitHub Actions or similar CI/CD platform
- Containerized test environment
- Automated test execution

### 3. Test Cases

#### 3.1 Unit Test Cases

##### 3.1.1 Domain Layer Tests

###### 3.1.1.1 Family Entity Tests
- Test family creation with valid data
- Test family creation with invalid data (e.g., no parents, too many parents)
- Test adding a parent to a family
- Test adding a child to a family
- Test removing a child from a family
- Test divorce process
- Test marking a parent as deceased

###### 3.1.1.2 Parent Entity Tests
- Test parent creation with valid data
- Test parent creation with invalid data (e.g., empty name, future birth date)
- Test marking a parent as deceased

###### 3.1.1.3 Child Entity Tests
- Test child creation with valid data
- Test child creation with invalid data (e.g., empty name, future birth date)
- Test marking a child as deceased

##### 3.1.2 Service Layer Tests

###### 3.1.2.1 Family Service Tests
- Test creating a family
- Test retrieving a family
- Test adding a parent to a family
- Test adding a child to a family
- Test removing a child from a family
- Test divorce process
- Test marking a parent as deceased
- Test error handling for various scenarios

##### 3.1.3 Repository Layer Tests

###### 3.1.3.1 MongoDB Repository Tests
- Test saving a family
- Test retrieving a family by ID
- Test finding families by parent ID
- Test finding a family by child ID
- Test error handling for various scenarios

###### 3.1.3.2 PostgreSQL Repository Tests
- Test saving a family
- Test retrieving a family by ID
- Test finding families by parent ID
- Test finding a family by child ID
- Test error handling for various scenarios

##### 3.1.4 GraphQL Resolver Tests
- Test createFamily mutation
- Test addParent mutation
- Test addChild mutation
- Test removeChild mutation
- Test markParentDeceased mutation
- Test divorce mutation
- Test getFamily query
- Test findFamiliesByParent query
- Test findFamilyByChild query
- Test error handling for various scenarios

#### 3.2 Integration Test Cases

##### 3.2.1 Repository Integration Tests

###### 3.2.1.1 MongoDB Repository Integration Tests
- Test saving and retrieving a family
- Test finding families by parent ID
- Test finding a family by child ID
- Test transaction handling

###### 3.2.1.2 PostgreSQL Repository Integration Tests
- Test saving and retrieving a family
- Test finding families by parent ID
- Test finding a family by child ID
- Test transaction handling

##### 3.2.2 Service Integration Tests
- Test family service with actual repositories
- Test complex workflows (e.g., create family, add parent, divorce)

##### 3.2.3 GraphQL Integration Tests
- Test GraphQL resolvers with actual services
- Test GraphQL schema validation

#### 3.3 End-to-End Test Cases

##### 3.3.1 API Tests
- Test GraphQL API with actual databases
- Test complete workflows (e.g., create family, add parent, divorce)
- Test error handling and validation

##### 3.3.2 Performance Tests
- Test API response time under load
- Test database performance under load
- Test concurrent operations

### 4. Test Data

#### 4.1 Test Data Generation
- Generate test data for families, parents, and children
- Generate test data for various family structures (single, married, divorced, widowed)
- Generate test data for edge cases (e.g., maximum number of parents, no children)

#### 4.2 Test Data Management
- Reset test data before each test run
- Use unique identifiers for test data
- Clean up test data after test runs

### 5. Test Coverage

#### 5.1 Code Coverage Targets
- Domain layer: 95%+ statement coverage
- Service layer: 90%+ statement coverage
- Repository layer: 85%+ statement coverage
- GraphQL layer: 80%+ statement coverage
- Overall: 90%+ statement coverage

#### 5.2 Functional Coverage
- All use cases from the SRS must be covered by tests
- All error conditions must be covered by tests
- All business rules must be covered by tests

### 6. Test Automation

#### 6.1 Test Frameworks
- Go testing package for unit tests
- Testify for assertions and mocks
- Custom test utilities for common testing tasks

#### 6.2 Continuous Integration
- Automated test execution on every pull request
- Code coverage reporting
- Test result reporting

### 7. Test Execution

#### 7.1 Test Execution Process
1. Run unit tests
2. Run integration tests
3. Run end-to-end tests
4. Generate test reports
5. Analyze test results

#### 7.2 Test Schedule
- Unit tests: Run on every code change
- Integration tests: Run on every pull request
- End-to-end tests: Run on every pull request to main branch

#### 7.3 Test Environment Setup
- Set up test databases (MongoDB and PostgreSQL)
- Configure test environment variables
- Initialize test data
- Create and configure the `secrets` folder with required credential files (see [Secrets Setup Guide](Secrets_Setup_Guide.md))

### 8. Test Reporting

#### 8.1 Test Result Reporting
- Generate test result reports
- Report test coverage
- Report test failures
- Report test execution time

#### 8.2 Defect Tracking
- Log defects in issue tracking system
- Link defects to test cases
- Track defect resolution

### 9. Entry and Exit Criteria

#### 9.1 Entry Criteria
- Code must compile without errors
- All dependencies must be available
- Test environment must be set up
- Secrets folder must be properly configured with all required files

#### 9.2 Exit Criteria
- All tests must pass
- Code coverage must meet targets
- No critical or high-severity defects

### 10. Test Scenarios

#### 10.1 Family Creation Scenarios
- Create a single-parent family
- Create a two-parent family
- Create a family with children
- Create a family with invalid data (should fail)

#### 10.2 Family Modification Scenarios
- Add a parent to a single-parent family
- Add a child to a family
- Remove a child from a family
- Mark a parent as deceased
- Process a divorce

#### 10.3 Query Scenarios
- Get a family by ID
- Find families by parent ID
- Find a family by child ID

#### 10.4 Edge Case Scenarios
- Create a family with maximum number of parents
- Add a parent to a family that already has two parents (should fail)
- Add a duplicate parent to a family (should fail)
- Add a duplicate child to a family (should fail)
- Mark an already deceased parent as deceased (should fail)
- Divorce a single-parent family (should fail)

### 11. Risks and Contingencies

#### 11.1 Risks
- Database connectivity issues
- Test data corruption
- Test environment instability
- Incomplete test coverage
- Missing or incorrectly configured secrets files

#### 11.2 Contingencies
- Backup test databases
- Implement test data reset mechanisms
- Monitor test environment health
- Regularly review and update test coverage
- Provide automated scripts to verify and set up the required secrets files

### 12. Appendices

#### 12.1 Test Process Diagram
![Test Process Diagram](diagrams/STP%20Test%20Process%20Activity%20Diagram.svg)

This activity diagram illustrates the test execution process, showing the flow through unit testing, integration testing, and end-to-end testing phases. It includes decision points for test failures and coverage targets.

#### 12.2 Test Coverage Diagram
![Test Coverage Diagram](diagrams/STP%20Test%20Coverage%20Component%20Diagram.svg)

This component diagram shows which parts of the system are covered by different types of tests, with coverage targets for each component. It illustrates the comprehensive testing strategy described in this document.

#### 12.3 Test Case Templates
- Unit test template
- Integration test template
- End-to-end test template

#### 12.4 Test Data Examples
- Example family data
- Example parent data
- Example child data

#### 12.5 Test Environment Setup Instructions
- MongoDB setup instructions
- PostgreSQL setup instructions
- Docker setup instructions
- Secrets folder setup instructions (see [Secrets Setup Guide](Secrets_Setup_Guide.md))
