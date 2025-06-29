// Copyright (c) 2025 A Bit of Help, Inc.

package dto

import (
	"fmt"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

const (
	// ISO8601DateFormat is the standard date format used across the application
	// Using RFC3339 format as required by the project guidelines
	ISO8601DateFormat = time.RFC3339
)

// FamilyMapper provides methods to convert between domain DTOs and GraphQL models
type FamilyMapper interface {
	ToDomain(input model.FamilyInput) (entity.FamilyDTO, error)
	ToGraphQL(dto entity.FamilyDTO) (*model.Family, error)
	ToParentDTO(input model.ParentInput) (entity.ParentDTO, error)
	ToChildDTO(input model.ChildInput) (entity.ChildDTO, error)
	ToParent(dto entity.ParentDTO) (*model.Parent, error)
	ToChild(dto entity.ChildDTO) (*model.Child, error)
}

// familyMapper implements FamilyMapper
type familyMapper struct{}

// NewFamilyMapper creates a new instance of FamilyMapper
func NewFamilyMapper() FamilyMapper {
	return &familyMapper{}
}

func (m *familyMapper) ToDomain(input model.FamilyInput) (entity.FamilyDTO, error) {
	if input.ID == "" {
		return entity.FamilyDTO{}, fmt.Errorf("invalid ID: ID cannot be empty")
	}

	// Validate status
	status := string(input.Status)
	if status != "ACTIVE" && status != "DIVORCED" {
		return entity.FamilyDTO{}, fmt.Errorf("invalid family status: %s", status)
	}

	// Convert parents
	parents := make([]entity.ParentDTO, 0, len(input.Parents))
	for _, p := range input.Parents {
		parent, err := m.ToParentDTO(*p)
		if err != nil {
			return entity.FamilyDTO{}, fmt.Errorf("invalid parent: %w", err)
		}
		parents = append(parents, parent)
	}

	// Convert children
	children := make([]entity.ChildDTO, 0, len(input.Children))
	for _, c := range input.Children {
		child, err := m.ToChildDTO(*c)
		if err != nil {
			return entity.FamilyDTO{}, fmt.Errorf("invalid child: %w", err)
		}
		children = append(children, child)
	}

	return entity.FamilyDTO{
		ID:       input.ID.String(),
		Status:   status,
		Parents:  parents,
		Children: children,
	}, nil
}

func (m *familyMapper) ToGraphQL(dto entity.FamilyDTO) (*model.Family, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("invalid ID: ID cannot be empty")
	}

	// Convert parents
	parents := make([]*model.Parent, 0, len(dto.Parents))
	for _, p := range dto.Parents {
		parent, err := m.ToParent(p)
		if err != nil {
			return nil, fmt.Errorf("invalid parent: %w", err)
		}
		parents = append(parents, parent)
	}

	// Convert children
	children := make([]*model.Child, 0, len(dto.Children))
	for _, c := range dto.Children {
		child, err := m.ToChild(c)
		if err != nil {
			return nil, fmt.Errorf("invalid child: %w", err)
		}
		children = append(children, child)
	}

	// Validate and convert status
	status := model.FamilyStatus(dto.Status)
	if status != "ACTIVE" && status != "DIVORCED" {
		return nil, fmt.Errorf("invalid family status: %s", dto.Status)
	}

	return &model.Family{
		ID:       identification.ID(dto.ID),
		Status:   status,
		Parents:  parents,
		Children: children,
	}, nil
}

func (m *familyMapper) ToParentDTO(input model.ParentInput) (entity.ParentDTO, error) {
	if input.ID == "" {
		return entity.ParentDTO{}, fmt.Errorf("invalid ID: ID cannot be empty")
	}

	birthDate, err := time.Parse(ISO8601DateFormat, input.BirthDate)
	if err != nil {
		return entity.ParentDTO{}, fmt.Errorf("invalid birth date: %w", err)
	}

	var deathDate *time.Time
	if input.DeathDate != nil {
		parsed, err := time.Parse(ISO8601DateFormat, *input.DeathDate)
		if err != nil {
			return entity.ParentDTO{}, fmt.Errorf("invalid death date: %w", err)
		}
		if parsed.Before(birthDate) {
			return entity.ParentDTO{}, fmt.Errorf("death date cannot be before birth date")
		}
		deathDate = &parsed
	}

	return entity.ParentDTO{
		ID:        input.ID.String(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		BirthDate: birthDate,
		DeathDate: deathDate,
	}, nil
}

func (m *familyMapper) ToChildDTO(input model.ChildInput) (entity.ChildDTO, error) {
	if input.ID == "" {
		return entity.ChildDTO{}, fmt.Errorf("invalid ID: ID cannot be empty")
	}

	birthDate, err := time.Parse(ISO8601DateFormat, input.BirthDate)
	if err != nil {
		return entity.ChildDTO{}, fmt.Errorf("invalid birth date: %w", err)
	}

	var deathDate *time.Time
	if input.DeathDate != nil {
		parsed, err := time.Parse(ISO8601DateFormat, *input.DeathDate)
		if err != nil {
			return entity.ChildDTO{}, fmt.Errorf("invalid death date: %w", err)
		}
		if parsed.Before(birthDate) {
			return entity.ChildDTO{}, fmt.Errorf("death date cannot be before birth date")
		}
		deathDate = &parsed
	}

	return entity.ChildDTO{
		ID:        input.ID.String(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		BirthDate: birthDate,
		DeathDate: deathDate,
	}, nil
}

func (m *familyMapper) ToParent(dto entity.ParentDTO) (*model.Parent, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("invalid ID: ID cannot be empty")
	}

	var deathDate *string
	if dto.DeathDate != nil {
		formatted := dto.DeathDate.Format(ISO8601DateFormat)
		deathDate = &formatted
	}

	return &model.Parent{
		ID:        identification.ID(dto.ID),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		BirthDate: dto.BirthDate.Format(ISO8601DateFormat),
		DeathDate: deathDate,
	}, nil
}

func (m *familyMapper) ToChild(dto entity.ChildDTO) (*model.Child, error) {
	if dto.ID == "" {
		return nil, fmt.Errorf("invalid ID: ID cannot be empty")
	}

	var deathDate *string
	if dto.DeathDate != nil {
		formatted := dto.DeathDate.Format(ISO8601DateFormat)
		deathDate = &formatted
	}

	return &model.Child{
		ID:        identification.ID(dto.ID),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		BirthDate: dto.BirthDate.Format(ISO8601DateFormat),
		DeathDate: deathDate,
	}, nil
}
