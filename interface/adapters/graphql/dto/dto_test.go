// Copyright (c) 2025 A Bit of Help, Inc.

package dto

import (
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/abitofhelp/servicelib/valueobject/identification"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFamilyMapper_ToDomain(t *testing.T) {
	// Setup test data
	familyID := uuid.New().String()
	parentID := uuid.New().String()
	childID := uuid.New().String()
	birthDate := "2000-01-01T00:00:00Z"
	deathDate := "2020-01-01T00:00:00Z"
	input := model.FamilyInput{
		ID:     identification.ID(familyID),
		Status: model.FamilyStatus("ACTIVE"),
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
				DeathDate: nil,
			},
		},
	}

	// Create mapper
	mapper := NewFamilyMapper()

	// Execute test
	result, err := mapper.ToDomain(input)

	// Assert results
	require.NoError(t, err)
	assert.Equal(t, input.ID.String(), result.ID)
	assert.Equal(t, "ACTIVE", result.Status)

	// Assert parents
	require.Len(t, result.Parents, 1)
	assert.Equal(t, input.Parents[0].ID.String(), result.Parents[0].ID)
	assert.Equal(t, input.Parents[0].FirstName, result.Parents[0].FirstName)
	assert.Equal(t, input.Parents[0].LastName, result.Parents[0].LastName)

	// Assert children
	require.Len(t, result.Children, 1)
	assert.Equal(t, input.Children[0].ID.String(), result.Children[0].ID)
	assert.Equal(t, input.Children[0].FirstName, result.Children[0].FirstName)
	assert.Equal(t, input.Children[0].LastName, result.Children[0].LastName)
}

func TestFamilyMapper_ToGraphQL(t *testing.T) {
	// Setup test data
	familyID := uuid.New().String()
	parentID := uuid.New().String()
	childID := uuid.New().String()
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	input := entity.FamilyDTO{
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
				DeathDate: nil,
			},
		},
	}

	// Create mapper
	mapper := NewFamilyMapper()

	// Execute test
	result, err := mapper.ToGraphQL(input)

	// Assert results
	require.NoError(t, err)
	assert.Equal(t, identification.ID(input.ID), result.ID)
	assert.Equal(t, model.FamilyStatus("ACTIVE"), result.Status)

	// Assert parents
	require.Len(t, result.Parents, 1)
	assert.Equal(t, identification.ID(input.Parents[0].ID), result.Parents[0].ID)
	assert.Equal(t, input.Parents[0].FirstName, result.Parents[0].FirstName)
	assert.Equal(t, input.Parents[0].LastName, result.Parents[0].LastName)
	assert.Equal(t, "2000-01-01T00:00:00Z", result.Parents[0].BirthDate)
	require.NotNil(t, result.Parents[0].DeathDate)
	assert.Equal(t, "2020-01-01T00:00:00Z", *result.Parents[0].DeathDate)

	// Assert children
	require.Len(t, result.Children, 1)
	assert.Equal(t, identification.ID(input.Children[0].ID), result.Children[0].ID)
	assert.Equal(t, input.Children[0].FirstName, result.Children[0].FirstName)
	assert.Equal(t, input.Children[0].LastName, result.Children[0].LastName)
	assert.Equal(t, "2000-01-01T00:00:00Z", result.Children[0].BirthDate)
	assert.Nil(t, result.Children[0].DeathDate)
}

func TestFamilyMapper_ToParentDTO(t *testing.T) {
	// Setup test data
	parentID := uuid.New().String()
	birthDate := "2000-01-01T00:00:00Z"
	deathDate := "2020-01-01T00:00:00Z"
	input := model.ParentInput{
		ID:        identification.ID(parentID),
		FirstName: "John",
		LastName:  "Doe",
		BirthDate: birthDate,
		DeathDate: &deathDate,
	}

	// Create mapper
	mapper := NewFamilyMapper()

	// Execute test
	result, err := mapper.ToParentDTO(input)

	// Assert results
	require.NoError(t, err)
	assert.Equal(t, input.ID.String(), result.ID)
	assert.Equal(t, input.FirstName, result.FirstName)
	assert.Equal(t, input.LastName, result.LastName)
	assert.Equal(t, birthDate, result.BirthDate.Format(time.RFC3339))
	require.NotNil(t, result.DeathDate)
	assert.Equal(t, deathDate, result.DeathDate.Format(time.RFC3339))
}

