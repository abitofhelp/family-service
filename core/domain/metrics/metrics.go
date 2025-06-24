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

	// Family size distribution
	FamilySizeDistribution = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "family_size_distribution",
			Help:    "Distribution of family sizes (total number of members)",
			Buckets: []float64{1, 2, 3, 4, 5, 6, 8, 10, 15, 20},
		},
		[]string{"status"},
	)

	// Parent age distribution
	ParentAgeDistribution = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "parent_age_distribution_years",
			Help:    "Distribution of parent ages in years",
			Buckets: []float64{18, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 80, 90},
		},
		[]string{"status"},
	)

	// Child age distribution
	ChildAgeDistribution = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "child_age_distribution_years",
			Help:    "Distribution of child ages in years",
			Buckets: []float64{0, 1, 2, 5, 10, 15, 18, 21, 25},
		},
		[]string{},
	)

	// Family operation errors
	FamilyOperationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "family_operation_errors_total",
			Help: "Total number of errors in family operations by error type",
		},
		[]string{"operation", "error_type"},
	)

	// Business rule violations
	BusinessRuleViolations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_rule_violations_total",
			Help: "Total number of business rule violations by rule type",
		},
		[]string{"rule_type"},
	)

	// API request metrics
	APIRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of API requests by operation and status",
		},
		[]string{"operation", "status"},
	)

	// API request duration
	APIRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duration of API requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
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

	// Repository operation errors
	RepositoryOperationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repository_operation_errors_total",
			Help: "Total number of errors in repository operations by error type",
		},
		[]string{"operation", "error_type"},
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
		FamilySizeDistribution,
		ParentAgeDistribution,
		ChildAgeDistribution,
		FamilyOperationErrors,
		BusinessRuleViolations,
		APIRequestsTotal,
		APIRequestDuration,
		RepositoryOperationsTotal,
		RepositoryOperationsDuration,
		RepositoryOperationErrors,
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

	FamilySizeDistribution = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "family_size_distribution",
			Help:    "Distribution of family sizes (total number of members)",
			Buckets: []float64{1, 2, 3, 4, 5, 6, 8, 10, 15, 20},
		},
		[]string{"status"},
	)

	ParentAgeDistribution = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "parent_age_distribution_years",
			Help:    "Distribution of parent ages in years",
			Buckets: []float64{18, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 80, 90},
		},
		[]string{"status"},
	)

	ChildAgeDistribution = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "child_age_distribution_years",
			Help:    "Distribution of child ages in years",
			Buckets: []float64{0, 1, 2, 5, 10, 15, 18, 21, 25},
		},
		[]string{},
	)

	FamilyOperationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "family_operation_errors_total",
			Help: "Total number of errors in family operations by error type",
		},
		[]string{"operation", "error_type"},
	)

	BusinessRuleViolations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_rule_violations_total",
			Help: "Total number of business rule violations by rule type",
		},
		[]string{"rule_type"},
	)

	APIRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of API requests by operation and status",
		},
		[]string{"operation", "status"},
	)

	APIRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duration of API requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
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

	RepositoryOperationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repository_operation_errors_total",
			Help: "Total number of errors in repository operations by error type",
		},
		[]string{"operation", "error_type"},
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

	// Family size distribution
	FamilySizeDistribution.WithLabelValues("single")
	FamilySizeDistribution.WithLabelValues("married")
	FamilySizeDistribution.WithLabelValues("divorced")
	FamilySizeDistribution.WithLabelValues("widowed")

	// Parent age distribution
	ParentAgeDistribution.WithLabelValues("alive")
	ParentAgeDistribution.WithLabelValues("deceased")

	// Child age distribution - no labels needed

	// Family operation errors
	FamilyOperationErrors.WithLabelValues("create_family", "validation_error")
	FamilyOperationErrors.WithLabelValues("create_family", "database_error")
	FamilyOperationErrors.WithLabelValues("get_family", "not_found")
	FamilyOperationErrors.WithLabelValues("get_family", "database_error")
	FamilyOperationErrors.WithLabelValues("add_parent", "validation_error")
	FamilyOperationErrors.WithLabelValues("add_parent", "database_error")
	FamilyOperationErrors.WithLabelValues("add_child", "validation_error")
	FamilyOperationErrors.WithLabelValues("add_child", "database_error")
	FamilyOperationErrors.WithLabelValues("remove_child", "not_found")
	FamilyOperationErrors.WithLabelValues("remove_child", "database_error")
	FamilyOperationErrors.WithLabelValues("mark_parent_deceased", "not_found")
	FamilyOperationErrors.WithLabelValues("mark_parent_deceased", "database_error")
	FamilyOperationErrors.WithLabelValues("divorce", "validation_error")
	FamilyOperationErrors.WithLabelValues("divorce", "database_error")

	// Business rule violations
	BusinessRuleViolations.WithLabelValues("too_many_parents")
	BusinessRuleViolations.WithLabelValues("invalid_parent_age")
	BusinessRuleViolations.WithLabelValues("invalid_child_age")
	BusinessRuleViolations.WithLabelValues("invalid_family_status_transition")
	BusinessRuleViolations.WithLabelValues("duplicate_member")

	// API request metrics
	APIRequestsTotal.WithLabelValues("query_get_family", StatusSuccess)
	APIRequestsTotal.WithLabelValues("query_get_family", StatusFailure)
	APIRequestsTotal.WithLabelValues("query_get_all_families", StatusSuccess)
	APIRequestsTotal.WithLabelValues("query_get_all_families", StatusFailure)
	APIRequestsTotal.WithLabelValues("query_find_families_by_parent", StatusSuccess)
	APIRequestsTotal.WithLabelValues("query_find_families_by_parent", StatusFailure)
	APIRequestsTotal.WithLabelValues("query_find_family_by_child", StatusSuccess)
	APIRequestsTotal.WithLabelValues("query_find_family_by_child", StatusFailure)
	APIRequestsTotal.WithLabelValues("mutation_create_family", StatusSuccess)
	APIRequestsTotal.WithLabelValues("mutation_create_family", StatusFailure)
	APIRequestsTotal.WithLabelValues("mutation_add_parent", StatusSuccess)
	APIRequestsTotal.WithLabelValues("mutation_add_parent", StatusFailure)
	APIRequestsTotal.WithLabelValues("mutation_add_child", StatusSuccess)
	APIRequestsTotal.WithLabelValues("mutation_add_child", StatusFailure)
	APIRequestsTotal.WithLabelValues("mutation_remove_child", StatusSuccess)
	APIRequestsTotal.WithLabelValues("mutation_remove_child", StatusFailure)
	APIRequestsTotal.WithLabelValues("mutation_mark_parent_deceased", StatusSuccess)
	APIRequestsTotal.WithLabelValues("mutation_mark_parent_deceased", StatusFailure)
	APIRequestsTotal.WithLabelValues("mutation_divorce", StatusSuccess)
	APIRequestsTotal.WithLabelValues("mutation_divorce", StatusFailure)

	// API request duration
	APIRequestDuration.WithLabelValues("query_get_family")
	APIRequestDuration.WithLabelValues("query_get_all_families")
	APIRequestDuration.WithLabelValues("query_find_families_by_parent")
	APIRequestDuration.WithLabelValues("query_find_family_by_child")
	APIRequestDuration.WithLabelValues("mutation_create_family")
	APIRequestDuration.WithLabelValues("mutation_add_parent")
	APIRequestDuration.WithLabelValues("mutation_add_child")
	APIRequestDuration.WithLabelValues("mutation_remove_child")
	APIRequestDuration.WithLabelValues("mutation_mark_parent_deceased")
	APIRequestDuration.WithLabelValues("mutation_divorce")

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

	// Repository operation errors
	RepositoryOperationErrors.WithLabelValues("save", "connection_error")
	RepositoryOperationErrors.WithLabelValues("save", "constraint_violation")
	RepositoryOperationErrors.WithLabelValues("save", "serialization_error")
	RepositoryOperationErrors.WithLabelValues("get_by_id", "connection_error")
	RepositoryOperationErrors.WithLabelValues("get_by_id", "not_found")
	RepositoryOperationErrors.WithLabelValues("get_all", "connection_error")
	RepositoryOperationErrors.WithLabelValues("find_by_parent_id", "connection_error")
	RepositoryOperationErrors.WithLabelValues("find_by_child_id", "connection_error")
}

// Initialize metrics with the default registry
func init() {
	// Register metrics with the default Prometheus registry
	RegisterMetrics(prometheus.DefaultRegisterer)
}
