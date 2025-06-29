// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/abitofhelp/servicelib/valueobject/identification"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockFamilyService is a mock implementation of the FamilyService interface
type MockFamilyService struct {
	mock.Mock
}

func (m *MockFamilyService) CreateFamily(ctx context.Context, family entity.FamilyDTO) (entity.FamilyDTO, error) {
	args := m.Called(ctx, family)
	return args.Get(0).(entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) GetFamily(ctx context.Context, id string) (entity.FamilyDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) UpdateFamily(ctx context.Context, family entity.FamilyDTO) (entity.FamilyDTO, error) {
	args := m.Called(ctx, family)
	return args.Get(0).(entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) DeleteFamily(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestResolver_CreateFamily_Integration(t *testing.T) {
	// Setup mock service
	mockService := new(MockFamilyService)
	resolver := NewResolver(mockService)

	// Test data
	familyID := uuid.New().String()
	parentID := uuid.New().String()
	childID := uuid.New().String()
	now := time.Now()
	birthDate := now.AddDate(-30, 0, 0).Format("2006-01-02")
	deathDate := now.AddDate(-1, 0, 0).Format("2006-01-02")

	input := model.FamilyInput{
		ID:     identification.ID(familyID),
		Status: model.FamilyStatusActive,
		Parents: []*model.ParentInput{
			{
				ID:        identification.ID(parentID),
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: birthDate,
				DeathDate: &deathDate,
			},
		},
		Children: []*model.ChildInput{
			{
				ID:        identification.ID(childID),
				FirstName: "Jane",
				LastName:  "Doe",
				BirthDate: birthDate,
			},
		},
	}

	// Expected domain DTO
	expectedDTO := entity.FamilyDTO{
		ID:     familyID,
		Status: "ACTIVE",
		Parents: []entity.ParentDTO{
			{
				ID:        parentID,
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: now.AddDate(-30, 0, 0),
				DeathDate: &now,
			},
		},
		Children: []entity.ChildDTO{
			{
				ID:        childID,
				FirstName: "Jane",
				LastName:  "Doe",
				BirthDate: now.AddDate(-30, 0, 0),
			},
		},
	}

	// Setup mock expectations
	mockService.On("CreateFamily", mock.Anything, mock.MatchedBy(func(dto entity.FamilyDTO) bool {
		return dto.ID == familyID &&
			dto.Status == "ACTIVE" &&
			len(dto.Parents) == 1 &&
			len(dto.Children) == 1 &&
			dto.Parents[0].ID == parentID &&
			dto.Children[0].ID == childID
	})).Return(expectedDTO, nil)

	// Execute test
	ctx := context.Background()
	result, err := resolver.Mutation().CreateFamily(ctx, input)

	// Assert results
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, identification.ID(familyID), result.ID)
	assert.Equal(t, model.FamilyStatusActive, result.Status)
	
	// Assert parents
	require.Len(t, result.Parents, 1)
	assert.Equal(t, identification.ID(parentID), result.Parents[0].ID)
	assert.Equal(t, "John", result.Parents[0].FirstName)
	assert.Equal(t, "Doe", result.Parents[0].LastName)
	assert.Equal(t, birthDate, result.Parents[0].BirthDate)
	require.NotNil(t, result.Parents[0].DeathDate)
	assert.Equal(t, deathDate, *result.Parents[0].DeathDate)

	// Assert children
	require.Len(t, result.Children, 1)
	assert.Equal(t, identification.ID(childID), result.Children[0].ID)
	assert.Equal(t, "Jane", result.Children[0].FirstName)
	assert.Equal(t, "Doe", result.Children[0].LastName)
	assert.Equal(t, birthDate, result.Children[0].BirthDate)
	assert.Nil(t, result.Children[0].DeathDate)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestResolver_GetFamily_Integration(t *testing.T) {
	// Setup mock service
	mockService := new(MockFamilyService)
	resolver := NewResolver(mockService)

	// Test data
	familyID := uuid.New().String()
	parentID := uuid.New().String()
	childID := uuid.New().String()
	now := time.Now()
	birthDate := now.AddDate(-30, 0, 0)
	deathDate := now.AddDate(-1, 0, 0)

	// Setup mock response
	familyDTO := entity.FamilyDTO{
		ID:     familyID,
		Status: "ACTIVE",
		Parents: []entity.ParentDTO{
			{
				ID:        parentID,
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: birthDate,
				DeathDate: &deathDate,
			},
		},
		Children: []entity.ChildDTO{
			{
				ID:        childID,
				FirstName: "Jane",
				LastName:  "Doe",
				BirthDate: birthDate,
			},
		},
	}

	// Setup mock expectations
	mockService.On("GetFamily", mock.Anything, familyID).Return(familyDTO, nil)

	// Execute test
	ctx := context.Background()
	result, err := resolver.Query().GetFamily(ctx, identification.ID(familyID))

	// Assert results
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, identification.ID(familyID), result.ID)
	assert.Equal(t, model.FamilyStatusActive, result.Status)

	// Assert parents
	require.Len(t, result.Parents, 1)
	assert.Equal(t, identification.ID(parentID), result.Parents[0].ID)
	assert.Equal(t, "John", result.Parents[0].FirstName)
	assert.Equal(t, "Doe", result.Parents[0].LastName)
	assert.Equal(t, birthDate.Format("2006-01-02"), result.Parents[0].BirthDate)
	require.NotNil(t, result.Parents[0].DeathDate)
	assert.Equal(t, deathDate.Format("2006-01-02"), *result.Parents[0].DeathDate)

	// Assert children
	require.Len(t, result.Children, 1)
	assert.Equal(t, identification.ID(childID), result.Children[0].ID)
	assert.Equal(t, "Jane", result.Children[0].FirstName)
	assert.Equal(t, "Doe", result.Children[0].LastName)
	assert.Equal(t, birthDate.Format("2006-01-02"), result.Children[0].BirthDate)
	assert.Nil(t, result.Children[0].DeathDate)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestResolver_UpdateFamily_Integration(t *testing.T) {
	// Setup mock service
	mockService := new(MockFamilyService)
	resolver := NewResolver(mockService)

	// Test data
	familyID := uuid.New().String()
	parentID := uuid.New().String()
	childID := uuid.New().String()
	now := time.Now()
	birthDate := now.AddDate(-30, 0, 0).Format("2006-01-02")
	deathDate := now.AddDate(-1, 0, 0).Format("2006-01-02")

	input := model.FamilyInput{
		ID:     identification.ID(familyID),
		Status: model.FamilyStatusDivorced, // Changed status
		Parents: []*model.ParentInput{
			{
				ID:        identification.ID(parentID),
				FirstName: "John",
				LastName:  "Smith", // Changed last name
				BirthDate: birthDate,
				DeathDate: &deathDate,
			},
		},
		Children: []*model.ChildInput{
			{
				ID:        identification.ID(childID),
				FirstName: "Jane",
				LastName:  "Smith", // Changed last name
				BirthDate: birthDate,
			},
		},
	}

	// Expected domain DTO
	expectedDTO := entity.FamilyDTO{
		ID:     familyID,
		Status: "DIVORCED",
		Parents: []entity.ParentDTO{
			{
				ID:        parentID,
				FirstName: "John",
				LastName:  "Smith",
				BirthDate: now.AddDate(-30, 0, 0),
				DeathDate: &now,
			},
		},
		Children: []entity.ChildDTO{
			{
				ID:        childID,
				FirstName: "Jane",
				LastName:  "Smith",
				BirthDate: now.AddDate(-30, 0, 0),
			},
		},
	}

	// Setup mock expectations
	mockService.On("UpdateFamily", mock.Anything, mock.MatchedBy(func(dto entity.FamilyDTO) bool {
		return dto.ID == familyID &&
			dto.Status == "DIVORCED" &&
			len(dto.Parents) == 1 &&
			len(dto.Children) == 1 &&
			dto.Parents[0].LastName == "Smith" &&
			dto.Children[0].LastName == "Smith"
	})).Return(expectedDTO, nil)

	// Execute test
	ctx := context.Background()
	result, err := resolver.Mutation().UpdateFamily(ctx, input)

	// Assert results
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, identification.ID(familyID), result.ID)
	assert.Equal(t, model.FamilyStatusDivorced, result.Status)

	// Assert parents
	require.Len(t, result.Parents, 1)
	assert.Equal(t, "Smith", result.Parents[0].LastName)

	// Assert children
	require.Len(t, result.Children, 1)
	assert.Equal(t, "Smith", result.Children[0].LastName)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}
