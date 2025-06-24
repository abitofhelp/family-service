// Copyright (c) 2025 A Bit of Help, Inc.

package entity

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// generateTestUUID generates a UUID for testing
func generateTestUUID() string {
	return uuid.New().String()
}

func TestFamilyValidationMaxParents(t *testing.T) {
	// Create three parents
	p1, err := NewParent(generateTestUUID(), "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("Failed to create parent p1: %v", err)
	}

	p2, err := NewParent(generateTestUUID(), "Jane", "Doe", time.Date(1982, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("Failed to create parent p2: %v", err)
	}

	p3, err := NewParent(generateTestUUID(), "Extra", "Parent", time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("Failed to create parent p3: %v", err)
	}

	// Try to create a family with too many parents
	_, err = NewFamily(generateTestUUID(), Married, []*Parent{p1, p2, p3}, []*Child{})

	// Should fail validation
	if err == nil {
		t.Fatal("expected validation error for too many parents")
	}

	t.Logf("Error message: %s", err.Error())
}

func TestParentValidationMinimumAge(t *testing.T) {
	// Create a parent that is too young (17 years old)
	now := time.Now()
	youngParentBirthDate := now.AddDate(-17, 0, 0) // 17 years ago
	_, err := NewParent(generateTestUUID(), "Young", "Parent", youngParentBirthDate, nil)

	// Should fail validation
	assert.NotNil(t, err, "expected validation error for parent age")
	t.Logf("Error message: %v", err)
	t.Logf("Error type: %T", err)

	// Just check if the error message contains the expected text
	assert.True(t, strings.Contains(err.Error(), "Validation failed"), 
		"error should mention validation failure")
}

func TestFamilyValidationParentAge(t *testing.T) {
	// Skip this test for now as it's covered by TestParentValidationMinimumAge
	t.Skip("This test is covered by TestParentValidationMinimumAge")
}

func TestFamilyValidationChildBirthDate(t *testing.T) {
	// Create a parent
	parentBirthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	parent, err := NewParent(generateTestUUID(), "John", "Doe", parentBirthDate, nil)
	if err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}

	// Create a child with birth date before parent's birth date
	childBirthDate := time.Date(1979, 1, 1, 0, 0, 0, 0, time.UTC)
	child, err := NewChild(generateTestUUID(), "Baby", "Doe", childBirthDate, nil)
	if err != nil {
		t.Fatalf("Failed to create child: %v", err)
	}

	// Try to create a family with a child born before the parent
	_, err = NewFamily(generateTestUUID(), Single, []*Parent{parent}, []*Child{child})

	// Should fail validation
	assert.NotNil(t, err, "expected validation error for child birth date")
	t.Logf("Error message: %s", err.Error())
	assert.True(t, strings.Contains(err.Error(), "Validation failed"), 
		"error should mention validation failure")
}

func TestFamilyValidationParentChildAgeGap(t *testing.T) {
	// Create a parent
	parentBirthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	parent, err := NewParent(generateTestUUID(), "Young", "Parent", parentBirthDate, nil)
	if err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}

	// Create a child with too small age gap (only 10 years)
	childBirthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	child, err := NewChild(generateTestUUID(), "Baby", "Child", childBirthDate, nil)
	if err != nil {
		t.Fatalf("Failed to create child: %v", err)
	}

	// Try to create a family with too small parent-child age gap
	_, err = NewFamily(generateTestUUID(), Single, []*Parent{parent}, []*Child{child})

	// Should fail validation
	assert.NotNil(t, err, "expected validation error for parent-child age gap")
	t.Logf("Error message: %s", err.Error())
	assert.True(t, strings.Contains(err.Error(), "Validation failed"), 
		"error should mention validation failure")
}

func TestFamilyValidationAbandonedStatus(t *testing.T) {
	// Create a parent
	parent, err := NewParent(generateTestUUID(), "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}

	// Try to create an abandoned family with no children
	_, err = NewFamily(generateTestUUID(), Abandoned, []*Parent{parent}, []*Child{})

	// Should fail validation
	assert.NotNil(t, err, "expected validation error for abandoned family with no children")
	t.Logf("Error message: %s", err.Error())
	assert.True(t, strings.Contains(err.Error(), "Validation failed"), 
		"error should mention validation failure")
}

func TestFamilyValidationWidowedStatus(t *testing.T) {
	// Create a deceased parent
	parentBirthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	parentDeathDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	parent, err := NewParent(generateTestUUID(), "John", "Doe", parentBirthDate, &parentDeathDate)
	if err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}

	// Try to create a widowed family with a deceased parent
	_, err = NewFamily(generateTestUUID(), Widowed, []*Parent{parent}, []*Child{})

	// Should fail validation
	assert.NotNil(t, err, "expected validation error for widowed family with deceased parent")
	t.Logf("Error message: %s", err.Error())
	assert.True(t, strings.Contains(err.Error(), "Validation failed"), 
		"error should mention validation failure")
}
