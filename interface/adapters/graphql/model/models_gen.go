// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// Input for creating or adding a child to a family.
type ChildInput struct {
	// Unique identifier for the child
	ID identification.ID `json:"id"`
	// First name of the child (1-50 characters)
	FirstName string `json:"firstName"`
	// Last name of the child (1-50 characters)
	LastName string `json:"lastName"`
	// Birth date of the child in RFC3339 format (YYYY-MM-DD).
	// Must not be in the future.
	BirthDate string `json:"birthDate"`
	// Death date of the child in RFC3339 format (YYYY-MM-DD), if applicable.
	// Must be after the birth date and not in the future.
	DeathDate *string `json:"deathDate,omitempty"`
}

// Error represents an error that occurred during a GraphQL operation.
// Errors provide information about what went wrong and where.
type Error struct {
	// Human-readable error message
	Message string `json:"message"`
	// Error code for programmatic handling (e.g., NOT_FOUND, VALIDATION_ERROR)
	Code *string `json:"code,omitempty"`
	// Path to the field that caused the error
	Path []string `json:"path,omitempty"`
}

// Family represents a family unit with parents and children.
// A family must have at least one parent and can have zero or more children.
// A family can have at most two parents.
type Family struct {
	// Unique identifier for the family
	ID identification.ID `json:"id"`
	// Current status of the family (SINGLE, MARRIED, DIVORCED, WIDOWED, or ABANDONED)
	Status FamilyStatus `json:"status"`
	// List of parents in the family (1-2 parents)
	Parents []*Parent `json:"parents"`
	// List of children in the family (0 or more)
	Children []*Child `json:"children"`
	// Number of parents in the family
	ParentCount int `json:"parentCount"`
	// Number of children in the family
	ChildrenCount int `json:"childrenCount"`
}

// Input for creating a new family.
// A family must have at least one parent and can have zero or more children.
// A family can have at most two parents.
type FamilyInput struct {
	// Unique identifier for the family
	ID identification.ID `json:"id"`
	// Status of the family (SINGLE, MARRIED, DIVORCED, WIDOWED, or ABANDONED).
	// Must be consistent with the number of parents:
	// - SINGLE: One parent
	// - MARRIED: Two parents
	// - Other statuses have specific business rules
	Status FamilyStatus `json:"status"`
	// List of parents in the family (1-2 parents).
	// Parents must be at least 18 years old.
	Parents []*ParentInput `json:"parents"`
	// List of children in the family (0 or more).
	Children []*ChildInput `json:"children"`
}

// Mutations for modifying family data.
// All mutations require appropriate authorization.
type Mutation struct {
}

// Input for creating or adding a parent to a family.
// Parents must be at least 18 years old.
type ParentInput struct {
	// Unique identifier for the parent
	ID identification.ID `json:"id"`
	// First name of the parent (1-50 characters)
	FirstName string `json:"firstName"`
	// Last name of the parent (1-50 characters)
	LastName string `json:"lastName"`
	// Birth date of the parent in RFC3339 format (YYYY-MM-DD).
	// Must be at least 18 years before the current date.
	BirthDate string `json:"birthDate"`
	// Death date of the parent in RFC3339 format (YYYY-MM-DD), if applicable.
	// Must be after the birth date and not in the future.
	DeathDate *string `json:"deathDate,omitempty"`
}

// Queries for retrieving family data.
// All queries require appropriate authorization.
type Query struct {
}

// FamilyStatus represents the current status of a family.
// The status affects what operations can be performed on the family.
type FamilyStatus string

const (
	// Single parent family with one parent
	FamilyStatusSingle FamilyStatus = "SINGLE"
	// Family with two parents who are married
	FamilyStatusMarried FamilyStatus = "MARRIED"
	// Family where the parents have divorced
	FamilyStatusDivorced FamilyStatus = "DIVORCED"
	// Family where one parent has died
	FamilyStatusWidowed FamilyStatus = "WIDOWED"
	// Family that has been abandoned
	FamilyStatusAbandoned FamilyStatus = "ABANDONED"
	// Active family status (used in tests)
	FamilyStatusActive FamilyStatus = "ACTIVE"
)

var AllFamilyStatus = []FamilyStatus{
	FamilyStatusSingle,
	FamilyStatusMarried,
	FamilyStatusDivorced,
	FamilyStatusWidowed,
	FamilyStatusAbandoned,
	FamilyStatusActive,
}

func (e FamilyStatus) IsValid() bool {
	switch e {
	case FamilyStatusSingle, FamilyStatusMarried, FamilyStatusDivorced, FamilyStatusWidowed, FamilyStatusAbandoned, FamilyStatusActive:
		return true
	}
	return false
}

func (e FamilyStatus) String() string {
	return string(e)
}

func (e *FamilyStatus) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FamilyStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FamilyStatus", str)
	}
	return nil
}

func (e FamilyStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *FamilyStatus) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e FamilyStatus) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

// Resource represents the resource type being accessed.
// Different resources may have different access controls.
type Resource string

const (
	// Family resource type
	ResourceFamily Resource = "FAMILY"
	// Parent resource type
	ResourceParent Resource = "PARENT"
	// Child resource type
	ResourceChild Resource = "CHILD"
)

var AllResource = []Resource{
	ResourceFamily,
	ResourceParent,
	ResourceChild,
}

func (e Resource) IsValid() bool {
	switch e {
	case ResourceFamily, ResourceParent, ResourceChild:
		return true
	}
	return false
}

func (e Resource) String() string {
	return string(e)
}

func (e *Resource) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Resource(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Resource", str)
	}
	return nil
}

func (e Resource) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *Resource) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e Resource) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

// Role represents the authorization role of a user.
// Different roles have different levels of access to the system.
type Role string

const (
	// Administrator with full access to all operations
	RoleAdmin Role = "ADMIN"
	// Editor with permission to modify data but not administer the system
	RoleEditor Role = "EDITOR"
	// Viewer with read-only access to data
	RoleViewer Role = "VIEWER"
)

var AllRole = []Role{
	RoleAdmin,
	RoleEditor,
	RoleViewer,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleEditor, RoleViewer:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *Role) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e Role) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

// Scope represents the permission scope for a resource.
// Scopes define what actions can be performed on resources.
type Scope string

const (
	// Permission to read/view a resource
	ScopeRead Scope = "READ"
	// Permission to modify an existing resource
	ScopeWrite Scope = "WRITE"
	// Permission to delete a resource
	ScopeDelete Scope = "DELETE"
	// Permission to create a new resource
	ScopeCreate Scope = "CREATE"
)

var AllScope = []Scope{
	ScopeRead,
	ScopeWrite,
	ScopeDelete,
	ScopeCreate,
}

func (e Scope) IsValid() bool {
	switch e {
	case ScopeRead, ScopeWrite, ScopeDelete, ScopeCreate:
		return true
	}
	return false
}

func (e Scope) String() string {
	return string(e)
}

func (e *Scope) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Scope(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Scope", str)
	}
	return nil
}

func (e Scope) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *Scope) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e Scope) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}
