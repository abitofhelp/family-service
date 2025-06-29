# Domain Metrics

## Overview

The Domain Metrics package provides Prometheus metrics for monitoring the Family Service application. It defines a comprehensive set of metrics for tracking family operations, family member counts, family status counts, family size distribution, parent and child age distributions, operation errors, business rule violations, API requests, and repository operations. These metrics provide valuable insights into the application's behavior and performance.

## Architecture

The Domain Metrics package is part of the core domain layer in the Clean Architecture and Hexagonal Architecture patterns. It sits at the center of the application and has no dependencies on other layers. The architecture follows these principles:

- **Domain-Driven Design (DDD)**: Metrics are designed based on the ubiquitous language of the domain
- **Clean Architecture**: Metrics are independent of infrastructure concerns
- **Hexagonal Architecture**: Metrics are used by both the domain and application layers
- **Dependency Inversion**: The package depends on abstractions (Prometheus interfaces) rather than concrete implementations

The package is organized into:

- **Metric Definitions**: Constants and variables that define the metrics
- **Registration Functions**: Functions that register metrics with the Prometheus registry
- **Helper Functions**: Functions that make it easier to use the metrics
- **Reset Functions**: Functions that reset metrics to their initial state (useful for testing)

## Implementation Details

The Domain Metrics package implements the following design patterns:

1. **Facade Pattern**: Provides a simplified interface to the Prometheus metrics library
2. **Singleton Pattern**: Ensures that metrics are registered only once
3. **Factory Pattern**: Creates and configures metrics with appropriate labels
4. **Observer Pattern**: Metrics observe and record system behavior

Key implementation details:

- **Prometheus Integration**: Uses the Prometheus client library for metric implementation
- **Label Pre-initialization**: Pre-initializes metric labels to avoid runtime overhead
- **Metric Types**: Uses appropriate metric types (Counter, Gauge, Histogram) for different measurements
- **Metric Naming**: Follows Prometheus naming conventions for metrics
- **Metric Documentation**: Includes help text for all metrics
- **Thread Safety**: All operations are thread-safe

## Features

- **Family Operation Metrics**: Track family operations like creation, addition of parents/children, and divorces
- **Family Member Metrics**: Monitor counts of parents and children
- **Family Status Metrics**: Track families by status (single, married, divorced, widowed)
- **Age Distribution Metrics**: Monitor age distributions of parents and children
- **Error Metrics**: Track operation errors and business rule violations
- **API Request Metrics**: Monitor API requests and their durations
- **Repository Operation Metrics**: Track repository operations and their durations
- **Pre-initialized Labels**: Optimize performance by pre-initializing metric labels
- **Metric Reset Capability**: Support for resetting metrics (useful for testing)

## API Documentation

### Core Metrics

#### Family Operations

Metrics for tracking family operations:

```
// Family operations counters
FamilyOperationsTotal = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "family_operations_total",
        Help: "Total number of family domain operations",
    },
    []string{"operation", "status"},
)

// Family operations duration
FamilyOperationsDuration = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "family_operations_duration_seconds",
        Help:    "Duration of family domain operations in seconds",
        Buckets: prometheus.DefBuckets,
    },
    []string{"operation"},
)
```

#### Family Member Counts

Metrics for tracking family member counts:

```
// Family member counts
FamilyMemberCounts = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
        Name: "family_member_counts",
        Help: "Current count of family members by type",
    },
    []string{"type"},
)

// Family status counts
FamilyStatusCounts = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
        Name: "family_status_counts",
        Help: "Current count of families by status",
    },
    []string{"status"},
)
```

### Key Functions

#### RegisterMetrics

Registers all metrics with the provided registry:

```
// RegisterMetrics registers all metrics with the provided registry
func RegisterMetrics(registry prometheus.Registerer)
```

#### ResetMetrics

Resets all metrics to their initial state:

```
// ResetMetrics resets all metrics to their initial state
// This is particularly useful for testing to ensure a clean state between tests
func ResetMetrics()
```

## Examples

There may be additional examples in the /EXAMPLES directory.

Example of using the metrics:

