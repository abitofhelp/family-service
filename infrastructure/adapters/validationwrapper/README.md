# Validation Wrapper

This package provides a wrapper around the `github.com/abitofhelp/servicelib/validation` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters).

## Purpose

The purpose of this wrapper is to:

1. Isolate the domain layer from external dependencies
2. Provide a consistent validation approach throughout the application
3. Make it easier to replace or update the underlying validation library in the future

## Usage

Instead of directly importing `github.com/abitofhelp/servicelib/validation`, import this wrapper:

```go
import "github.com/abitofhelp/family-service/infrastructure/adapters/validationwrapper"
```

Then use the wrapper functions and interfaces:

```go
// Create a validation rule
rule := validationwrapper.NewValidationRule("myRule", func(entity interface{}) error {
    // Validation logic
    return nil
})

// Create a validation pipeline
pipeline := validationwrapper.NewValidationPipeline()
pipeline.AddRule(rule)

// Validate an entity
err := pipeline.Validate(myEntity)

// Create a composite rule
compositeRule := validationwrapper.NewCompositeRule("myCompositeRule")
compositeRule.AddRule(rule1)
compositeRule.AddRule(rule2)

// Use helper functions
err := validationwrapper.ValidateNotNil(value, "fieldName")
err := validationwrapper.ValidateNotEmpty(value, "fieldName")
err := validationwrapper.ValidateMinLength(value, 5, "fieldName")
err := validationwrapper.ValidateMaxLength(value, 10, "fieldName")
```

## Components

### ValidationRule

The `ValidationRule` interface defines a validation rule that can be applied to an entity:

```go
type ValidationRule interface {
    Validate(entity interface{}) error
}
```

### ValidationPipeline

The `ValidationPipeline` interface defines a pipeline of validation rules:

```go
type ValidationPipeline interface {
    AddRule(rule ValidationRule) ValidationPipeline
    Validate(entity interface{}) error
}
```

### CompositeRule

The `CompositeRule` interface defines a composite validation rule that combines multiple rules:

```go
type CompositeRule interface {
    ValidationRule
    AddRule(rule ValidationRule) CompositeRule
}
```

### Helper Functions

The package provides several helper functions for common validations:

- `ValidateNotNil`: Validates that a value is not nil
- `ValidateNotEmpty`: Validates that a string is not empty
- `ValidateMinLength`: Validates that a string has a minimum length
- `ValidateMaxLength`: Validates that a string has a maximum length