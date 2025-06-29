# Validation Wrapper

## Overview

The Validation Wrapper package provides a wrapper around the `github.com/abitofhelp/servicelib/validation` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters), allowing the domain layer to remain isolated from external dependencies.

## Architecture

The Validation Wrapper package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/validation` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

The package sits in the infrastructure layer of the application and is used by the domain layer through interfaces defined in the domain layer. The architecture follows these principles:

- **Dependency Inversion**: The domain layer depends on abstractions, not concrete implementations
- **Adapter Pattern**: This package adapts the external library to the domain's needs
- **Composite Pattern**: Validation rules can be composed to create complex validation logic
- **Chain of Responsibility**: Validation rules are chained together in a pipeline

## Implementation Details

The Validation Wrapper package implements the following design patterns:

1. **Adapter Pattern**: Adapts the external library to the domain's needs
2. **Composite Pattern**: Validation rules can be composed to create complex validation logic
3. **Chain of Responsibility**: Validation rules are chained together in a pipeline
4. **Strategy Pattern**: Different validation strategies can be implemented as validation rules

Key implementation details:

- **Interface-Based Design**: The package defines interfaces for validation rules and pipelines
- **Composition Over Inheritance**: Validation rules can be composed to create complex validation logic
- **Fluent Interface**: The validation pipeline provides a fluent interface for adding rules
- **Helper Functions**: Common validation logic is encapsulated in helper functions

The package uses the `github.com/abitofhelp/servicelib/validation` package internally but exposes its own API to the domain layer, ensuring that the domain layer doesn't directly depend on the external library.

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Family Service Example](../../../examples/family_service/README.md) - Shows how to use the validation wrapper

Example of using the validation wrapper:

```
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

## Configuration

The Validation Wrapper package doesn't require any specific configuration. It provides a set of interfaces and implementations that can be used as-is. However, you can configure the validation behavior by:

- **Creating Custom Rules**: Implement the `ValidationRule` interface to create custom validation rules
- **Composing Rules**: Use the `CompositeRule` to compose multiple rules into a single rule
- **Building Pipelines**: Use the `ValidationPipeline` to build a pipeline of validation rules
- **Using Helper Functions**: Use the provided helper functions for common validations

## Testing

The Validation Wrapper package is tested through:

1. **Unit Tests**: Each function and interface implementation has unit tests
2. **Integration Tests**: Tests that verify the wrapper works correctly with the underlying library
3. **Validation Logic Tests**: Tests that verify the validation logic works correctly

Key testing approaches:

- **Rule Testing**: Tests that verify validation rules work correctly
- **Pipeline Testing**: Tests that verify validation pipelines work correctly
- **Composite Rule Testing**: Tests that verify composite rules work correctly
- **Helper Function Testing**: Tests that verify helper functions work correctly

Example of a test case:

```
// Test that the validation rule validates correctly
func TestValidationRule_Validate(t *testing.T) {
    // Create a validation rule
    rule := validationwrapper.NewValidationRule("myRule", func(entity interface{}) error {
        // Validation logic that always passes
        return nil
    })

    // Validate an entity
    err := rule.Validate("test entity")

    // Verify the validation passed
    assert.NoError(t, err)
}
```

## Design Notes

1. **Interface-Based Design**: The package uses interfaces to define validation rules and pipelines
2. **Composition Over Inheritance**: Validation rules can be composed to create complex validation logic
3. **Fluent Interface**: The validation pipeline provides a fluent interface for adding rules
4. **Helper Functions**: Common validation logic is encapsulated in helper functions
5. **Dependency Inversion**: The package follows the Dependency Inversion Principle by ensuring that the domain layer depends on abstractions rather than concrete implementations
6. **Error Handling**: Validation errors are returned as errors with descriptive messages

## API Documentation

### ValidationRule

The `ValidationRule` interface defines a validation rule that can be applied to an entity:

```
type ValidationRule interface {
    Validate(entity interface{}) error
}
```

### ValidationPipeline

The `ValidationPipeline` interface defines a pipeline of validation rules:

```
type ValidationPipeline interface {
    AddRule(rule ValidationRule) ValidationPipeline
    Validate(entity interface{}) error
}
```

### CompositeRule

The `CompositeRule` interface defines a composite validation rule that combines multiple rules:

```
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

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Composite Pattern](https://en.wikipedia.org/wiki/Composite_pattern)
- [Chain of Responsibility Pattern](https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern)
- [Domain Validation](../../../core/domain/validation/README.md) - Uses this wrapper for domain validation
- [Domain Entities](../../../core/domain/entity/README.md) - Entities that are validated using this wrapper
- [Error Wrapper](../errorswrapper/README.md) - Used for creating validation errors
