// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

// MockFamilyService is a mock implementation of the FamilyApplicationService interface
type MockFamilyService struct {
	mock.Mock
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

// Additional methods that might be needed for tests
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

// Methods required by ApplicationService[*entity.Family, *entity.FamilyDTO] interface
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

// Method required by di.ApplicationService interface
func (m *MockFamilyService) GetID() string {
	args := m.Called()
	return args.String(0)
}
