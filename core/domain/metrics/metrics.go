// Copyright (c) 2025 A Bit of Help, Inc.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Domain operation metrics
var (
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

	// Repository operation metrics
	RepositoryOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repository_operations_total",
			Help: "Total number of repository operations",
		},
		[]string{"operation", "status"},
	)

	// Repository operations duration
	RepositoryOperationsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "repository_operations_duration_seconds",
			Help:    "Duration of repository operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// Operation status constants
const (
	StatusSuccess = "success"
	StatusFailure = "failure"
)

// RegisterMetrics registers all metrics with the provided registry
func RegisterMetrics(registry prometheus.Registerer) {
	// Register metrics with the provided registry
	registry.MustRegister(
		FamilyOperationsTotal,
		FamilyOperationsDuration,
		FamilyMemberCounts,
		FamilyStatusCounts,
		RepositoryOperationsTotal,
		RepositoryOperationsDuration,
	)

	// Pre-create metric labels to avoid runtime initialization
	initializeMetricLabels()
}

// ResetMetrics resets all metrics to their initial state
// This is particularly useful for testing to ensure a clean state between tests
func ResetMetrics() {
	// Reset all metrics by creating new instances
	FamilyOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "family_operations_total",
			Help: "Total number of family domain operations",
		},
		[]string{"operation", "status"},
	)

	FamilyOperationsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "family_operations_duration_seconds",
			Help:    "Duration of family domain operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	FamilyMemberCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "family_member_counts",
			Help: "Current count of family members by type",
		},
		[]string{"type"},
	)

	FamilyStatusCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "family_status_counts",
			Help: "Current count of families by status",
		},
		[]string{"status"},
	)

	RepositoryOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repository_operations_total",
			Help: "Total number of repository operations",
		},
		[]string{"operation", "status"},
	)

	RepositoryOperationsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "repository_operations_duration_seconds",
			Help:    "Duration of repository operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
}

// Initialize metric labels to avoid runtime initialization
func initializeMetricLabels() {
	// Family operations counters
	FamilyOperationsTotal.WithLabelValues("create_family", StatusSuccess)
	FamilyOperationsTotal.WithLabelValues("create_family", StatusFailure)
	FamilyOperationsTotal.WithLabelValues("get_family", StatusSuccess)
	FamilyOperationsTotal.WithLabelValues("get_family", StatusFailure)
	FamilyOperationsTotal.WithLabelValues("add_parent", StatusSuccess)
	FamilyOperationsTotal.WithLabelValues("add_parent", StatusFailure)
	FamilyOperationsTotal.WithLabelValues("add_child", StatusSuccess)
	FamilyOperationsTotal.WithLabelValues("add_child", StatusFailure)
	FamilyOperationsTotal.WithLabelValues("remove_child", StatusSuccess)
	FamilyOperationsTotal.WithLabelValues("remove_child", StatusFailure)
	FamilyOperationsTotal.WithLabelValues("mark_parent_deceased", StatusSuccess)
	FamilyOperationsTotal.WithLabelValues("mark_parent_deceased", StatusFailure)
	FamilyOperationsTotal.WithLabelValues("divorce", StatusSuccess)
	FamilyOperationsTotal.WithLabelValues("divorce", StatusFailure)

	// Family operations duration
	FamilyOperationsDuration.WithLabelValues("create_family")
	FamilyOperationsDuration.WithLabelValues("get_family")
	FamilyOperationsDuration.WithLabelValues("add_parent")
	FamilyOperationsDuration.WithLabelValues("add_child")
	FamilyOperationsDuration.WithLabelValues("remove_child")
	FamilyOperationsDuration.WithLabelValues("mark_parent_deceased")
	FamilyOperationsDuration.WithLabelValues("divorce")

	// Family member counts
	FamilyMemberCounts.WithLabelValues("parents")
	FamilyMemberCounts.WithLabelValues("children")

	// Family status counts
	FamilyStatusCounts.WithLabelValues("single")
	FamilyStatusCounts.WithLabelValues("married")
	FamilyStatusCounts.WithLabelValues("divorced")
	FamilyStatusCounts.WithLabelValues("widowed")

	// Repository operations
	RepositoryOperationsTotal.WithLabelValues("save", StatusSuccess)
	RepositoryOperationsTotal.WithLabelValues("save", StatusFailure)
	RepositoryOperationsTotal.WithLabelValues("get_by_id", StatusSuccess)
	RepositoryOperationsTotal.WithLabelValues("get_by_id", StatusFailure)
	RepositoryOperationsTotal.WithLabelValues("get_all", StatusSuccess)
	RepositoryOperationsTotal.WithLabelValues("get_all", StatusFailure)
	RepositoryOperationsTotal.WithLabelValues("find_by_parent_id", StatusSuccess)
	RepositoryOperationsTotal.WithLabelValues("find_by_parent_id", StatusFailure)
	RepositoryOperationsTotal.WithLabelValues("find_by_child_id", StatusSuccess)
	RepositoryOperationsTotal.WithLabelValues("find_by_child_id", StatusFailure)

	// Repository operations duration
	RepositoryOperationsDuration.WithLabelValues("save")
	RepositoryOperationsDuration.WithLabelValues("get_by_id")
	RepositoryOperationsDuration.WithLabelValues("get_all")
	RepositoryOperationsDuration.WithLabelValues("find_by_parent_id")
	RepositoryOperationsDuration.WithLabelValues("find_by_child_id")
}

// Initialize metrics with the default registry
func init() {
	// Register metrics with the default Prometheus registry
	RegisterMetrics(prometheus.DefaultRegisterer)
}
