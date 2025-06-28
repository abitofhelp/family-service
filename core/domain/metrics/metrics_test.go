// Copyright (c) 2025 A Bit of Help, Inc.

package metrics_test

import (
	"testing"

	"github.com/abitofhelp/family-service/core/domain/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetricsInitialization(t *testing.T) {
	// Test that metrics are properly initialized
	assert.NotNil(t, metrics.FamilyOperationsTotal)
	assert.NotNil(t, metrics.FamilyOperationsDuration)
	assert.NotNil(t, metrics.FamilyMemberCounts)
	assert.NotNil(t, metrics.FamilyStatusCounts)
	assert.NotNil(t, metrics.RepositoryOperationsTotal)
	assert.NotNil(t, metrics.RepositoryOperationsDuration)
}

func TestMetricsIncrement(t *testing.T) {
	// Create a new registry for this test
	registry := prometheus.NewRegistry()

	// Create a new counter for this test
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_counter",
			Help: "Test counter for unit tests",
		},
		[]string{"operation", "status"},
	)

	// Register the counter with the test registry
	registry.MustRegister(counter)

	// Increment the counter
	counter.WithLabelValues("create_family", "success").Inc()

	// Check that the counter was incremented
	count, err := testutil.GatherAndCount(
		registry, 
		"test_counter",
		"operation=\"create_family\",status=\"success\"",
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Expected counter to be incremented")
}

func TestMetricsObserve(t *testing.T) {
	// Create a new registry for this test
	registry := prometheus.NewRegistry()

	// Create a new histogram for this test
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "test_histogram",
			Help:    "Test histogram for unit tests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// Register the histogram with the test registry
	registry.MustRegister(histogram)

	// Observe a duration
	histogram.WithLabelValues("create_family").Observe(0.5)

	// Check that the histogram was updated
	count, err := testutil.GatherAndCount(
		registry, 
		"test_histogram",
		"operation=\"create_family\"",
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Expected histogram to be updated")
}

func TestMetricsGauge(t *testing.T) {
	// Create a new registry for this test
	registry := prometheus.NewRegistry()

	// Create a new gauge for this test
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "test_gauge",
			Help: "Test gauge for unit tests",
		},
		[]string{"type"},
	)

	// Register the gauge with the test registry
	registry.MustRegister(gauge)

	// Set a gauge value
	gauge.WithLabelValues("parents").Set(2)

	// Check that the gauge was set
	count, err := testutil.GatherAndCount(
		registry, 
		"test_gauge",
		"type=\"parents\"",
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Expected gauge to be set")
}
