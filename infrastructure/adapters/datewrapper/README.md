# Date Wrapper

## Overview

The `datewrapper` package provides a wrapper around the `servicelib/date` package to avoid violating the dependency inversion principle. It offers utilities for working with dates and times in RFC3339 format, which is the standard format used throughout the application.

## Architecture

The `datewrapper` package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/date` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

## Implementation Details

The package implements simple wrapper functions around the `servicelib/date` package's functionality, ensuring consistent date formatting and parsing throughout the application. It handles both required and optional (nullable) dates.

Key features include:
- Parse date strings in RFC3339 format
- Format time.Time objects as RFC3339 strings
- Handle optional (nullable) dates
- Consistent date formatting across the application

Constants:
- `StandardDateFormat`: The standard date format used throughout the application (RFC3339)

Functions:
- `ParseDate`: Parses a date string in RFC3339 format
- `ParseOptionalDate`: Parses an optional date string in RFC3339 format
- `FormatDate`: Formats a time.Time as a string in RFC3339 format
- `FormatOptionalDate`: Formats an optional time.Time as a string in RFC3339 format

## Examples

```go
import "github.com/abitofhelp/family-service/infrastructure/adapters/datewrapper"

// Parse a date string
date, err := datewrapper.ParseDate("2023-04-15T14:30:00Z")
if err != nil {
    // Handle error
}

// Format a time.Time as a string
dateStr := datewrapper.FormatDate(time.Now())

// Parse an optional date string
var optionalDateStr *string
// optionalDateStr = &someString // or nil
optionalDate, err := datewrapper.ParseOptionalDate(optionalDateStr)
if err != nil {
    // Handle error
}

// Format an optional time.Time as a string
var optionalTime *time.Time
// optionalTime = &someTime // or nil
optionalTimeStr := datewrapper.FormatOptionalDate(optionalTime)
```

## Configuration

The `datewrapper` package doesn't require any specific configuration. It uses the RFC3339 format for all date operations, which is defined as a constant in the package.

## Testing

The package includes unit tests that verify the correct parsing and formatting of dates in RFC3339 format. Tests cover both required and optional date handling, including edge cases like nil pointers and invalid date strings.

## Design Notes

1. **RFC3339 Format**: The package standardizes on RFC3339 format for all date operations to ensure consistency across the application.
2. **Nullable Dates**: Special handling is provided for optional (nullable) dates to simplify working with dates that may or may not be present.
3. **Error Handling**: All parsing functions return appropriate errors when invalid date strings are provided.
4. **Immutability**: The package treats dates as immutable values, never modifying the input dates.

## References

- [RFC3339 Specification](https://tools.ietf.org/html/rfc3339)
- [Go time Package](https://golang.org/pkg/time/)
- [Domain Entities](../../../core/domain/entity/README.md) - Uses these date utilities for birth and death dates
- [Domain Services](../../../core/domain/services/README.md) - Uses these date utilities for business operations
