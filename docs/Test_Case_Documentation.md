# Test Case Documentation Guide

## Overview
This document provides detailed documentation for test cases and coverage goals for the Family Service project. It serves as a guide for writing effective tests, achieving coverage targets, and maintaining consistent test documentation across the codebase.

## Table of Contents
1. [Test Documentation Standards](#test-documentation-standards)
2. [Coverage Goals and Strategies](#coverage-goals-and-strategies)
3. [Test Case Templates](#test-case-templates)
4. [Examples of Well-Documented Tests](#examples-of-well-documented-tests)
5. [Package-Specific Testing Guidelines](#package-specific-testing-guidelines)
6. [Troubleshooting Common Testing Issues](#troubleshooting-common-testing-issues)

## Test Documentation Standards

### Test Function Documentation
Each test function should include:

1. **Purpose**: A clear description of what the test is verifying
2. **Coverage**: Which parts of the code are being covered
3. **Edge Cases**: What edge cases are being tested
4. **Dependencies**: Any dependencies or setup required

Example:

    // TestSaveFamily tests the Save method of the FamilyRepository.
    // It covers:
    // - Successful saving of a family with parents and children
    // - Error handling for invalid family data
    // - Transaction rollback on failure
    // Edge cases:
    // - Family with maximum number of parents
    // - Family with no children
    // Dependencies:
    // - Requires a properly initialized database connection
    func TestSaveFamily(t *testing.T) {
        // Test implementation
    }

### Test Case Documentation
For table-driven tests, each test case should include:

1. **Name**: A descriptive name that explains the scenario
2. **Input**: The input data and conditions
3. **Expected Output**: The expected result or behavior
4. **Reason**: Why this test case is important

Example:

    testCases := []struct {
        name           string // Descriptive name of the test case
        family         *entity.Family // Input family data
        expectedErr    bool // Whether an error is expected
        expectedErrMsg string // Expected error message if an error is expected
        reason         string // Why this test case is important
    }{
        {
            name:           "Valid family with two parents",
            family:         createValidFamilyWithTwoParents(),
            expectedErr:    false,
            expectedErrMsg: "",
            reason:         "Tests the happy path with a valid family structure",
        },
        {
            name:           "Family with no parents",
            family:         createFamilyWithNoParents(),
            expectedErr:    true,
            expectedErrMsg: "family must have at least one parent",
            reason:         "Tests validation of the minimum parent requirement",
        },
    }

## Coverage Goals and Strategies

### Overall Coverage Targets
As specified in the Software Test Plan (STP), the coverage targets are:

- Domain layer: 95%+ statement coverage
- Service layer: 90%+ statement coverage
- Repository layer: 85%+ statement coverage
- GraphQL layer: 80%+ statement coverage
- Overall: 90%+ statement coverage

### Strategies for Achieving Coverage

#### Domain Layer (95%+)
1. **Entity Tests**:
   - Test all constructors with valid and invalid inputs
   - Test all methods that modify entity state
   - Test validation rules
   - Test business logic methods

2. **Value Object Tests**:
   - Test creation with valid and invalid values
   - Test equality and comparison methods
   - Test formatting and parsing methods

3. **Domain Service Tests**:
   - Test all business operations
   - Test all validation rules
   - Test interaction between entities

#### Service Layer (90%+)
1. **Application Service Tests**:
   - Test all use cases
   - Test error handling
   - Test transaction management
   - Mock dependencies (repositories, domain services)

2. **Integration Tests**:
   - Test interaction between services and repositories
   - Test complete workflows

#### Repository Layer (85%+)
1. **Repository Implementation Tests**:
   - Test all CRUD operations
   - Test query methods
   - Test error handling
   - Test transaction management

2. **Database Adapter Tests**:
   - Test connection management
   - Test query execution
   - Test error handling

#### GraphQL Layer (80%+)
1. **Resolver Tests**:
   - Test all queries and mutations
   - Test error handling
   - Test authorization
   - Mock dependencies (services)

2. **Schema Tests**:
   - Test schema validation
   - Test type definitions

### Measuring and Reporting Coverage
1. **Running Coverage Tests**:

       go test -coverprofile=coverage.out ./...
       go tool cover -html=coverage.out -o coverage.html

2. **Interpreting Coverage Reports**:
   - Identify uncovered lines and functions
   - Prioritize coverage of error handling and edge cases
   - Focus on business-critical code paths

3. **Continuous Integration**:
   - Set up coverage thresholds in CI/CD pipeline
   - Fail builds that don't meet coverage targets
   - Generate coverage reports for each build

## Test Case Templates

### Unit Test Template

    // Test<FunctionName> tests the <FunctionName> function/method.
    // It covers:
    // - <What aspects of the function are being tested>
    // Edge cases:
    // - <What edge cases are being tested>
    // Dependencies:
    // - <Any dependencies or setup required>
    func Test<FunctionName>(t *testing.T) {
        // Setup
        // ...

        // Test cases
        testCases := []struct {
            name           string
            input          <InputType>
            expected       <ExpectedType>
            expectedErr    bool
            expectedErrMsg string
            reason         string
        }{
            // Test cases...
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                // Test execution
                result, err := <FunctionName>(tc.input)

                // Assertions
                if tc.expectedErr {
                    require.Error(t, err)
                    assert.Contains(t, err.Error(), tc.expectedErrMsg)
                } else {
                    require.NoError(t, err)
                    assert.Equal(t, tc.expected, result)
                }
            })
        }
    }

### Integration Test Template

    // Test<Component>Integration tests the integration between <Component> and its dependencies.
    // It covers:
    // - <What aspects of the integration are being tested>
    // Edge cases:
    // - <What edge cases are being tested>
    // Dependencies:
    // - <Any dependencies or setup required>
    func Test<Component>Integration(t *testing.T) {
        // Setup
        // ...

        // Test cases
        testCases := []struct {
            name           string
            input          <InputType>
            expected       <ExpectedType>
            expectedErr    bool
            expectedErrMsg string
            reason         string
        }{
            // Test cases...
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                // Test execution
                result, err := <Component>.<Method>(tc.input)

                // Assertions
                if tc.expectedErr {
                    require.Error(t, err)
                    assert.Contains(t, err.Error(), tc.expectedErrMsg)
                } else {
                    require.NoError(t, err)
                    assert.Equal(t, tc.expected, result)
                }
            })
        }
    }

## Examples of Well-Documented Tests

### Example 1: Entity Test

    // TestFamilyCreation tests the creation of a Family entity.
    // It covers:
    // - Creating a family with valid data
    // - Validation of family data during creation
    // - Error handling for invalid family data
    // Edge cases:
    // - Family with maximum number of parents
    // - Family with no children
    // Dependencies:
    // - None
    func TestFamilyCreation(t *testing.T) {
        // Test cases
        testCases := []struct {
            name           string
            parents        []entity.Parent
            children       []entity.Child
            status         entity.Status
            expectedErr    bool
            expectedErrMsg string
            reason         string
        }{
            {
                name:           "Valid family with two parents and children",
                parents:        createValidParents(2),
                children:       createValidChildren(2),
                status:         entity.Married,
                expectedErr:    false,
                expectedErrMsg: "",
                reason:         "Tests the happy path with a valid family structure",
            },
            {
                name:           "Family with no parents",
                parents:        []entity.Parent{},
                children:       createValidChildren(2),
                status:         entity.Single,
                expectedErr:    true,
                expectedErrMsg: "family must have at least one parent",
                reason:         "Tests validation of the minimum parent requirement",
            },
            {
                name:           "Family with too many parents",
                parents:        createValidParents(3),
                children:       createValidChildren(2),
                status:         entity.Married,
                expectedErr:    true,
                expectedErrMsg: "family cannot have more than two parents",
                reason:         "Tests validation of the maximum parent limit",
            },
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                // Test execution
                family, err := entity.NewFamily(
                    uuid.New().String(),
                    tc.status,
                    tc.parents,
                    tc.children,
                )

                // Assertions
                if tc.expectedErr {
                    require.Error(t, err)
                    assert.Contains(t, err.Error(), tc.expectedErrMsg)
                    assert.Nil(t, family)
                } else {
                    require.NoError(t, err)
                    assert.NotNil(t, family)
                    assert.Equal(t, tc.status, family.Status())
                    assert.Equal(t, len(tc.parents), len(family.Parents()))
                    assert.Equal(t, len(tc.children), len(family.Children()))
                }
            })
        }
    }

