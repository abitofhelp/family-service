# Rate Limiting Wrapper

## Overview

The Rate Limiting Wrapper package provides functionality for limiting the rate of requests to the application. It implements various rate limiting algorithms and strategies to protect the application from excessive load and potential denial-of-service attacks.

## Architecture

This package is part of the infrastructure layer in the Clean Architecture and Hexagonal Architecture patterns. It provides adapters for rate limiting that can be used by the application.

## Implementation Details

The Rate Limiting Wrapper implements the following design patterns:
- Decorator Pattern: Wraps existing functionality with rate limiting capabilities
- Strategy Pattern: Allows different rate limiting strategies to be used
- Factory Pattern: Creates instances of rate limiters

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:
- [Rate Limiting Example](../../../examples/rate_limiting/README.md) - Shows how to use the rate limiting wrapper

## Configuration

The Rate Limiting Wrapper can be configured with the following options:
- Rate Limit: Configure the maximum number of requests allowed per time window
- Time Window: Configure the time window for rate limiting
- Rate Limiting Strategy: Configure which rate limiting algorithm to use (token bucket, leaky bucket, fixed window, etc.)
- Burst Size: Configure the maximum burst size allowed

## Testing

The Rate Limiting Wrapper is tested through:
1. Unit Tests: Each rate limiting method has comprehensive unit tests
2. Integration Tests: Tests that verify the rate limiting wrapper works correctly with the application
3. Load Tests: Tests that verify the rate limiting wrapper can handle high load

## Design Notes

1. The Rate Limiting Wrapper supports multiple rate limiting algorithms
2. Rate limits can be configured per endpoint, per user, or globally
3. The wrapper provides feedback on rate limit status (remaining requests, reset time)
4. Rate limiting can be bypassed for certain users or endpoints

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Token Bucket Algorithm](https://en.wikipedia.org/wiki/Token_bucket)
- [Leaky Bucket Algorithm](https://en.wikipedia.org/wiki/Leaky_bucket)