// Copyright (c) 2025 A Bit of Help, Inc.

package services

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap/zaptest"
)

// mockRepo is a mock implementation of the FamilyRepository interface for testing
type mockRepo struct {
	store map[string]*entity.Family
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	fam, ok := m.store[id]
	if !ok {
		return nil, errors.NewNotFoundError("Family", id, nil)
	}
	return fam, nil
}

func (m *mockRepo) Save(ctx context.Context, f *entity.Family) error {
	if f == nil {
		return errors.NewValidationError("family cannot be nil", "family", nil)
	}
	m.store[f.ID()] = f
	return nil
}

func (m *mockRepo) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error) {
	var result []*entity.Family
	for _, fam := range m.store {
		for _, p := range fam.Parents() {
			if p.ID() == parentID {
				result = append(result, fam)
				break
			}
		}
	}
	return result, nil
}

func (m *mockRepo) FindByChildID(ctx context.Context, childID string) (*entity.Family, error) {
	for _, fam := range m.store {
		for _, c := range fam.Children() {
			if c.ID() == childID {
				return fam, nil
			}
		}
	}
	return nil, errors.NewNotFoundError("Family with Child", childID, nil)
}

func (m *mockRepo) GetAll(ctx context.Context) ([]*entity.Family, error) {
	var families []*entity.Family
	for _, fam := range m.store {
		families = append(families, fam)
	}
	return families, nil
}

