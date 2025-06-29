# Repository Mocks

## Overview

The Repository Mocks package provides mock implementations of the repository interfaces defined in the domain ports package. These mocks are generated using GoMock and are designed to facilitate testing of components that depend on these interfaces. By using these mocks, developers can easily simulate different repository behaviors and test edge cases without relying on actual database implementations.

## Architecture

The Repository Mocks package is part of the testing infrastructure for the domain layer in the Clean Architecture and Hexagonal Architecture patterns. It sits alongside the domain ports package and provides mock implementations of the interfaces defined there. The architecture follows these principles:

- **Hexagonal Architecture (Ports and Adapters)**: Provides mock implementations of the ports (interfaces) defined in the domain layer
- **Dependency Inversion Principle**: Allows tests to inject mock implementations of the interfaces that the domain layer depends on
- **Clean Architecture**: Supports testing of the domain layer without relying on infrastructure implementations
- **Test-Driven Development**: Facilitates writing tests before implementing the actual infrastructure adapters

The package is organized into:

- **Generated Mock Types**: Mock implementations of the repository interfaces
- **Mock Recorders**: Types that record method calls and expectations
- **Factory Functions**: Functions that create new mock instances

## Implementation Details

The Repository Mocks package implements the following design patterns:

1. **Mock Object Pattern**: Provides objects that simulate the behavior of real objects for testing
2. **Factory Pattern**: Factory functions create new mock instances
3. **Recorder Pattern**: Recorders track method calls and expectations
4. **Proxy Pattern**: Mocks act as proxies for the real implementations

Key implementation details:

- **GoMock Generation**: Mocks are automatically generated using GoMock's mockgen tool
- **Interface Reflection**: GoMock uses reflection to analyze the interfaces and generate appropriate mocks
- **Expectation API**: Mocks provide a fluent API for setting expectations
- **Argument Matching**: Mocks support flexible argument matching using GoMock's matchers
- **Call Counting**: Mocks track the number of times each method is called
- **Return Value Specification**: Mocks allow specifying return values for method calls

## Features

- **Generated Mocks**: Automatically generated using GoMock for consistent behavior
- **Full Interface Coverage**: Implements all methods defined in the repository interfaces
- **Expectation Setting**: Allows setting expectations for method calls
- **Return Value Control**: Enables specifying return values for method calls
- **Call Verification**: Supports verifying that expected methods were called
- **Argument Matching**: Provides flexible argument matching for method calls

## API Documentation

### Core Types

#### MockFamilyRepository

A mock implementation of the FamilyRepository interface.

```
// MockFamilyRepository is a mock of FamilyRepository interface
type MockFamilyRepository struct {
    ctrl     *gomock.Controller
    recorder *MockFamilyRepositoryMockRecorder
}

// NewMockFamilyRepository creates a new mock instance
func NewMockFamilyRepository(ctrl *gomock.Controller) *MockFamilyRepository
```

### Key Methods

#### EXPECT

Returns an object that allows setting expectations on the mock.

```
// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFamilyRepository) EXPECT() *MockFamilyRepositoryMockRecorder
```

#### FindByChildID

Mocks the FindByChildID method of the FamilyRepository interface.

```
// FindByChildID mocks base method
func (m *MockFamilyRepository) FindByChildID(ctx context.Context, childID string) (*entity.Family, error)
```

#### FindByParentID

Mocks the FindByParentID method of the FamilyRepository interface.

```
// FindByParentID mocks base method
func (m *MockFamilyRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error)
```

## Examples

There may be additional examples in the /EXAMPLES directory.

Example of using the MockFamilyRepository:

```
// Create a mock controller
ctrl := gomock.NewController(t)
defer ctrl.Finish()

// Create a mock repository
mockRepo := mock.NewMockFamilyRepository(ctrl)

// Set up expectations for GetByID
mockRepo.EXPECT().
    GetByID(gomock.Any(), "family-123").
    Return(&entity.Family{
        // Family fields
    }, nil).
    Times(1)

// Set up expectations for FindByParentID
mockRepo.EXPECT().
    FindByParentID(gomock.Any(), "parent-456").
    Return([]*entity.Family{
        // Family list
    }, nil).
    Times(1)

// Set up expectations for Save
mockRepo.EXPECT().
    Save(gomock.Any(), gomock.Any()).
    Return(nil).
    Times(1)

// Use the mock repository in your test
family, err := mockRepo.GetByID(ctx, "family-123")
assert.NoError(t, err)
assert.NotNil(t, family)

families, err := mockRepo.FindByParentID(ctx, "parent-456")
assert.NoError(t, err)
assert.NotEmpty(t, families)

err = mockRepo.Save(ctx, family)
assert.NoError(t, err)
```

