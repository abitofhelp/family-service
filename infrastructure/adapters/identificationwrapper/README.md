# Identification Wrapper

This package provides a wrapper around the `github.com/abitofhelp/servicelib/valueobject/identification` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters).

## Purpose

The purpose of this wrapper is to:

1. Isolate the domain layer from external dependencies
2. Provide a consistent approach to identification throughout the application
3. Make it easier to replace or update the underlying identification library in the future

## Usage

Instead of directly importing `github.com/abitofhelp/servicelib/valueobject/identification`, import this wrapper:

```go
import "github.com/abitofhelp/family-service/infrastructure/adapters/identificationwrapper"
```

Then use the wrapper types and functions:

```go
// Create a new ID
id := identificationwrapper.NewID()

// Create an ID from a string
id, err := identificationwrapper.NewIDFromString("123e4567-e89b-12d3-a456-426614174000")

// Get the string representation
str := id.String()

// Check if an ID is empty
if id.IsEmpty() {
    // Handle empty ID
}

// Compare IDs
if id1.Equals(id2) {
    // IDs are equal
}
```

## Components

### ID

The `ID` type represents a unique identifier:

```go
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