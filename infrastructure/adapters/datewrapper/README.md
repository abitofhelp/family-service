# Infrastructure Adapters - Date Wrapper

## Overview

The Date Wrapper adapter provides implementations for date and time-related ports defined in the core domain and application layers. This adapter connects the application to date and time libraries and utilities, following the Ports and Adapters (Hexagonal) architecture pattern. 

> **For Junior Developers**: Think of this adapter as a translator between your business logic and the date/time functionality. It allows your core business code to work with dates without knowing the specific date library being used.

By isolating date and time implementations in adapter classes, the core business logic remains independent of specific date/time technologies, making the system more maintainable, testable, and flexible.

## Features

- Date and time creation and manipulation
- Time zone handling
- Date formatting and parsing
- Duration calculations
- Date comparisons
- Clock abstraction for testability
- RFC3339 compliance
- Date range operations

## Getting Started

If you're new to this codebase, follow these steps to start using the Date Wrapper:

1. **Understand the purpose**: The Date Wrapper handles all date and time operations in a consistent way
2. **Learn the interfaces**: Look at the domain ports to understand what operations are available
3. **Ask questions**: If something isn't clear, ask a more experienced developer

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/datewrapper
```

## Configuration

The date wrapper can be configured according to specific requirements. Here's an example of configuring the date wrapper:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a date wrapper

// 1. Import necessary packages
import date, config, logging

// 2. Create a logger
// This is needed for the date wrapper to log any issues
logger = logging.NewLogger()

// 3. Configure the date wrapper
// These settings determine how dates are handled by default
dateConfig = {
    defaultTimeZone: "UTC",                        // Always use UTC for backend operations
    defaultDateFormat: "2006-01-02",               // Go's date format uses this specific date
    defaultTimeFormat: "15:04:05",                 // 24-hour time format
    defaultDateTimeFormat: "2006-01-02T15:04:05Z07:00"  // RFC3339 format
}

// 4. Create the date wrapper
// This is the object you'll use for all date operations
dateWrapper = date.NewDateWrapper(dateConfig, logger)

// 5. Use the date wrapper
// Examples of common operations:
now = dateWrapper.Now()                           // Get current time
formatted = dateWrapper.Format(now, "RFC3339")    // Format a date as string
parsed = dateWrapper.Parse("2023-06-15T14:30:00Z") // Parse a string to date
inTwoHours = dateWrapper.AddHours(now, 2)         // Add time to a date
isAfter = dateWrapper.IsAfter(inTwoHours, now)    // Compare dates
```

## API Documentation

### Core Concepts

> **For Junior Developers**: These concepts are fundamental to understanding how the Date Wrapper works. Take time to understand each one before diving into the code.

The date wrapper follows these core concepts:

1. **Adapter Pattern**: Implements date and time ports defined in the core domain or application layer
   - This means the Date Wrapper implements interfaces defined elsewhere
   - The business logic only knows about these interfaces, not the implementation details

2. **Dependency Injection**: Receives dependencies through constructor injection
   - Dependencies like loggers are passed in when creating the wrapper
   - This makes testing easier and components more loosely coupled

3. **Configuration**: Configured through a central configuration system
   - Settings like default time zones and formats are defined in configuration
   - This allows changing behavior without changing code

4. **Logging**: Uses a consistent logging approach
   - All operations are logged for debugging and monitoring
   - Context information is included in logs when available

5. **Error Handling**: Handles date and time errors gracefully
   - Invalid inputs are detected and reported clearly
   - Errors include helpful messages for debugging

### Key Adapter Functions

Here are the main functions you'll use when working with the Date Wrapper:

```
// Pseudocode example - not actual Go code
// This demonstrates a date wrapper implementation

// Date wrapper structure
type DateWrapper {
    config        // Date wrapper configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    clock         // Clock for testability (makes unit testing possible)
}

// Constructor for the date wrapper
// This is how you create a new instance of the wrapper
function NewDateWrapper(config, logger) {
    return new DateWrapper {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        clock: new RealClock()  // Uses real system time by default
    }
}

// Method to get the current time
// Use this when you need the current date/time
function DateWrapper.Now() {
    // Implementation would include:
    // 1. Using the clock to get the current time
    // 2. Applying the default time zone if needed
    // 3. Returning the time
}

// Method to format a date
// Use this to convert a date object to a string
function DateWrapper.Format(date, format) {
    // Implementation would include:
    // 1. Validating the date
    // 2. Resolving the format string (using predefined formats or the provided one)
    // 3. Formatting the date
    // 4. Handling formatting errors
    // 5. Returning the formatted string
}
```

### Common Date Operations

Here are some common operations you might need to perform:

1. **Getting the current time**:
   ```
   now = dateWrapper.Now()
   ```

2. **Formatting a date as a string**:
   ```
   formatted = dateWrapper.Format(date, "RFC3339")
   ```

3. **Parsing a string into a date**:
   ```
   date = dateWrapper.Parse("2023-06-15T14:30:00Z")
   ```

4. **Adding time to a date**:
   ```
   tomorrow = dateWrapper.AddDays(today, 1)
   ```

5. **Comparing dates**:
   ```
   if dateWrapper.IsAfter(date1, date2) {
       // date1 is after date2
   }
   ```

