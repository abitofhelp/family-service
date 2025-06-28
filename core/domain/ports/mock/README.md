# Repository Mocks

## Overview

The Repository Mocks package provides mock implementations of the repository interfaces defined in the domain ports package. These mocks are generated using GoMock and are designed to facilitate testing of components that depend on these interfaces. By using these mocks, developers can easily simulate different repository behaviors and test edge cases without relying on actual database implementations.

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

## Best Practices

1. **Initialize in Test Setup**: Create mock instances in test setup functions
2. **Set Expectations Before Use**: Set expectations before calling the code under test
3. **Verify After Use**: Verify that all expectations were met after the test
4. **Use Argument Matchers**: Use argument matchers for flexible matching
5. **Keep Tests Focused**: Test one behavior at a time with specific expectations

## Related Components

- [Domain Ports](../README.md) - Defines the interfaces that these mocks implement
- [Domain Services](../../services/README.md) - Services that can be tested using these mocks
- [Application Services](../../../application/services/README.md) - Higher-level services that can be tested using these mocks
- [GoMock Documentation](https://github.com/golang/mock) - Documentation for the GoMock library