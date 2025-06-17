// Copyright (c) 2025 A Bit of Help, Inc.

package entity

import (
	"testing"
	"time"
)

func TestFamilyValidationMaxParents(t *testing.T) {
	// Create three parents
	p1, _ := NewParent("p1", "John", "Doe", time.Now(), nil)
	p2, _ := NewParent("p2", "Jane", "Doe", time.Now(), nil)
	p3, _ := NewParent("p3", "Extra", "Parent", time.Now(), nil)

	// Try to create a family with too many parents
	_, err := NewFamily("test1", Married, []*Parent{p1, p2, p3}, []*Child{})

	// Should fail validation
	if err == nil {
		t.Fatal("expected validation error for too many parents")
	}
}