func TestCreateFamily(t *testing.T) {
	// Setup
	repo := &mockRepo{store: make(map[string]*entity.Family)}
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(repo, contextLogger)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("00000000-0000-0000-0000-000000000001", "John", "Doe", birthDate, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a family DTO
	dto := entity.FamilyDTO{
		ID:       "00000000-0000-0000-0000-000000000002",
		Status:   string(entity.Single),
		Parents:  []entity.ParentDTO{p.ToDTO()},
		Children: []entity.ChildDTO{},
	}

	// Test
	result, err := svc.CreateFamily(context.Background(), dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify
	if result.ID != dto.ID {
		t.Errorf("expected ID %s, got %s", dto.ID, result.ID)
	}
	if result.Status != dto.Status {
		t.Errorf("expected Status %s, got %s", dto.Status, result.Status)
	}
	if len(result.Parents) != 1 {
		t.Errorf("expected 1 parent, got %d", len(result.Parents))
	}
}

func TestGetFamily(t *testing.T) {
	// Setup
	repo := &mockRepo{store: make(map[string]*entity.Family)}
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(repo, contextLogger)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("00000000-0000-0000-0000-000000000003", "John", "Doe", birthDate, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a family
	fam, err := entity.NewFamily("00000000-0000-0000-0000-000000000004", entity.Single, []*entity.Parent{p}, []*entity.Child{})
	if err != nil {
		t.Fatalf("failed to create family: %v", err)
	}

	// Save the family
	err = repo.Save(context.Background(), fam)
	if err != nil {
		t.Fatalf("failed to save family: %v", err)
	}

	// Test
	result, err := svc.GetFamily(context.Background(), "00000000-0000-0000-0000-000000000004")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify
	if result.ID != fam.ID() {
		t.Errorf("expected ID %s, got %s", fam.ID(), result.ID)
	}
	if result.Status != string(fam.Status()) {
		t.Errorf("expected Status %s, got %s", fam.Status(), result.Status)
	}
	if len(result.Parents) != 1 {
		t.Errorf("expected 1 parent, got %d", len(result.Parents))
	}
}

func TestAddParent(t *testing.T) {
	// Setup
	repo := &mockRepo{store: make(map[string]*entity.Family)}
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(repo, contextLogger)

	// Create a parent
	birthDate1 := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p1, err := entity.NewParent("00000000-0000-0000-0000-000000000005", "John", "Doe", birthDate1, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a family
	fam, err := entity.NewFamily("00000000-0000-0000-0000-000000000006", entity.Single, []*entity.Parent{p1}, []*entity.Child{})
	if err != nil {
		t.Fatalf("failed to create family: %v", err)
	}

	// Save the family
	err = repo.Save(context.Background(), fam)
	if err != nil {
		t.Fatalf("failed to save family: %v", err)
	}

	// Create a second parent
	birthDate2 := time.Date(1982, 2, 2, 0, 0, 0, 0, time.UTC)
	p2DTO := entity.ParentDTO{
		ID:        "00000000-0000-0000-0000-000000000007",
		FirstName: "Jane",
		LastName:  "Doe",
		BirthDate: birthDate2,
		DeathDate: nil,
	}

	// Test
	result, err := svc.AddParent(context.Background(), "00000000-0000-0000-0000-000000000006", p2DTO)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify
	if len(result.Parents) != 2 {
		t.Errorf("expected 2 parents, got %d", len(result.Parents))
	}
	if result.Status != string(entity.Married) {
		t.Errorf("expected Status %s, got %s", entity.Married, result.Status)
	}
}

func TestAddChild(t *testing.T) {
	// Setup
	repo := &mockRepo{store: make(map[string]*entity.Family)}
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(repo, contextLogger)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("00000000-0000-0000-0000-000000000008", "John", "Doe", birthDate, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a family
	fam, err := entity.NewFamily("00000000-0000-0000-0000-000000000009", entity.Single, []*entity.Parent{p}, []*entity.Child{})
	if err != nil {
		t.Fatalf("failed to create family: %v", err)
	}

	// Save the family
	err = repo.Save(context.Background(), fam)
	if err != nil {
		t.Fatalf("failed to save family: %v", err)
	}

	// Create a child
	childBirthDate := time.Date(2010, 3, 3, 0, 0, 0, 0, time.UTC)
	childDTO := entity.ChildDTO{
		ID:        "00000000-0000-0000-0000-00000000000a",
		FirstName: "Baby",
		LastName:  "Doe",
		BirthDate: childBirthDate,
		DeathDate: nil,
	}

	// Test
	result, err := svc.AddChild(context.Background(), "00000000-0000-0000-0000-000000000009", childDTO)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify
	if len(result.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(result.Children))
	}
}

func TestDivorce(t *testing.T) {
	// Setup
	repo := &mockRepo{store: make(map[string]*entity.Family)}
	logger := zaptest.NewLogger(t)
	contextLogger := logging.NewContextLogger(logger)
	svc := NewFamilyDomainService(repo, contextLogger)

	// Create parents
	birthDate1 := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p1, err := entity.NewParent("00000000-0000-0000-0000-00000000000b", "John", "Doe", birthDate1, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	birthDate2 := time.Date(1982, 2, 2, 0, 0, 0, 0, time.UTC)
	p2, err := entity.NewParent("00000000-0000-0000-0000-00000000000c", "Jane", "Doe", birthDate2, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a child
	childBirthDate := time.Date(2010, 3, 3, 0, 0, 0, 0, time.UTC)
	c, err := entity.NewChild("00000000-0000-0000-0000-00000000000d", "Baby", "Doe", childBirthDate, nil)
	if err != nil {
		t.Fatalf("failed to create child: %v", err)
	}

	// Create a family
	fam, err := entity.NewFamily("00000000-0000-0000-0000-00000000000e", entity.Married, []*entity.Parent{p1, p2}, []*entity.Child{c})
	if err != nil {
		t.Fatalf("failed to create family: %v", err)
	}

	// Save the family
	err = repo.Save(context.Background(), fam)
	if err != nil {
		t.Fatalf("failed to save family: %v", err)
	}

	// Test
	result, err := svc.Divorce(context.Background(), "00000000-0000-0000-0000-00000000000e", "00000000-0000-0000-0000-00000000000b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the original family (now with custodial parent)
	if result.Status != string(entity.Divorced) {
		t.Errorf("expected Status %s, got %s", entity.Divorced, result.Status)
	}
	if len(result.Parents) != 1 {
		t.Errorf("expected 1 parent, got %d", len(result.Parents))
	}
	if len(result.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(result.Children))
	}
	// Verify the family with custodial parent keeps the original ID
	if result.ID != "00000000-0000-0000-0000-00000000000e" {
		t.Errorf("expected family with custodial parent to keep the original ID, got %s", result.ID)
	}
	// Verify the parent is the custodial parent
	if result.Parents[0].ID != "00000000-0000-0000-0000-00000000000b" {
		t.Errorf("expected custodial parent ID to be 00000000-0000-0000-0000-00000000000b, got %s", result.Parents[0].ID)
	}

	// Find the new family with the remaining parent
	// We need to get all families and find the one that's not the original
	allFamilies, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var remainingFam *entity.Family
	for _, f := range allFamilies {
		if f.ID() != "00000000-0000-0000-0000-00000000000e" {
			remainingFam = f
			break
		}
	}

	if remainingFam == nil {
		t.Fatalf("could not find family with remaining parent")
	}

	// Verify the family with remaining parent
	if remainingFam.Status() != entity.Divorced {
		t.Errorf("expected Status %s, got %s", entity.Divorced, remainingFam.Status())
	}
	if len(remainingFam.Parents()) != 1 {
		t.Errorf("expected 1 parent, got %d", len(remainingFam.Parents()))
	}
	if len(remainingFam.Children()) != 0 {
		t.Errorf("expected 0 children, got %d", len(remainingFam.Children()))
	}
	// Verify the parent is the remaining parent
	if remainingFam.Parents()[0].ID() != "00000000-0000-0000-0000-00000000000c" {
		t.Errorf("expected remaining parent ID to be 00000000-0000-0000-0000-00000000000c, got %s", remainingFam.Parents()[0].ID())
	}
}
