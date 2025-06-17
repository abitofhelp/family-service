// Copyright (c) 2025 A Bit of Help, Inc.

package integration

import (
	"context"
	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/services"
	"testing"
	"time"
)

type mockRepo struct {
	store map[string]*entity.Family
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	fam, ok := m.store[id]
	if !ok {
		return nil, nil // In a real implementation, this would return an error
	}
	return fam, nil
}

func (m *mockRepo) Save(ctx context.Context, f *entity.Family) error {
	m.store[f.ID()] = f
	return nil
}

func (m *mockRepo) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error) {
	// Implementation not needed for this test
	return nil, nil
}

func (m *mockRepo) FindByChildID(ctx context.Context, childID string) (*entity.Family, error) {
	// Implementation not needed for this test
	return nil, nil
}

func TestCreateFamily(t *testing.T) {
	repo := &mockRepo{store: make(map[string]*entity.Family)}
	svc := services.NewFamilyDomainService(repo)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	parentDTO := entity.ParentDTO{
		ID:        "p1",
		FirstName: "John",
		LastName:  "Doe",
		BirthDate: birthDate,
		DeathDate: nil,
	}

	// Create a family DTO
	dto := entity.FamilyDTO{
		ID:       "abc123",
		Status:   string(entity.Single),
		Parents:  []entity.ParentDTO{parentDTO},
		Children: []entity.ChildDTO{},
	}

	result, err := svc.CreateFamily(context.Background(), dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != dto.ID {
		t.Errorf("expected ID %s, got %s", dto.ID, result.ID)
	}
}