## Configuration

The Repository Mocks package doesn't require any specific configuration as it provides mock implementations for testing. However, there are some configuration options for the GoMock library itself:

- **Controller Configuration**: Configure the behavior of the mock controller
- **Expectation Configuration**: Configure how expectations are matched and verified
- **Argument Matcher Configuration**: Configure how arguments are matched

Example configuration:

```
// Create a mock controller with custom options
ctrl := gomock.NewController(t)
defer ctrl.Finish()

// Configure the controller to be strict (default)
// This means that unexpected calls will cause the test to fail
ctrl.Strict = true

// Create a mock repository with the configured controller
mockRepo := mock.NewMockFamilyRepository(ctrl)
```

## Testing

The Repository Mocks package is itself used for testing other components. It is tested through:

1. **Generated Tests**: GoMock generates tests for the mock implementations
2. **Integration Tests**: Tests that verify the mocks work correctly with the components they are testing
3. **Example Tests**: Tests that demonstrate how to use the mocks

Key testing approaches:

- **Expectation Testing**: Tests that verify expectations are correctly set and matched
- **Argument Matching Testing**: Tests that verify argument matchers work correctly
- **Return Value Testing**: Tests that verify return values are correctly specified and returned
- **Call Count Testing**: Tests that verify call counts are correctly tracked

Example of testing a component using the mocks:

```
func TestFamilyService_GetByID(t *testing.T) {
    // Create a mock controller
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Create a mock repository
    mockRepo := mock.NewMockFamilyRepository(ctrl)

    // Set up expectations
    expectedFamily := &entity.Family{
        // Family fields
    }
    mockRepo.EXPECT().
        GetByID(gomock.Any(), "family-123").
        Return(expectedFamily, nil).
        Times(1)

    // Create the service with the mock repository
    service := services.NewFamilyService(mockRepo)

    // Call the method being tested
    family, err := service.GetByID(context.Background(), "family-123")

    // Verify the results
    assert.NoError(t, err)
    assert.Equal(t, expectedFamily, family)
}
```

## Design Notes

1. **Generated Code**: The mocks are generated code and should not be modified manually
2. **Interface Coupling**: The mocks are tightly coupled to the interfaces they implement
3. **Regeneration**: The mocks need to be regenerated when the interfaces change
4. **Testing Focus**: The mocks are designed for testing, not for production use
5. **Behavior Verification**: The mocks support behavior verification (verifying that methods are called with specific arguments)
6. **State Verification**: The mocks can also support state verification (verifying that the system under test changes state correctly)
7. **Mockgen Tool**: The mocks are generated using the mockgen tool from the GoMock library

## Best Practices

1. **Initialize in Test Setup**: Create mock instances in test setup functions
2. **Set Expectations Before Use**: Set expectations before calling the code under test
3. **Verify After Use**: Verify that all expectations were met after the test
4. **Use Argument Matchers**: Use argument matchers for flexible matching
5. **Keep Tests Focused**: Test one behavior at a time with specific expectations

## References

- [GoMock Documentation](https://github.com/golang/mock) - Documentation for the GoMock library
- [Mockgen Tool](https://github.com/golang/mock/tree/master/mockgen) - Tool for generating mocks
- [Test-Driven Development](https://en.wikipedia.org/wiki/Test-driven_development) - Development methodology that relies heavily on mocks
- [Behavior-Driven Development](https://en.wikipedia.org/wiki/Behavior-driven_development) - Development methodology that uses mocks for behavior verification
- [Domain Ports](../README.md) - Defines the interfaces that these mocks implement
- [Domain Services](../../services/README.md) - Services that can be tested using these mocks
- [Application Services](../../../application/services/README.md) - Higher-level services that can be tested using these mocks
- [Testing Best Practices](https://golang.org/doc/testing) - Go's testing best practices