func TestFamilyMapper_ToChildDTO(t *testing.T) {
	// Setup test data
	childID := uuid.New().String()
	birthDate := "2000-01-01T00:00:00Z"
	input := model.ChildInput{
		ID:        identification.ID(childID),
		FirstName: "Jane",
		LastName:  "Doe",
		BirthDate: birthDate,
		DeathDate: nil,
	}

	// Create mapper
	mapper := NewFamilyMapper()

	// Execute test
	result, err := mapper.ToChildDTO(input)

	// Assert results
	require.NoError(t, err)
	assert.Equal(t, input.ID.String(), result.ID)
	assert.Equal(t, input.FirstName, result.FirstName)
	assert.Equal(t, input.LastName, result.LastName)
	assert.Equal(t, birthDate, result.BirthDate.Format(time.RFC3339))
	assert.Nil(t, result.DeathDate)
}

func TestFamilyMapper_ToParent(t *testing.T) {
	// Setup test data
	parentID := uuid.New().String()
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	input := entity.ParentDTO{
		ID:        parentID,
		FirstName: "John",
		LastName:  "Doe",
		BirthDate: birthDate,
		DeathDate: &deathDate,
	}

	// Create mapper
	mapper := NewFamilyMapper()

	// Execute test
	result, err := mapper.ToParent(input)

	// Assert results
	require.NoError(t, err)
	assert.Equal(t, identification.ID(input.ID), result.ID)
	assert.Equal(t, input.FirstName, result.FirstName)
	assert.Equal(t, input.LastName, result.LastName)
	assert.Equal(t, "2000-01-01T00:00:00Z", result.BirthDate)
	require.NotNil(t, result.DeathDate)
	assert.Equal(t, "2020-01-01T00:00:00Z", *result.DeathDate)
}

func TestFamilyMapper_ToChild(t *testing.T) {
	// Setup test data
	childID := uuid.New().String()
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	input := entity.ChildDTO{
		ID:        childID,
		FirstName: "Jane",
		LastName:  "Doe",
		BirthDate: birthDate,
		DeathDate: nil,
	}

	// Create mapper
	mapper := NewFamilyMapper()

	// Execute test
	result, err := mapper.ToChild(input)

	// Assert results
	require.NoError(t, err)
	assert.Equal(t, identification.ID(input.ID), result.ID)
	assert.Equal(t, input.FirstName, result.FirstName)
	assert.Equal(t, input.LastName, result.LastName)
	assert.Equal(t, "2000-01-01T00:00:00Z", result.BirthDate)
	assert.Nil(t, result.DeathDate)
}

