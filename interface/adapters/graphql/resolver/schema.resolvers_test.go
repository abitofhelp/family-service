package resolver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/dto"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/abitofhelp/servicelib/valueobject/identification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFamilyService is a mock implementation of FamilyApplicationService
type MockFamilyService struct {
	mock.Mock
}

func (m *MockFamilyService) Create(ctx context.Context, dto *entity.FamilyDTO) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) GetByID(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) GetAll(ctx context.Context) ([]*entity.FamilyDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) GetFamily(ctx context.Context, id string) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) GetAllFamilies(ctx context.Context) ([]*entity.FamilyDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) UpdateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) DeleteFamily(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFamilyService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, familyID, parentDTO)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) AddChild(ctx context.Context, familyID string, childDTO entity.ChildDTO) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, familyID, childDTO)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) RemoveChild(ctx context.Context, familyID string, childID string) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, familyID, childID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) MarkParentDeceased(ctx context.Context, familyID string, parentID string, deathDate time.Time) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, familyID, parentID, deathDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, familyID, custodialParentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) FindFamiliesByParent(ctx context.Context, parentID string) ([]*entity.FamilyDTO, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.FamilyDTO), args.Error(1)
}

func (m *MockFamilyService) FindFamilyByChild(ctx context.Context, childID string) (*entity.FamilyDTO, error) {
	args := m.Called(ctx, childID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FamilyDTO), args.Error(1)
}

// Helper function to create a test family DTO
func createTestFamilyDTO() *entity.FamilyDTO {
	return &entity.FamilyDTO{
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
	}
}

func TestQueryResolver_GetFamily(t *testing.T) {
	// Create mock service
	mockService := new(MockFamilyService)
	resolver := NewResolver(mockService)

	// Create test data
	ctx := context.Background()
	familyID := identification.ID("family1")
	testFamily := createTestFamilyDTO()

	// Set up mock expectations
	mockService.On("GetFamily", ctx, familyID.String()).Return(testFamily, nil)

	// Execute the resolver
	result, err := resolver.Query().GetFamily(ctx, familyID)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, familyID, result.ID)
	assert.Equal(t, model.FamilyStatus("ACTIVE"), result.Status)
	assert.Len(t, result.Parents, 1)
	assert.Len(t, result.Children, 1)

	// Verify mock
	mockService.AssertExpectations(t)
}

func TestQueryResolver_GetFamily_Error(t *testing.T) {
	// Create mock service
	mockService := new(MockFamilyService)
	resolver := NewResolver(mockService)

	// Create test data
	ctx := context.Background()
	familyID := identification.ID("nonexistent")
	expectedErr := fmt.Errorf("family not found")

	// Set up mock expectations
	mockService.On("GetFamily", ctx, familyID.String()).Return(nil, expectedErr)

	// Execute the resolver
	result, err := resolver.Query().GetFamily(ctx, familyID)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), expectedErr.Error())

	// Verify mock
	mockService.AssertExpectations(t)
}

func TestMutationResolver_CreateFamily(t *testing.T) {
	// Create mock service
	mockService := new(MockFamilyService)
	resolver := NewResolver(mockService)

	// Create test data
	ctx := context.Background()
	input := model.FamilyInput{
		ID:     identification.ID("family1"),
		Status: model.FamilyStatusActive,
		Parents: []*model.ParentInput{
			{
				ID:        identification.ID("parent1"),
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: "1980-01-01",
			},
		},
		Children: []*model.ChildInput{
			{
				ID:        identification.ID("child1"),
				FirstName: "Jane",
				LastName:  "Doe",
				BirthDate: "2010-01-01",
			},
		},
	}

	testFamily := createTestFamilyDTO()

	// Convert input to domain DTO for mock expectation
	mapper := dto.NewFamilyMapper()
	expectedDTO, err := mapper.ToDomain(input)
	assert.NoError(t, err)

	// Set up mock expectations
	mockService.On("CreateFamily", ctx, expectedDTO).Return(testFamily, nil)

	// Execute the resolver
	result, err := resolver.Mutation().CreateFamily(ctx, input)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, input.ID, result.ID)
	assert.Equal(t, input.Status, result.Status)
	assert.Len(t, result.Parents, 1)
	assert.Len(t, result.Children, 1)

	// Verify mock
	mockService.AssertExpectations(t)
}

func TestMutationResolver_AddChild(t *testing.T) {
	// Create mock service
	mockService := new(MockFamilyService)
	resolver := NewResolver(mockService)

	// Create test data
	ctx := context.Background()
	familyID := identification.ID("family1")
	input := model.ChildInput{
		ID:        identification.ID("child2"),
		FirstName: "Jim",
		LastName:  "Doe",
		BirthDate: "2012-01-01",
	}

	testFamily := createTestFamilyDTO()
	testFamily.Children = append(testFamily.Children, entity.ChildDTO{
		ID:        input.ID.String(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		BirthDate: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	// Convert input to domain DTO for mock expectation
	mapper := dto.NewFamilyMapper()
	expectedDTO, err := mapper.ToChildDTO(input)
	assert.NoError(t, err)

	// Set up mock expectations
	mockService.On("AddChild", ctx, familyID.String(), expectedDTO).Return(testFamily, nil)

	// Execute the resolver
	result, err := resolver.Mutation().AddChild(ctx, familyID, input)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Children, 2)
	assert.Equal(t, input.ID, result.Children[1].ID)
	assert.Equal(t, input.FirstName, result.Children[1].FirstName)
	assert.Equal(t, input.LastName, result.Children[1].LastName)

	// Verify mock
	mockService.AssertExpectations(t)
}

func TestQueryResolver_CountChildren(t *testing.T) {
	// Create mock service
	mockService := new(MockFamilyService)
	resolver := NewResolver(mockService)

	// Create test data
	ctx := context.Background()
	testFamilies := []*entity.FamilyDTO{
		createTestFamilyDTO(),
		{
			ID:     "family2",
			Status: "ACTIVE",
			Children: []entity.ChildDTO{
				{
					ID:        "child2",
					FirstName: "Jim",
					LastName:  "Smith",
					BirthDate: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "child3",
					FirstName: "Sarah",
					LastName:  "Smith",
					BirthDate: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	// Set up mock expectations
	mockService.On("GetAllFamilies", ctx).Return(testFamilies, nil)

	// Execute the resolver
	count, err := resolver.Query().CountChildren(ctx)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, 3, count) // Total unique children across all families

	// Verify mock
	mockService.AssertExpectations(t)
}
