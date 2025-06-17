package services

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/servicelib/errors"
)

// mockRepo is a mock implementation of the FamilyRepository interface for testing
type mockRepo struct {
	store map[string]*entity.Family
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	fam, ok := m.store[id]
	if !ok {
		return nil, errors.NewNotFoundError("Family", id)
	}
	return fam, nil
}

func (m *mockRepo) Save(ctx context.Context, f *entity.Family) error {
	if f == nil {
		return errors.NewValidationError("family cannot be nil")
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
	return nil, errors.NewNotFoundError("Family with Child", childID)
}

func TestCreateFamily(t *testing.T) {
	// Setup
	repo := &mockRepo{store: make(map[string]*entity.Family)}
	svc := NewFamilyDomainService(repo)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("p1", "John", "Doe", birthDate, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a family DTO
	dto := entity.FamilyDTO{
		ID:       "abc123",
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
	svc := NewFamilyDomainService(repo)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("p1", "John", "Doe", birthDate, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a family
	fam, err := entity.NewFamily("abc123", entity.Single, []*entity.Parent{p}, []*entity.Child{})
	if err != nil {
		t.Fatalf("failed to create family: %v", err)
	}

	// Save the family
	err = repo.Save(context.Background(), fam)
	if err != nil {
		t.Fatalf("failed to save family: %v", err)
	}

	// Test
	result, err := svc.GetFamily(context.Background(), "abc123")
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
	svc := NewFamilyDomainService(repo)

	// Create a parent
	birthDate1 := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p1, err := entity.NewParent("p1", "John", "Doe", birthDate1, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a family
	fam, err := entity.NewFamily("abc123", entity.Single, []*entity.Parent{p1}, []*entity.Child{})
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
		ID:        "p2",
		FirstName: "Jane",
		LastName:  "Doe",
		BirthDate: birthDate2,
		DeathDate: nil,
	}

	// Test
	result, err := svc.AddParent(context.Background(), "abc123", p2DTO)
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
	svc := NewFamilyDomainService(repo)

	// Create a parent
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p, err := entity.NewParent("p1", "John", "Doe", birthDate, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a family
	fam, err := entity.NewFamily("abc123", entity.Single, []*entity.Parent{p}, []*entity.Child{})
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
		ID:        "c1",
		FirstName: "Baby",
		LastName:  "Doe",
		BirthDate: childBirthDate,
		DeathDate: nil,
	}

	// Test
	result, err := svc.AddChild(context.Background(), "abc123", childDTO)
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
	svc := NewFamilyDomainService(repo)

	// Create parents
	birthDate1 := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	p1, err := entity.NewParent("p1", "John", "Doe", birthDate1, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	birthDate2 := time.Date(1982, 2, 2, 0, 0, 0, 0, time.UTC)
	p2, err := entity.NewParent("p2", "Jane", "Doe", birthDate2, nil)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create a child
	childBirthDate := time.Date(2010, 3, 3, 0, 0, 0, 0, time.UTC)
	c, err := entity.NewChild("c1", "Baby", "Doe", childBirthDate, nil)
	if err != nil {
		t.Fatalf("failed to create child: %v", err)
	}

	// Create a family
	fam, err := entity.NewFamily("abc123", entity.Married, []*entity.Parent{p1, p2}, []*entity.Child{c})
	if err != nil {
		t.Fatalf("failed to create family: %v", err)
	}

	// Save the family
	err = repo.Save(context.Background(), fam)
	if err != nil {
		t.Fatalf("failed to save family: %v", err)
	}

	// Test
	result, err := svc.Divorce(context.Background(), "abc123", "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify
	if result.Status != string(entity.Divorced) {
		t.Errorf("expected Status %s, got %s", entity.Divorced, result.Status)
	}
	if len(result.Parents) != 1 {
		t.Errorf("expected 1 parent, got %d", len(result.Parents))
	}
	if len(result.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(result.Children))
	}

	// Check original family
	origFam, err := svc.GetFamily(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if origFam.Status != string(entity.Divorced) {
		t.Errorf("expected Status %s, got %s", entity.Divorced, origFam.Status)
	}
	if len(origFam.Parents) != 1 {
		t.Errorf("expected 1 parent, got %d", len(origFam.Parents))
	}
	if len(origFam.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(origFam.Children))
	}
}
