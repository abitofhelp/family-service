# Date Wrapper

## Overview

The `datewrapper` package provides a wrapper around the `servicelib/date` package to avoid violating the dependency inversion principle. It offers utilities for working with dates and times in RFC3339 format, which is the standard format used throughout the application.

## Features

- Parse date strings in RFC3339 format
- Format time.Time objects as RFC3339 strings
- Handle optional (nullable) dates
- Consistent date formatting across the application

## Usage

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

## Constants

- `StandardDateFormat`: The standard date format used throughout the application (RFC3339)

## Functions

- `ParseDate`: Parses a date string in RFC3339 format
- `ParseOptionalDate`: Parses an optional date string in RFC3339 format
- `FormatDate`: Formats a time.Time as a string in RFC3339 format
- `FormatOptionalDate`: Formats an optional time.Time as a string in RFC3339 format