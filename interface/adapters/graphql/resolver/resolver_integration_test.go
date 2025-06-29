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

// MockFamilyService is defined in mock_family_service.go

func TestResolver_CreateFamily_Integration(t *testing.T) {
	// Setup mock service and mapper
	mockService := new(MockFamilyService)
	mockMapper := NewMockFamilyMapper()
	resolver := NewResolver(mockService, mockMapper)

	// Test data
	familyID := uuid.New().String()
	parentID := uuid.New().String()
	childID := uuid.New().String()
	now := time.Now()
	birthDate := now.AddDate(-30, 0, 0).Format(time.RFC3339)
	deathDate := now.AddDate(-1, 0, 0).Format(time.RFC3339)

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
	mockService.On("CreateFamily", mock.Anything, mock.AnythingOfType("entity.FamilyDTO")).Return(&expectedDTO, nil)

	// Setup mock expectations for the mapper
	// Create the expected GraphQL model that should be returned by the mapper
	expectedGraphQLFamily := &model.Family{
		ID:     identification.ID(familyID),
		Status: model.FamilyStatusActive,
		Parents: []*model.Parent{
			{
				ID:        identification.ID(parentID),
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: birthDate,
				DeathDate: &deathDate,
			},
		},
		Children: []*model.Child{
			{
				ID:        identification.ID(childID),
				FirstName: "Jane",
				LastName:  "Doe",
				BirthDate: birthDate,
			},
		},
	}

	// Set up the mock to return the expected GraphQL model when ToGraphQL is called with any DTO
	mockMapper.On("ToGraphQL", mock.AnythingOfType("entity.FamilyDTO")).Return(expectedGraphQLFamily, nil)

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
	// Setup mock service and mapper
	mockService := new(MockFamilyService)
	mockMapper := NewMockFamilyMapper()
	resolver := NewResolver(mockService, mockMapper)

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
	mockService.On("GetFamily", mock.Anything, familyID).Return(&familyDTO, nil)

	// Setup mock expectations for the mapper
	// Create the expected GraphQL model that should be returned by the mapper
	expectedGraphQLFamily := &model.Family{
		ID:     identification.ID(familyID),
		Status: model.FamilyStatusActive,
		Parents: []*model.Parent{
			{
				ID:        identification.ID(parentID),
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: birthDate.Format(time.RFC3339),
				DeathDate: func() *string { s := deathDate.Format(time.RFC3339); return &s }(),
			},
		},
		Children: []*model.Child{
			{
				ID:        identification.ID(childID),
				FirstName: "Jane",
				LastName:  "Doe",
				BirthDate: birthDate.Format(time.RFC3339),
			},
		},
	}

	// Set up the mock to return the expected GraphQL model when ToGraphQL is called with any DTO
	mockMapper.On("ToGraphQL", mock.AnythingOfType("entity.FamilyDTO")).Return(expectedGraphQLFamily, nil)

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
	assert.Equal(t, birthDate.Format(time.RFC3339), result.Parents[0].BirthDate)
	require.NotNil(t, result.Parents[0].DeathDate)
	assert.Equal(t, deathDate.Format(time.RFC3339), *result.Parents[0].DeathDate)

	// Assert children
	require.Len(t, result.Children, 1)
	assert.Equal(t, identification.ID(childID), result.Children[0].ID)
	assert.Equal(t, "Jane", result.Children[0].FirstName)
	assert.Equal(t, "Doe", result.Children[0].LastName)
	assert.Equal(t, birthDate.Format(time.RFC3339), result.Children[0].BirthDate)
	assert.Nil(t, result.Children[0].DeathDate)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestResolver_UpdateFamily_Integration(t *testing.T) {
	// Skip this test since UpdateFamily is not part of the MutationResolver interface
	t.Skip("UpdateFamily is not part of the MutationResolver interface")
}
