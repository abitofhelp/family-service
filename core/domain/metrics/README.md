# Domain Metrics

## Overview

The Domain Metrics package provides Prometheus metrics for monitoring the Family Service application. It defines a comprehensive set of metrics for tracking family operations, family member counts, family status counts, family size distribution, parent and child age distributions, operation errors, business rule violations, API requests, and repository operations. These metrics provide valuable insights into the application's behavior and performance.

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

## Best Practices

1. **Consistent Naming**: Use consistent naming conventions for metrics
2. **Appropriate Labels**: Choose labels that provide meaningful dimensions for analysis
3. **Optimized Cardinality**: Avoid high cardinality labels that can cause performance issues
4. **Pre-initialized Labels**: Pre-initialize metric labels to avoid runtime overhead
5. **Documentation**: Document metrics clearly to ensure they are understood and used correctly

## Related Components

- [Domain Entities](../entity/README.md) - The entities being monitored by these metrics
- [Domain Services](../services/README.md) - Services that use these metrics to track operations
- [Application Services](../../application/services/README.md) - Higher-level services that also use these metrics
- [Prometheus Integration](../../../infrastructure/adapters/telemetry/README.md) - Infrastructure for collecting and exposing these metrics