### Example 2: Repository Test

    // TestMongoFamilyRepository_GetByID tests the GetByID method of the MongoFamilyRepository.
    // It covers:
    // - Retrieving a family by ID
    // - Error handling for non-existent families
    // - Error handling for database connection issues
    // Edge cases:
    // - Invalid ID format
    // - ID that doesn't exist in the database
    // Dependencies:
    // - MongoDB connection
    // - Properly initialized repository
    func TestMongoFamilyRepository_GetByID(t *testing.T) {
        // Setup
        ctx := context.Background()
        repo, client := setupMongoTest(t)
        defer client.Disconnect(ctx)

        // Create test data
        family := createTestFamily()
        err := repo.Save(ctx, family)
        require.NoError(t, err)

        // Test cases
        testCases := []struct {
            name           string
            id             string
            expectedErr    bool
            expectedErrMsg string
            reason         string
        }{
            {
                name:           "Existing family ID",
                id:             family.ID(),
                expectedErr:    false,
                expectedErrMsg: "",
                reason:         "Tests retrieving an existing family",
            },
            {
                name:           "Non-existent family ID",
                id:             uuid.New().String(),
                expectedErr:    true,
                expectedErrMsg: "family not found",
                reason:         "Tests error handling for non-existent families",
            },
            {
                name:           "Invalid family ID",
                id:             "invalid-id",
                expectedErr:    true,
                expectedErrMsg: "invalid family ID",
                reason:         "Tests error handling for invalid ID format",
            },
        }

        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                // Test execution
                result, err := repo.GetByID(ctx, tc.id)

                // Assertions
                if tc.expectedErr {
                    require.Error(t, err)
                    assert.Contains(t, err.Error(), tc.expectedErrMsg)
                    assert.Nil(t, result)
                } else {
                    require.NoError(t, err)
                    assert.NotNil(t, result)
                    assert.Equal(t, tc.id, result.ID())
                }
            })
        }
    }

