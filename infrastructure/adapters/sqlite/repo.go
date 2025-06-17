// Copyright (c) 2025 A Bit of Help, Inc.

package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/servicelib/errors"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SQLiteFamilyRepository implements the ports.FamilyRepository interface for SQLite
type SQLiteFamilyRepository struct {
	DB *sql.DB
}

// Ensure SQLiteFamilyRepository implements ports.FamilyRepository
var _ ports.FamilyRepository = (*SQLiteFamilyRepository)(nil)

// NewSQLiteFamilyRepository creates a new SQLiteFamilyRepository
func NewSQLiteFamilyRepository(db *sql.DB) *SQLiteFamilyRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &SQLiteFamilyRepository{DB: db}
}

// ensureTableExists creates the families table if it doesn't exist
func (r *SQLiteFamilyRepository) ensureTableExists(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS families (
		id TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		parents TEXT NOT NULL,
		children TEXT NOT NULL
	);
	`
	_, err := r.DB.ExecContext(ctx, query)
	if err != nil {
		return errors.NewRepositoryError(err, "failed to create families table", "SQLITE_ERROR")
	}
	return nil
}

// GetByID retrieves a family by its ID
func (r *SQLiteFamilyRepository) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	if id == "" {
		return nil, errors.NewValidationError("id is required")
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	var famID string
	var statusStr string
	var parentsData, childrenData string

	query := "SELECT id, status, parents, children FROM families WHERE id = ?"
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&famID, &statusStr, &parentsData, &childrenData)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Family", id)
		}
		return nil, errors.NewRepositoryError(err, "failed to get family from SQLite", "SQLITE_ERROR")
	}

	// Parse parents JSON
	var parentDTOs []entity.ParentDTO
	if err := json.Unmarshal([]byte(parentsData), &parentDTOs); err != nil {
		return nil, errors.NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
	}

	// Convert parent DTOs to domain entities
	parents := make([]*entity.Parent, 0, len(parentDTOs))
	for _, dto := range parentDTOs {
		p, err := entity.ParentFromDTO(dto)
		if err != nil {
			return nil, errors.NewRepositoryError(err, "failed to convert parent DTO to entity", "CONVERSION_ERROR")
		}
		parents = append(parents, p)
	}

	// Parse children JSON
	var childDTOs []entity.ChildDTO
	if err := json.Unmarshal([]byte(childrenData), &childDTOs); err != nil {
		return nil, errors.NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
	}

	// Convert child DTOs to domain entities
	children := make([]*entity.Child, 0, len(childDTOs))
	for _, dto := range childDTOs {
		c, err := entity.ChildFromDTO(dto)
		if err != nil {
			return nil, errors.NewRepositoryError(err, "failed to convert child DTO to entity", "CONVERSION_ERROR")
		}
		children = append(children, c)
	}

	// Create family entity
	return entity.NewFamily(famID, entity.Status(statusStr), parents, children)
}

// Save persists a family
func (r *SQLiteFamilyRepository) Save(ctx context.Context, fam *entity.Family) error {
	if fam == nil {
		return errors.NewValidationError("family cannot be nil")
	}

	if err := fam.Validate(); err != nil {
		return err
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return err
	}

	// Begin transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.NewRepositoryError(err, "failed to begin transaction", "SQLITE_ERROR")
	}

	// Ensure transaction is rolled back if an error occurs
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Log rollback error, but don't return it as it would mask the original error
				fmt.Printf("Error rolling back transaction: %v\n", rollbackErr)
			}
		}
	}()

	// Convert parents to DTOs for JSON serialization
	parentDTOs := make([]entity.ParentDTO, 0, len(fam.Parents()))
	for _, p := range fam.Parents() {
		parentDTOs = append(parentDTOs, p.ToDTO())
	}

	// Convert children to DTOs for JSON serialization
	childDTOs := make([]entity.ChildDTO, 0, len(fam.Children()))
	for _, c := range fam.Children() {
		childDTOs = append(childDTOs, c.ToDTO())
	}

	// Marshal to JSON
	parentsJSON, err := json.Marshal(parentDTOs)
	if err != nil {
		return errors.NewRepositoryError(err, "failed to marshal parents to JSON", "JSON_ERROR")
	}

	childrenJSON, err := json.Marshal(childDTOs)
	if err != nil {
		return errors.NewRepositoryError(err, "failed to marshal children to JSON", "JSON_ERROR")
	}

	// Check if family exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT 1 FROM families WHERE id = ?", fam.ID()).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return errors.NewRepositoryError(err, "failed to check if family exists", "SQLITE_ERROR")
	}

	var query string
	var args []interface{}

	if err == sql.ErrNoRows {
		// Insert new family
		query = "INSERT INTO families (id, status, parents, children) VALUES (?, ?, ?, ?)"
		args = []interface{}{fam.ID(), string(fam.Status()), parentsJSON, childrenJSON}
	} else {
		// Update existing family
		query = "UPDATE families SET status = ?, parents = ?, children = ? WHERE id = ?"
		args = []interface{}{string(fam.Status()), parentsJSON, childrenJSON, fam.ID()}
	}

	// Execute SQL
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.NewRepositoryError(err, "failed to save family to SQLite", "SQLITE_ERROR")
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return errors.NewRepositoryError(err, "failed to commit transaction", "SQLITE_ERROR")
	}

	return nil
}

// FindByParentID finds families that contain a specific parent
func (r *SQLiteFamilyRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error) {
	if parentID == "" {
		return nil, errors.NewValidationError("parent ID is required")
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	// SQLite doesn't have native JSON path operators like PostgreSQL,
	// so we need to fetch all families and filter in application code
	rows, err := r.DB.QueryContext(ctx, "SELECT id, status, parents, children FROM families")
	if err != nil {
		return nil, errors.NewRepositoryError(err, "failed to query families", "SQLITE_ERROR")
	}
	defer rows.Close()

	var families []*entity.Family

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData string

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to scan family row", "SQLITE_ERROR")
		}

		// Parse parents JSON
		var parentDTOs []entity.ParentDTO
		if err := json.Unmarshal([]byte(parentsData), &parentDTOs); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
		}

		// Check if any parent has the specified ID
		hasParent := false
		for _, dto := range parentDTOs {
			if dto.ID == parentID {
				hasParent = true
				break
			}
		}

		if !hasParent {
			continue // Skip this family if it doesn't have the parent
		}

		// Convert parent DTOs to domain entities
		parents := make([]*entity.Parent, 0, len(parentDTOs))
		for _, dto := range parentDTOs {
			p, err := entity.ParentFromDTO(dto)
			if err != nil {
				return nil, errors.NewRepositoryError(err, "failed to convert parent DTO to entity", "CONVERSION_ERROR")
			}
			parents = append(parents, p)
		}

		// Parse children JSON
		var childDTOs []entity.ChildDTO
		if err := json.Unmarshal([]byte(childrenData), &childDTOs); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
		}

		// Convert child DTOs to domain entities
		children := make([]*entity.Child, 0, len(childDTOs))
		for _, dto := range childDTOs {
			c, err := entity.ChildFromDTO(dto)
			if err != nil {
				return nil, errors.NewRepositoryError(err, "failed to convert child DTO to entity", "CONVERSION_ERROR")
			}
			children = append(children, c)
		}

		// Create family entity
		fam, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
		if err != nil {
			return nil, errors.NewRepositoryError(err, "failed to create family entity", "CONVERSION_ERROR")
		}

		families = append(families, fam)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewRepositoryError(err, "error iterating over family rows", "SQLITE_ERROR")
	}

	return families, nil
}

// GetAll retrieves all families
func (r *SQLiteFamilyRepository) GetAll(ctx context.Context) ([]*entity.Family, error) {
	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	// Query all families
	rows, err := r.DB.QueryContext(ctx, "SELECT id, status, parents, children FROM families")
	if err != nil {
		return nil, errors.NewRepositoryError(err, "failed to query families", "SQLITE_ERROR")
	}
	defer rows.Close()

	var families []*entity.Family

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData string

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to scan family row", "SQLITE_ERROR")
		}

		// Parse parents JSON
		var parentDTOs []entity.ParentDTO
		if err := json.Unmarshal([]byte(parentsData), &parentDTOs); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
		}

		// Convert parent DTOs to domain entities
		parents := make([]*entity.Parent, 0, len(parentDTOs))
		for _, dto := range parentDTOs {
			p, err := entity.ParentFromDTO(dto)
			if err != nil {
				return nil, errors.NewRepositoryError(err, "failed to convert parent DTO to entity", "CONVERSION_ERROR")
			}
			parents = append(parents, p)
		}

		// Parse children JSON
		var childDTOs []entity.ChildDTO
		if err := json.Unmarshal([]byte(childrenData), &childDTOs); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
		}

		// Convert child DTOs to domain entities
		children := make([]*entity.Child, 0, len(childDTOs))
		for _, dto := range childDTOs {
			c, err := entity.ChildFromDTO(dto)
			if err != nil {
				return nil, errors.NewRepositoryError(err, "failed to convert child DTO to entity", "CONVERSION_ERROR")
			}
			children = append(children, c)
		}

		// Create family entity
		fam, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
		if err != nil {
			return nil, errors.NewRepositoryError(err, "failed to create family entity", "CONVERSION_ERROR")
		}

		families = append(families, fam)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewRepositoryError(err, "error iterating over family rows", "SQLITE_ERROR")
	}

	return families, nil
}

// FindByChildID finds the family that contains a specific child
func (r *SQLiteFamilyRepository) FindByChildID(ctx context.Context, childID string) (*entity.Family, error) {
	if childID == "" {
		return nil, errors.NewValidationError("child ID is required")
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	// SQLite doesn't have native JSON path operators like PostgreSQL,
	// so we need to fetch all families and filter in application code
	rows, err := r.DB.QueryContext(ctx, "SELECT id, status, parents, children FROM families")
	if err != nil {
		return nil, errors.NewRepositoryError(err, "failed to query families", "SQLITE_ERROR")
	}
	defer rows.Close()

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData string

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to scan family row", "SQLITE_ERROR")
		}

		// Parse children JSON
		var childDTOs []entity.ChildDTO
		if err := json.Unmarshal([]byte(childrenData), &childDTOs); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
		}

		// Check if any child has the specified ID
		hasChild := false
		for _, dto := range childDTOs {
			if dto.ID == childID {
				hasChild = true
				break
			}
		}

		if !hasChild {
			continue // Skip this family if it doesn't have the child
		}

		// Parse parents JSON
		var parentDTOs []entity.ParentDTO
		if err := json.Unmarshal([]byte(parentsData), &parentDTOs); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
		}

		// Convert parent DTOs to domain entities
		parents := make([]*entity.Parent, 0, len(parentDTOs))
		for _, dto := range parentDTOs {
			p, err := entity.ParentFromDTO(dto)
			if err != nil {
				return nil, errors.NewRepositoryError(err, "failed to convert parent DTO to entity", "CONVERSION_ERROR")
			}
			parents = append(parents, p)
		}

		// Convert child DTOs to domain entities
		children := make([]*entity.Child, 0, len(childDTOs))
		for _, dto := range childDTOs {
			c, err := entity.ChildFromDTO(dto)
			if err != nil {
				return nil, errors.NewRepositoryError(err, "failed to convert child DTO to entity", "CONVERSION_ERROR")
			}
			children = append(children, c)
		}

		// Create family entity
		fam, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
		if err != nil {
			return nil, errors.NewRepositoryError(err, "failed to create family entity", "CONVERSION_ERROR")
		}

		return fam, nil
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewRepositoryError(err, "error iterating over family rows", "SQLITE_ERROR")
	}

	return nil, errors.NewNotFoundError("Family with Child", childID)
}
