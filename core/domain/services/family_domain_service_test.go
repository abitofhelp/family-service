// Copyright (c) 2025 A Bit of Help, Inc.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports/mock"
	"github.com/abitofhelp/servicelib/logging"
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
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("00000000-0000-0000-0000-000000000001", "John", "Doe", birthDate, nil)
	require.NoError(t, err, "failed to create parent")

	// Create a family DTO
	dto := entity.FamilyDTO{
		ID:       "00000000-0000-0000-0000-000000000002",
		Status:   string(entity.Single),
		Parents:  []entity.ParentDTO{p.ToDTO()},
		Children: []entity.ChildDTO{},
	}

	// Set expectations
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, family *entity.Family) error {
			// Verify the family being saved
			assert.Equal(t, dto.ID, family.ID(), "Family ID should match")
			assert.Equal(t, entity.Status(dto.Status), family.Status(), "Family status should match")
			assert.Len(t, family.Parents(), 1, "Should have 1 parent")
			return nil
		})

	// Test
	result, err := svc.CreateFamily(context.Background(), dto)
	require.NoError(t, err, "unexpected error")

	// Verify
	assert.Equal(t, dto.ID, result.ID, "ID should match")
	assert.Equal(t, dto.Status, result.Status, "Status should match")
	assert.Len(t, result.Parents, 1, "Should have 1 parent")
}

func TestGetFamily(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("00000000-0000-0000-0000-000000000003", "John", "Doe", birthDate, nil)
	require.NoError(t, err, "failed to create parent")

	// Create a family
	fam, err := entity.NewFamily("00000000-0000-0000-0000-000000000004", entity.Single, []*entity.Parent{p}, []*entity.Child{})
	require.NoError(t, err, "failed to create family")

	// Set expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), "00000000-0000-0000-0000-000000000004").Return(fam, nil)

	// Test
	result, err := svc.GetFamily(context.Background(), "00000000-0000-0000-0000-000000000004")
	require.NoError(t, err, "unexpected error")

	// Verify
	assert.Equal(t, fam.ID(), result.ID, "ID should match")
	assert.Equal(t, string(fam.Status()), result.Status, "Status should match")
	assert.Len(t, result.Parents, 1, "Should have 1 parent")
}

func TestAddParent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create a parent
	birthDate1 := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p1, err := entity.NewParent("00000000-0000-0000-0000-000000000005", "John", "Doe", birthDate1, nil)
	require.NoError(t, err, "failed to create parent")

	// Create a family
	fam, err := entity.NewFamily("00000000-0000-0000-0000-000000000006", entity.Single, []*entity.Parent{p1}, []*entity.Child{})
	require.NoError(t, err, "failed to create family")

	// Create a second parent
	birthDate2 := time.Date(1982, 2, 2, 0, 0, 0, 0, time.UTC)
	p2DTO := entity.ParentDTO{
		ID:        "00000000-0000-0000-0000-000000000007",
		FirstName: "Jane",
		LastName:  "Doe",
		BirthDate: birthDate2,
		DeathDate: nil,
	}

	// Set expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), "00000000-0000-0000-0000-000000000006").Return(fam, nil)
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, family *entity.Family) error {
			// Verify the family being saved
			assert.Equal(t, "00000000-0000-0000-0000-000000000006", family.ID(), "Family ID should match")
			assert.Equal(t, entity.Married, family.Status(), "Family status should be Married")
			assert.Len(t, family.Parents(), 2, "Should have 2 parents")
			return nil
		})

	// Test
	result, err := svc.AddParent(context.Background(), "00000000-0000-0000-0000-000000000006", p2DTO)
	require.NoError(t, err, "unexpected error")

	// Verify
	assert.Len(t, result.Parents, 2, "Should have 2 parents")
	assert.Equal(t, string(entity.Married), result.Status, "Status should be Married")
}