```
// Import the metrics package
import "github.com/abitofhelp/family-service/core/domain/metrics"

// Register metrics with the Prometheus registry
metrics.RegisterMetrics(prometheus.DefaultRegisterer)

// Record a family creation operation
metrics.FamilyOperationsTotal.WithLabelValues("create", "success").Inc()

// Record the duration of a family creation operation
timer := prometheus.NewTimer(metrics.FamilyOperationsDuration.WithLabelValues("create"))
defer timer.ObserveDuration()

// Update family member counts
metrics.FamilyMemberCounts.WithLabelValues("parent").Add(2)
metrics.FamilyMemberCounts.WithLabelValues("child").Add(1)

// Update family status counts
metrics.FamilyStatusCounts.WithLabelValues("married").Inc()
```

## Configuration

The Domain Metrics package doesn't require any specific configuration itself, but it integrates with Prometheus, which can be configured with various options:

- **Metric Registration**: Metrics must be registered with a Prometheus registry
- **Label Cardinality**: The cardinality of labels can be configured to balance detail and performance
- **Histogram Buckets**: Histogram buckets can be configured to capture the appropriate range of values
- **Metric Prefixes**: Metric names can be prefixed to provide additional context

Example configuration:

```
// Configure custom histogram buckets
customBuckets := []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5, 10}

// Create a custom histogram with the configured buckets
customHistogram := prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "custom_operation_duration_seconds",
        Help:    "Duration of custom operations in seconds",
        Buckets: customBuckets,
    },
    []string{"operation"},
)

// Register the custom histogram
prometheus.MustRegister(customHistogram)
```

## Testing

The Domain Metrics package is tested through:

1. **Unit Tests**: Each metric and function has unit tests
2. **Integration Tests**: Tests that verify metrics are correctly registered and updated
3. **Reset Testing**: Tests that verify metrics can be reset to their initial state

Key testing approaches:

- **Metric Registration Testing**: Tests that verify metrics are correctly registered
- **Metric Update Testing**: Tests that verify metrics are correctly updated
- **Metric Reset Testing**: Tests that verify metrics can be reset to their initial state
- **Label Testing**: Tests that verify labels are correctly applied

Example of a test case:

```
func TestFamilyOperationsTotal(t *testing.T) {
    // Reset metrics before the test
    metrics.ResetMetrics()

    // Record a family creation operation
    metrics.FamilyOperationsTotal.WithLabelValues("create", "success").Inc()

    // Verify the metric was updated
    metricFamily, err := metrics.FamilyOperationsTotal.MetricVec.Collect()
    assert.NoError(t, err)
    assert.Len(t, metricFamily, 1)

    // Verify the metric value
    metric := metricFamily[0].Metric[0]
    assert.Equal(t, float64(1), *metric.Counter.Value)

    // Verify the metric labels
    assert.Equal(t, "create", *metric.Label[0].Value)
    assert.Equal(t, "success", *metric.Label[1].Value)
}
```

## Design Notes

1. **Metric Types**: Different metric types (Counter, Gauge, Histogram) are used for different measurements
2. **Label Cardinality**: Labels are carefully chosen to balance detail and performance
3. **Metric Naming**: Metric names follow Prometheus naming conventions
4. **Metric Documentation**: All metrics include help text to explain their purpose
5. **Thread Safety**: All operations are thread-safe to support concurrent access
6. **Performance Optimization**: Label values are pre-initialized to avoid runtime overhead
7. **Reset Capability**: Metrics can be reset to their initial state, which is useful for testing

## Best Practices

1. **Consistent Naming**: Use consistent naming conventions for metrics
2. **Appropriate Labels**: Choose labels that provide meaningful dimensions for analysis
3. **Optimized Cardinality**: Avoid high cardinality labels that can cause performance issues
4. **Pre-initialized Labels**: Pre-initialize metric labels to avoid runtime overhead
5. **Documentation**: Document metrics clearly to ensure they are understood and used correctly

## References

- [Prometheus Documentation](https://prometheus.io/docs/introduction/overview/)
- [Prometheus Client Library for Go](https://github.com/prometheus/client_golang)
- [Prometheus Naming Conventions](https://prometheus.io/docs/practices/naming/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/instrumentation/)
- [Domain Entities](../entity/README.md) - The entities being monitored by these metrics
- [Domain Services](../services/README.md) - Services that use these metrics to track operations
- [Application Services](../../application/services/README.md) - Higher-level services that also use these metrics
- [Prometheus Integration](../../../infrastructure/adapters/telemetry/README.md) - Infrastructure for collecting and exposing these metrics