## Best Practices

> **For Junior Developers**: Following these best practices will help you avoid common pitfalls and write more maintainable code.

1. **Separation of Concerns**: Keep date and time logic separate from domain logic
   - **Why?** Your business logic shouldn't need to know how dates are formatted or parsed
   - **Example:** Don't put date formatting code in your domain entities

2. **Interface Segregation**: Define focused date and time interfaces in the domain layer
   - **Why?** Small, specific interfaces are easier to understand and implement
   - **Example:** Have separate interfaces for date creation, formatting, and comparison

3. **Dependency Injection**: Use constructor injection for adapter dependencies
   - **Why?** This makes testing easier and components more loosely coupled
   - **Example:** Pass the logger and configuration to the date wrapper constructor

4. **Error Handling**: Handle date and time errors gracefully
   - **Why?** Date operations can fail in many ways (invalid formats, out-of-range values)
   - **Example:** Check for errors when parsing dates from user input

5. **Consistent Formatting**: Use consistent date and time formats
   - **Why?** Consistency makes code more predictable and easier to maintain
   - **Example:** Always use RFC3339 for date-time strings in APIs

6. **Time Zone Awareness**: Be explicit about time zones
   - **Why?** Implicit time zones lead to bugs that are hard to track down
   - **Example:** Always specify the time zone when creating or formatting dates

7. **Testing**: Use a clock abstraction for testability
   - **Why?** This allows you to control time in your tests
   - **Example:** Inject a mock clock that returns a fixed time for deterministic tests

8. **RFC3339 Compliance**: Use RFC3339 format for date-time strings
   - **Why?** This is a standard format that works well across systems
   - **Example:** Use "2023-06-15T14:30:00Z" instead of custom formats

## Common Mistakes to Avoid

1. **Using local time for backend operations**
   - **Problem:** This leads to inconsistent behavior across different servers
   - **Solution:** Always use UTC for backend operations

2. **Hardcoding date formats**
   - **Problem:** Makes it hard to change formats later and may not work in all locales
   - **Solution:** Use configuration for date formats

3. **Not handling time zones properly**
   - **Problem:** Can lead to incorrect calculations and comparisons
   - **Solution:** Always be explicit about time zones

4. **Direct use of time libraries in domain code**
   - **Problem:** Creates tight coupling to specific libraries
   - **Solution:** Use the date wrapper adapter instead

## Troubleshooting

### Common Issues

#### Time Zone Issues

If you encounter time zone issues, consider the following:

- **Always be explicit about time zones**
  - **Example:** `dateWrapper.FormatWithTimeZone(date, "RFC3339", "UTC")`

- **Store dates in UTC and convert to local time zones only for display**
  - **Why?** UTC provides a consistent reference point
  - **Example:** `storedDate = dateWrapper.ToUTC(inputDate)`

- **Be aware of daylight saving time transitions**
  - **Problem:** Adding 24 hours might not give you the same time tomorrow during DST changes
  - **Solution:** Use `AddDays(1)` instead of `AddHours(24)`

- **Use time zone identifiers rather than offsets**
  - **Example:** Use "America/New_York" instead of "-05:00"
  - **Why?** Identifiers handle DST automatically

#### Parsing Errors

If you encounter date parsing errors, consider the following:

- **Validate input formats before parsing**
  - **Example:** Check if the string matches the expected pattern

- **Use clear error messages for parsing failures**
  - **Example:** "Date '2023-13-45' is invalid: month must be between 1 and 12"

- **Consider providing multiple format options for flexible parsing**
  - **Example:** Try parsing with several common formats before giving up

- **Be aware of locale-specific date formats**
  - **Problem:** "06/07/2023" could be June 7 or July 6 depending on locale
  - **Solution:** Use unambiguous formats like ISO 8601 (YYYY-MM-DD)

- **Use strict parsing mode when appropriate**
  - **Why?** Prevents silent acceptance of invalid dates
  - **Example:** `dateWrapper.ParseStrict("2023-06-15", "YYYY-MM-DD")`

## Related Components

> **For Junior Developers**: Understanding how components relate to each other is crucial for working effectively in this codebase.

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the date and time ports
  - This is where the interfaces that the Date Wrapper implements are defined
  - Look here to understand what operations are available

- [Application Layer](../../core/application/README.md) - The application layer that uses date and time operations
  - This layer contains the business logic that uses the Date Wrapper
  - See how date operations are used in business processes

- [Interface Adapters](../../interface/adapters/README.md) - The interface adapters that use date and time
  - These adapters handle external interactions and may format dates for display
  - Look here to see how dates are presented to users

## Glossary of Terms

- **Adapter Pattern**: A design pattern that allows incompatible interfaces to work together
- **Port**: An interface defined in the domain or application layer
- **Dependency Injection**: A technique where an object receives its dependencies from outside
- **RFC3339**: A standard format for representing dates and times (e.g., "2023-06-15T14:30:00Z")
- **UTC**: Coordinated Universal Time, the primary time standard used worldwide
- **Time Zone**: A region that observes a uniform standard time
- **DST**: Daylight Saving Time, the practice of advancing clocks during summer months

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