func TestAddChild(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("00000000-0000-0000-0000-000000000008", "John", "Doe", birthDate, nil)
	require.NoError(t, err, "failed to create parent")

	// Create a family
	fam, err := entity.NewFamily("00000000-0000-0000-0000-000000000009", entity.Single, []*entity.Parent{p}, []*entity.Child{})
	require.NoError(t, err, "failed to create family")

	// Create a child
	childBirthDate := time.Date(2010, 3, 3, 0, 0, 0, 0, time.UTC)
	childDTO := entity.ChildDTO{
		ID:        "00000000-0000-0000-0000-00000000000a",
		FirstName: "Baby",
		LastName:  "Doe",
		BirthDate: childBirthDate,
		DeathDate: nil,
	}

	// Set expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), "00000000-0000-0000-0000-000000000009").Return(fam, nil)
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, family *entity.Family) error {
			// Verify the family being saved
			assert.Equal(t, "00000000-0000-0000-0000-000000000009", family.ID(), "Family ID should match")
			assert.Len(t, family.Children(), 1, "Should have 1 child")
			return nil
		})

	// Test
	result, err := svc.AddChild(context.Background(), "00000000-0000-0000-0000-000000000009", childDTO)
	require.NoError(t, err, "unexpected error")

	// Verify
	assert.Len(t, result.Children, 1, "Should have 1 child")
}

func TestDivorce(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFamilyRepository(ctrl)
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(mockRepo, contextLogger)

	// Create parents
	birthDate1 := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p1, err := entity.NewParent("00000000-0000-0000-0000-00000000000b", "John", "Doe", birthDate1, nil)
	require.NoError(t, err, "failed to create parent")

	birthDate2 := time.Date(1982, 2, 2, 0, 0, 0, 0, time.UTC)
	p2, err := entity.NewParent("00000000-0000-0000-0000-00000000000c", "Jane", "Doe", birthDate2, nil)
	require.NoError(t, err, "failed to create parent")

	// Create a child
	childBirthDate := time.Date(2010, 3, 3, 0, 0, 0, 0, time.UTC)
	c, err := entity.NewChild("00000000-0000-0000-0000-00000000000d", "Baby", "Doe", childBirthDate, nil)
	require.NoError(t, err, "failed to create child")

	// Create a family
	fam, err := entity.NewFamily("00000000-0000-0000-0000-00000000000e", entity.Married, []*entity.Parent{p1, p2}, []*entity.Child{c})
	require.NoError(t, err, "failed to create family")

	// Create a variable to capture the new family created during divorce
	var newFamily *entity.Family

	// Set expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), "00000000-0000-0000-0000-00000000000e").Return(fam, nil)

	// Expect two Save calls - one for the original family and one for the new family
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, family *entity.Family) error {
			// This is the original family being updated
			if family.ID() == "00000000-0000-0000-0000-00000000000e" {
				// Verify the family being saved
				assert.Equal(t, entity.Divorced, family.Status(), "Original family status should be Divorced")
				assert.Len(t, family.Parents(), 1, "Original family should have 1 parent")
				assert.Equal(t, "00000000-0000-0000-0000-00000000000b", family.Parents()[0].ID(), "Original family should have the custodial parent")
				assert.Len(t, family.Children(), 1, "Original family should have 1 child")
			} else {
				// This is the new family being created
				newFamily = family
				assert.Equal(t, entity.Divorced, family.Status(), "New family status should be Divorced")
				assert.Len(t, family.Parents(), 1, "New family should have 1 parent")
				assert.Equal(t, "00000000-0000-0000-0000-00000000000c", family.Parents()[0].ID(), "New family should have the remaining parent")
				assert.Len(t, family.Children(), 0, "New family should have 0 children")
			}
			return nil
		}).Times(2)

	// Test
	result, err := svc.Divorce(context.Background(), "00000000-0000-0000-0000-00000000000e", "00000000-0000-0000-0000-00000000000b")
	require.NoError(t, err, "unexpected error")

	// Verify the original family (now with custodial parent)
	assert.Equal(t, string(entity.Divorced), result.Status, "Status should be Divorced")
	assert.Len(t, result.Parents, 1, "Should have 1 parent")
	assert.Len(t, result.Children, 1, "Should have 1 child")
	assert.Equal(t, "00000000-0000-0000-0000-00000000000e", result.ID, "Family with custodial parent should keep the original ID")
	assert.Equal(t, "00000000-0000-0000-0000-00000000000b", result.Parents[0].ID, "Custodial parent ID should be correct")

	// Verify the new family was created
	require.NotNil(t, newFamily, "New family should have been created")
	assert.Equal(t, entity.Divorced, newFamily.Status(), "New family status should be Divorced")
	assert.Len(t, newFamily.Parents(), 1, "New family should have 1 parent")
	assert.Len(t, newFamily.Children(), 0, "New family should have 0 children")
	assert.Equal(t, "00000000-0000-0000-0000-00000000000c", newFamily.Parents()[0].ID(), "Remaining parent ID should be correct")
}