## Package-Specific Testing Guidelines

### infrastructure/adapters/mongo
1. **Test Setup**:
   - Use a MongoDB test container or mock
   - Initialize the repository with proper configuration
   - Clean up test data after each test

2. **Test Coverage**:
   - Test all repository methods (GetByID, Save, FindByParentID, FindByChildID, GetAll)
   - Test document conversion (entityToDocument, documentToEntity)
   - Test error handling for database operations
   - Test transaction management

3. **Common Issues**:
   - Connection failures
   - Document conversion errors
   - Validation errors during entity creation

### infrastructure/adapters/postgres
1. **Test Setup**:
   - Use a PostgreSQL test container or mock
   - Initialize the repository with proper configuration
   - Clean up test data after each test

2. **Test Coverage**:
   - Test all repository methods (GetByID, Save, FindByParentID, FindByChildID, GetAll)
   - Test JSON serialization and deserialization
   - Test error handling for database operations
   - Test transaction management

3. **Common Issues**:
   - Connection failures
   - JSON parsing errors
   - SQL query errors

### infrastructure/adapters/sqlite
1. **Test Setup**:
   - Use an in-memory SQLite database
   - Initialize the repository with proper configuration
   - Create tables before running tests

2. **Test Coverage**:
   - Test all repository methods (GetByID, Save, FindByParentID, FindByChildID, GetAll)
   - Test table creation and schema management
   - Test error handling for database operations
   - Test transaction management

3. **Common Issues**:
   - Table creation failures
   - SQLite-specific query syntax
   - Transaction isolation issues

### core/domain/entity
1. **Test Coverage**:
   - Test entity creation with valid and invalid data
   - Test all methods that modify entity state
   - Test validation rules
   - Test business logic methods

2. **Common Issues**:
   - Validation failures
   - Business rule violations
   - Immutability issues

### core/domain/services
1. **Test Coverage**:
   - Test all business operations
   - Test validation rules
   - Test interaction between entities
   - Test error handling

2. **Common Issues**:
   - Business rule violations
   - Dependency issues
   - Error handling edge cases

### core/application/services
1. **Test Coverage**:
   - Test all use cases
   - Test error handling
   - Test transaction management
   - Test interaction with domain services and repositories

2. **Common Issues**:
   - Transaction management issues
   - Error propagation
   - Dependency issues

### interface/adapters/graphql
1. **Test Coverage**:
   - Test all queries and mutations
   - Test error handling
   - Test authorization
   - Test schema validation

2. **Common Issues**:
   - Schema validation errors
   - Authorization issues
   - Error handling edge cases

## Troubleshooting Common Testing Issues

### Skipped Tests
If tests are being skipped, it's important to understand why and address the underlying issues:

1. **Validation Issues**:
   - Review the validation rules in the domain entities
   - Ensure test data meets all validation requirements
   - Consider using test-specific factory methods to create valid entities

2. **Environment Issues**:
   - Ensure all required environment variables are set
   - Check database connections and configurations
   - Verify that all dependencies are properly initialized

3. **Timing Issues**:
   - Add appropriate timeouts to tests
   - Use context with deadlines for database operations
   - Consider using retry mechanisms for flaky tests

### Low Coverage
If coverage is below the target, consider these strategies:

1. **Identify Uncovered Code**:
   - Use coverage reports to identify uncovered lines and functions
   - Focus on business-critical code paths first
   - Pay special attention to error handling paths

2. **Add Missing Tests**:
   - Add tests for uncovered functions and methods
   - Add test cases for error conditions
   - Add tests for edge cases

3. **Refactor Existing Tests**:
   - Combine similar tests to reduce duplication
   - Use table-driven tests to cover more scenarios
   - Use parameterized tests for similar test cases with different inputs

### Test Failures
When tests fail, follow these steps to diagnose and fix the issues:

1. **Understand the Failure**:
   - Read the error message carefully
   - Check the line number and function where the failure occurred
   - Review the test case and expected behavior

2. **Isolate the Issue**:
   - Run the failing test in isolation
   - Add debug logging to understand the execution flow
   - Use a debugger if necessary

3. **Fix the Issue**:
   - Update the test if the expected behavior has changed
   - Fix the code if there's a bug
   - Update dependencies if they're causing the issue