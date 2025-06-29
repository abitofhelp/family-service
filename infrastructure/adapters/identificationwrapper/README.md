# Identification Wrapper

## Overview

The Identification Wrapper package provides a wrapper around the `github.com/abitofhelp/servicelib/valueobject/identification` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters), allowing the domain layer to remain isolated from external dependencies.

## Architecture

The Identification Wrapper package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/valueobject/identification` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

The package sits in the infrastructure layer of the application and is used by the domain layer through interfaces defined in the domain layer. The architecture follows these principles:

- **Dependency Inversion**: The domain layer depends on abstractions, not concrete implementations
- **Adapter Pattern**: This package adapts the external library to the domain's needs
- **Value Object Pattern**: The ID type is implemented as a value object with identity semantics

## Implementation Details

The Identification Wrapper package implements the following design patterns:

1. **Adapter Pattern**: Adapts the external library to the domain's needs
2. **Value Object Pattern**: The ID type is implemented as a value object with identity semantics
3. **Factory Pattern**: Factory methods create new ID instances

Key implementation details:

- **Type Alias**: The ID type is a type alias for string, making it type-safe
- **Factory Methods**: Methods like NewID() and NewIDFromString() create new ID instances
- **Value Object Methods**: Methods like Equals() and IsEmpty() provide value object semantics
- **JSON Marshaling**: The ID type implements the json.Marshaler and json.Unmarshaler interfaces

The package uses the `github.com/abitofhelp/servicelib/valueobject/identification` package internally but exposes its own API to the domain layer, ensuring that the domain layer doesn't directly depend on the external library.

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Family Service Example](../../../examples/family_service/README.md) - Shows how to use the identification wrapper

Example of using the identification wrapper:

```
// Create a new ID
id := identificationwrapper.NewID()

// Create an ID from a string
id2, err := identificationwrapper.NewIDFromString("123e4567-e89b-12d3-a456-426614174000")
if err != nil {
    // Handle error
}

// Get the string representation
str := id.String()

// Check if an ID is empty
if id.IsEmpty() {
    // Handle empty ID
}

// Compare IDs
if id.Equals(id2) {
    // IDs are equal
}
```

## Configuration

The Identification Wrapper package doesn't require any specific configuration. It uses the default configuration of the underlying `github.com/abitofhelp/servicelib/valueobject/identification` package.

## Testing

The Identification Wrapper package is tested through:

1. **Unit Tests**: Each function and method has unit tests
2. **Property-Based Testing**: Tests with randomized inputs to find edge cases
3. **Integration Tests**: Tests that verify the wrapper works correctly with the underlying library

Key testing approaches:

- **Factory Method Testing**: Tests that verify factory methods create valid IDs
- **Value Object Testing**: Tests that verify value object semantics (equality, emptiness)
- **JSON Marshaling Testing**: Tests that verify JSON marshaling and unmarshaling
- **Error Handling Testing**: Tests that verify error handling for invalid inputs

Example of a test case:

```
// Test function for NewIDFromString
// Valid UUID
validUUID := "123e4567-e89b-12d3-a456-426614174000"
id, err := identificationwrapper.NewIDFromString(validUUID)
assert.NoError(t, err)
assert.Equal(t, validUUID, id.String())

// Invalid UUID
invalidUUID := "not-a-uuid"
_, err = identificationwrapper.NewIDFromString(invalidUUID)
assert.Error(t, err)
```

## Design Notes

1. **Type Safety**: The ID type is a type alias for string, providing type safety
2. **Value Object Semantics**: The ID type has value object semantics (equality, immutability)
3. **Factory Methods**: Factory methods ensure that IDs are created correctly
4. **Error Handling**: Methods that can fail return errors that can be handled by the caller
5. **JSON Support**: The ID type implements the json.Marshaler and json.Unmarshaler interfaces for JSON support
6. **Dependency Inversion**: The package follows the Dependency Inversion Principle by ensuring that the domain layer depends on abstractions rather than concrete implementations

## API Documentation

### ID

The `ID` type represents a unique identifier:

```
type ID string
```

### Functions

The package provides the following functions:

- `NewID()`: Creates a new ID with a random UUID
- `NewIDFromString(id string)`: Creates a new ID from a string
- `(id ID) String()`: Returns the string representation of the ID
- `(id ID) IsEmpty()`: Checks if the ID is empty
- `(id ID) Equals(other ID)`: Checks if the ID equals another ID
- `(id ID) MarshalJSON()`: Implements the json.Marshaler interface
- `(id *ID) UnmarshalJSON(data []byte)`: Implements the json.Unmarshaler interface

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Value Objects](https://martinfowler.com/bliki/ValueObject.html)
- [UUID Specification](https://tools.ietf.org/html/rfc4122)
- [Domain Entities](../../../core/domain/entity/README.md) - Uses these IDs for entity identification
- [Domain Services](../../../core/domain/services/README.md) - Uses these IDs for service operations
- [Repository Wrapper](../repositorywrapper/README.md) - Uses these IDs for entity retrieval
