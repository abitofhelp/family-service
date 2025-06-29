// Copyright (c) 2025 A Bit of Help, Inc.

package repository

import (
	"github.com/abitofhelp/family-service/core/domain/entity"
)

// EntityConverter defines the interface for converting between domain entities and database-specific formats
type EntityConverter interface {
	// EntityToDocument converts a domain entity to a database-specific document
	EntityToDocument(entity interface{}) (interface{}, error)

	// DocumentToEntity converts a database-specific document to a domain entity
	DocumentToEntity(document interface{}) (interface{}, error)
}

// FamilyConverter defines the interface for converting between Family entities and database-specific formats
type FamilyConverter interface {
	// FamilyToDocument converts a Family entity to a database-specific document
	FamilyToDocument(family *entity.Family) (interface{}, error)

	// DocumentToFamily converts a database-specific document to a Family entity
	DocumentToFamily(document interface{}) (*entity.Family, error)

	// ParentToDocument converts a Parent entity to a database-specific document
	ParentToDocument(parent *entity.Parent) (interface{}, error)

	// DocumentToParent converts a database-specific document to a Parent entity
	DocumentToParent(document interface{}) (*entity.Parent, error)

	// ChildToDocument converts a Child entity to a database-specific document
	ChildToDocument(child *entity.Child) (interface{}, error)

	// DocumentToChild converts a database-specific document to a Child entity
	DocumentToChild(document interface{}) (*entity.Child, error)
}

// DateFormat is the standard date format used across all repositories
const DateFormat = "2006-01-02T15:04:05Z07:00" // RFC3339
