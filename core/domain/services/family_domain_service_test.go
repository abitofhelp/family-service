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
		ID:     "test-family-id",
		Status: "SINGLE",
		Parents: []entity.ParentDTO{
			{
				ID:        "test-parent-id",
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
	familyID := "test-family-id"
	family, _ := entity.NewFamily(familyID, entity.Single, []*entity.Parent{}, []*entity.Child{})

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
	familyID := "test-family-id"
	parent1, _ := entity.NewParent("parent1", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	family, _ := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent1}, []*entity.Child{})

	parentDTO := entity.ParentDTO{
		ID:        "parent2",
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
	familyID := "test-family-id"
	parent, _ := entity.NewParent("parent1", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	family, _ := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent}, []*entity.Child{})

	childDTO := entity.ChildDTO{
		ID:        "child1",
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
	familyID := "test-family-id"
	parent1, _ := entity.NewParent("parent1", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	parent2, _ := entity.NewParent("parent2", "Jane", "Doe", time.Date(1982, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	child, _ := entity.NewChild("child1", "Baby", "Doe", time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	family, _ := entity.NewFamily(familyID, entity.Married, []*entity.Parent{parent1, parent2}, []*entity.Child{child})

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), familyID).Return(family, nil)
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(2) // Save both families

	// Execute
	result, err := svc.Divorce(context.Background(), familyID, "parent1")

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, familyID, result.ID)
	assert.Equal(t, 1, result.ParentCount)
	assert.Equal(t, 1, result.ChildrenCount)
	assert.Equal(t, "DIVORCED", result.Status)
}