func TestFamilyMapper_Error_Cases(t *testing.T) {
	tests := []struct {
		name          string
		setupInvalid  func() interface{}
		testFunction  func(interface{}) error
		expectedError string
	}{
		{
			name: "Invalid birth date in ParentInput",
			setupInvalid: func() interface{} {
				return model.ParentInput{
					ID:        identification.ID(uuid.New().String()),
					BirthDate: "invalid-date",
				}
			},
			testFunction: func(input interface{}) error {
				mapper := NewFamilyMapper()
				_, err := mapper.ToParentDTO(input.(model.ParentInput))
				return err
			},
			expectedError: "parsing time",
		},
		{
			name: "Invalid death date in ParentInput",
			setupInvalid: func() interface{} {
				invalidDate := "invalid-date"
				return model.ParentInput{
					ID:        identification.ID(uuid.New().String()),
					BirthDate: "2000-01-01T00:00:00Z",
					DeathDate: &invalidDate,
				}
			},
			testFunction: func(input interface{}) error {
				mapper := NewFamilyMapper()
				_, err := mapper.ToParentDTO(input.(model.ParentInput))
				return err
			},
			expectedError: "parsing time",
		},
		{
			name: "Invalid birth date in ChildInput",
			setupInvalid: func() interface{} {
				return model.ChildInput{
					ID:        identification.ID(uuid.New().String()),
					BirthDate: "invalid-date",
				}
			},
			testFunction: func(input interface{}) error {
				mapper := NewFamilyMapper()
				_, err := mapper.ToChildDTO(input.(model.ChildInput))
				return err
			},
			expectedError: "parsing time",
		},
		{
			name: "Invalid FamilyStatus in FamilyInput",
			setupInvalid: func() interface{} {
				return model.FamilyInput{
					ID:     identification.ID(uuid.New().String()),
					Status: model.FamilyStatus("INVALID_STATUS"),
				}
			},
			testFunction: func(input interface{}) error {
				mapper := NewFamilyMapper()
				_, err := mapper.ToDomain(input.(model.FamilyInput))
				return err
			},
			expectedError: "invalid family status",
		},
		{
			name: "Empty ID in ParentInput",
			setupInvalid: func() interface{} {
				return model.ParentInput{
					ID:        "",
					BirthDate: "2000-01-01T00:00:00Z",
				}
			},
			testFunction: func(input interface{}) error {
				mapper := NewFamilyMapper()
				_, err := mapper.ToParentDTO(input.(model.ParentInput))
				return err
			},
			expectedError: "invalid ID",
		},
		{
			name: "Death date before birth date in ParentInput",
			setupInvalid: func() interface{} {
				deathDate := "1999-01-01T00:00:00Z" // Before birth date
				return model.ParentInput{
					ID:        identification.ID(uuid.New().String()),
					BirthDate: "2000-01-01T00:00:00Z",
					DeathDate: &deathDate,
				}
			},
			testFunction: func(input interface{}) error {
				mapper := NewFamilyMapper()
				_, err := mapper.ToParentDTO(input.(model.ParentInput))
				return err
			},
			expectedError: "death date cannot be before birth date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.setupInvalid()
			err := tt.testFunction(input)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestFamilyMapper_Edge_Cases(t *testing.T) {
	t.Run("Empty parents and children lists", func(t *testing.T) {
		input := entity.FamilyDTO{
			ID:       uuid.New().String(),
			Status:   "ACTIVE",
			Parents:  []entity.ParentDTO{},
			Children: []entity.ChildDTO{},
		}

		mapper := NewFamilyMapper()
		result, err := mapper.ToGraphQL(input)

		require.NoError(t, err)
		assert.Equal(t, identification.ID(input.ID), result.ID)
		assert.Equal(t, model.FamilyStatus("ACTIVE"), result.Status)
		assert.Empty(t, result.Parents)
		assert.Empty(t, result.Children)
	})

	t.Run("Future dates are allowed", func(t *testing.T) {
		futureDate := time.Now().AddDate(1, 0, 0).Format(RFC3339DateFormat)
		input := model.ParentInput{
			ID:        identification.ID(uuid.New().String()),
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: futureDate,
		}

		mapper := NewFamilyMapper()
		result, err := mapper.ToParentDTO(input)

		require.NoError(t, err)
		assert.Equal(t, input.ID.String(), result.ID)
		assert.Equal(t, futureDate, result.BirthDate.Format(RFC3339DateFormat))
	})

	t.Run("Very long names are allowed", func(t *testing.T) {
		longName := "ThisIsAVeryLongNameThatShouldStillBeAllowedInTheSystem"
		input := model.ChildInput{
			ID:        identification.ID(uuid.New().String()),
			FirstName: longName,
			LastName:  longName,
			BirthDate: "2000-01-01T00:00:00Z",
		}

		mapper := NewFamilyMapper()
		result, err := mapper.ToChildDTO(input)

		require.NoError(t, err)
		assert.Equal(t, input.FirstName, result.FirstName)
		assert.Equal(t, input.LastName, result.LastName)
	})
}
