// Copyright (c) 2025 A Bit of Help, Inc.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports/mock"
	"github.com/abitofhelp/family-service/infrastructure/adapters/loggingwrapper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestCreateFamily(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := loggingwrapper.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create test data
	dto := entity.FamilyDTO{
		ID:     "f47ac10b-58cc-4372-a567-0e02b2c3d479", // Valid UUID
		Status: "SINGLE",
		Parents: []entity.ParentDTO{
			{
				ID:        "38f5b8ed-1eb0-4a20-9f0e-7c3b3c3f3f3f", // Valid UUID
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	// Setup expectations
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	result, err := svc.CreateFamily(context.Background(), dto)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.Status, result.Status)
	assert.Equal(t, 1, result.ParentCount)
	assert.Equal(t, 0, result.ChildrenCount)
}

func TestGetFamily(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := loggingwrapper.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create test data
	familyID := "f47ac10b-58cc-4372-a567-0e02b2c3d479" // Valid UUID
	parent, _ := entity.NewParent("38f5b8ed-1eb0-4a20-9f0e-7c3b3c3f3f3f", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	family, _ := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent}, []*entity.Child{})

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), familyID).Return(family, nil)

	// Execute
	result, err := svc.GetFamily(context.Background(), familyID)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, familyID, result.ID)
}

func TestAddParent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := loggingwrapper.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create test data
	familyID := "f47ac10b-58cc-4372-a567-0e02b2c3d479" // Valid UUID
	parent1, _ := entity.NewParent("38f5b8ed-1eb0-4a20-9f0e-7c3b3c3f3f3f", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	family, _ := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent1}, []*entity.Child{})

	parentDTO := entity.ParentDTO{
		ID:        "a47ac10b-58cc-4372-a567-0e02b2c3d480", // Valid UUID
		FirstName: "Jane",
		LastName:  "Doe",
		BirthDate: time.Date(1982, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), familyID).Return(family, nil)
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	result, err := svc.AddParent(context.Background(), familyID, parentDTO)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, familyID, result.ID)
	assert.Equal(t, 2, result.ParentCount)
	assert.Equal(t, "MARRIED", result.Status) // Status should change to MARRIED when adding a second parent
}

func TestAddChild(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := loggingwrapper.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create test data
	familyID := "f47ac10b-58cc-4372-a567-0e02b2c3d479" // Valid UUID
	parent, _ := entity.NewParent("38f5b8ed-1eb0-4a20-9f0e-7c3b3c3f3f3f", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	family, _ := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent}, []*entity.Child{})

	childDTO := entity.ChildDTO{
		ID:        "b47ac10b-58cc-4372-a567-0e02b2c3d481", // Valid UUID
		FirstName: "Baby",
		LastName:  "Doe",
		BirthDate: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), familyID).Return(family, nil)
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	result, err := svc.AddChild(context.Background(), familyID, childDTO)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, familyID, result.ID)
	assert.Equal(t, 1, result.ParentCount)
	assert.Equal(t, 1, result.ChildrenCount)
}

func TestDivorce(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := loggingwrapper.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create test data
	familyID := "f47ac10b-58cc-4372-a567-0e02b2c3d479" // Valid UUID
	parent1, _ := entity.NewParent("38f5b8ed-1eb0-4a20-9f0e-7c3b3c3f3f3f", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	parent2, _ := entity.NewParent("a47ac10b-58cc-4372-a567-0e02b2c3d480", "Jane", "Doe", time.Date(1982, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	child, _ := entity.NewChild("b47ac10b-58cc-4372-a567-0e02b2c3d481", "Baby", "Doe", time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	family, _ := entity.NewFamily(familyID, entity.Married, []*entity.Parent{parent1, parent2}, []*entity.Child{child})

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), familyID).Return(family, nil)
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(2) // Save both families

	// Execute
	result, err := svc.Divorce(context.Background(), familyID, "38f5b8ed-1eb0-4a20-9f0e-7c3b3c3f3f3f")

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, familyID, result.ID)
	assert.Equal(t, 1, result.ParentCount)
	assert.Equal(t, 1, result.ChildrenCount)
	assert.Equal(t, "DIVORCED", result.Status)
}
