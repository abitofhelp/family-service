// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/stretchr/testify/mock"
)

// MockFamilyMapper is a mock implementation of the FamilyMapper interface
type MockFamilyMapper struct {
	mock.Mock
}

func (m *MockFamilyMapper) ToDomain(input model.FamilyInput) (entity.FamilyDTO, error) {
	args := m.Called(input)
	if fn, ok := args.Get(0).(func(mock.Arguments) (entity.FamilyDTO, error)); ok {
		return fn(args)
	}
	return args.Get(0).(entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyMapper) ToGraphQL(dto entity.FamilyDTO) (*model.Family, error) {
	args := m.Called(dto)
	if fn, ok := args.Get(0).(func(mock.Arguments) (*model.Family, error)); ok {
		return fn(args)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Family), args.Error(1)
}

func (m *MockFamilyMapper) ToParentDTO(input model.ParentInput) (entity.ParentDTO, error) {
	args := m.Called(input)
	return args.Get(0).(entity.ParentDTO), args.Error(1)
}

func (m *MockFamilyMapper) ToChildDTO(input model.ChildInput) (entity.ChildDTO, error) {
	args := m.Called(input)
	return args.Get(0).(entity.ChildDTO), args.Error(1)
}

func (m *MockFamilyMapper) ToParent(dto entity.ParentDTO) (*model.Parent, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Parent), args.Error(1)
}

func (m *MockFamilyMapper) ToChild(dto entity.ChildDTO) (*model.Child, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Child), args.Error(1)
}

// NewMockFamilyMapper creates a new instance of MockFamilyMapper with default implementations
func NewMockFamilyMapper() *MockFamilyMapper {
	mapper := new(MockFamilyMapper)

	// Setup default implementations for common methods
	// For TestResolver_CreateFamily_Integration

	// We need to modify our approach to handle dynamic values
	// Instead of using a function, we'll set up the mock to match specific inputs

	// We'll set up specific mocks in each test instead of a default mock

	// We'll set up specific mocks in each test instead of using a matcher

	// Setup default implementation for ToDomain
	// We need to use a specific matcher for each test case
	// For TestResolver_CreateFamily_Integration

	// Set up a fixed DTO for testing
	mapper.On("ToDomain", mock.AnythingOfType("model.FamilyInput")).Return(
		entity.FamilyDTO{
			ID:     "family1",
			Status: "ACTIVE",
			Parents: []entity.ParentDTO{
				{
					ID:        "parent1",
					FirstName: "John",
					LastName:  "Doe",
					BirthDate: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			Children: []entity.ChildDTO{
				{
					ID:        "child1",
					FirstName: "Jane",
					LastName:  "Doe",
					BirthDate: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		}, nil)

	// Set up mock for ToChildDTO
	mapper.On("ToChildDTO", mock.AnythingOfType("model.ChildInput")).Return(
		entity.ChildDTO{
			ID:        "child2",
			FirstName: "Jim",
			LastName:  "Doe",
			BirthDate: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)

	return mapper
}
