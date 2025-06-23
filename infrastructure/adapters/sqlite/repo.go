// Copyright (c) 2025 A Bit of Help, Inc.

package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"go.uber.org/zap"
)

// NewRepositoryError is a helper function that wraps errors.NewDatabaseError
// to maintain compatibility with the old errors.NewDatabaseError function.
func NewRepositoryError(err error, message string, code string) error {
	// Map the code to an appropriate operation and table
	operation := "operation"
	table := "families"

	// Map common codes to operations
	switch code {
	case "SQLITE_ERROR":
		operation = "query"
	case "JSON_ERROR":
		operation = "unmarshal"
	case "DATA_FORMAT_ERROR":
		operation = "parse"
	case "CONVERSION_ERROR":
		operation = "convert"
	}

	return errors.NewDatabaseError(message, operation, table, err)
}

// SQLiteFamilyRepository implements the ports.FamilyRepository interface for SQLite
type SQLiteFamilyRepository struct {
	DB     *sql.DB
	logger *logging.ContextLogger
}

// Ensure SQLiteFamilyRepository implements ports.FamilyRepository
var _ ports.FamilyRepository = (*SQLiteFamilyRepository)(nil)

// NewSQLiteFamilyRepository creates a new SQLiteFamilyRepository
func NewSQLiteFamilyRepository(db *sql.DB, logger *logging.ContextLogger) *SQLiteFamilyRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	if logger == nil {
		panic("logger cannot be nil")
	}
	return &SQLiteFamilyRepository{
		DB:     db,
		logger: logger,
	}
}

// ensureTableExists creates the families table if it doesn't exist
func (r *SQLiteFamilyRepository) ensureTableExists(ctx context.Context) error {
	r.logger.Debug(ctx, "Ensuring families table exists in SQLite")

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
		r.logger.Error(ctx, "Failed to create families table in SQLite", zap.Error(err))
		return NewRepositoryError(err, "failed to create families table", "SQLITE_ERROR")
	}

	r.logger.Debug(ctx, "Families table exists in SQLite")
	return nil
}

// GetByID retrieves a family by its ID
func (r *SQLiteFamilyRepository) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	r.logger.Debug(ctx, "Getting family by ID from SQLite", zap.String("family_id", id))

	if id == "" {
		r.logger.Warn(ctx, "Family ID is required for GetByID")
		return nil, errors.NewValidationError("id is required", "id", nil)
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
			r.logger.Info(ctx, "Family not found in SQLite", zap.String("family_id", id))
			return nil, errors.NewNotFoundError("Family", id, nil)
		}
		r.logger.Error(ctx, "Failed to get family from SQLite", zap.Error(err), zap.String("family_id", id))
		return nil, NewRepositoryError(err, "failed to get family from SQLite", "SQLITE_ERROR")
	}

	// Parse parents JSON
	var parentDTOs []entity.ParentDTO
	if err := json.Unmarshal([]byte(parentsData), &parentDTOs); err != nil {
		r.logger.Error(ctx, "Failed to unmarshal parents data", zap.Error(err), zap.String("family_id", id))
		return nil, NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
	}

	// Convert parent DTOs to domain entities
	parents := make([]*entity.Parent, 0, len(parentDTOs))
	for _, dto := range parentDTOs {
		p, err := entity.ParentFromDTO(dto)
		if err != nil {
			r.logger.Error(ctx, "Failed to convert parent DTO to entity",
				zap.Error(err),
				zap.String("family_id", id),
				zap.String("parent_id", dto.ID))
			return nil, NewRepositoryError(err, "failed to convert parent DTO to entity", "CONVERSION_ERROR")
		}
		parents = append(parents, p)
	}

	// Parse children JSON
	var childDTOs []entity.ChildDTO
	if err := json.Unmarshal([]byte(childrenData), &childDTOs); err != nil {
		r.logger.Error(ctx, "Failed to unmarshal children data", zap.Error(err), zap.String("family_id", id))
		return nil, NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
	}

	// Convert child DTOs to domain entities
	children := make([]*entity.Child, 0, len(childDTOs))
	for _, dto := range childDTOs {
		c, err := entity.ChildFromDTO(dto)
		if err != nil {
			r.logger.Error(ctx, "Failed to convert child DTO to entity",
				zap.Error(err),
				zap.String("family_id", id),
				zap.String("child_id", dto.ID))
			return nil, NewRepositoryError(err, "failed to convert child DTO to entity", "CONVERSION_ERROR")
		}
		children = append(children, c)
	}

	// Create family entity
	family, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
	if err != nil {
		r.logger.Error(ctx, "Failed to create family entity", zap.Error(err), zap.String("family_id", id))
		return nil, err
	}

	r.logger.Debug(ctx, "Successfully retrieved family from SQLite",
		zap.String("family_id", id),
		zap.String("status", statusStr),
		zap.Int("parent_count", len(parents)),
		zap.Int("children_count", len(children)))
	return family, nil
}

// Save persists a family
func (r *SQLiteFamilyRepository) Save(ctx context.Context, fam *entity.Family) error {
	r.logger.Debug(ctx, "Saving family to SQLite", zap.String("family_id", fam.ID()))

	if fam == nil {
		r.logger.Warn(ctx, "Family cannot be nil for Save")
		return errors.NewValidationError("family cannot be nil", "family", nil)
	}

	if err := fam.Validate(); err != nil {
		r.logger.Error(ctx, "Family validation failed", zap.Error(err), zap.String("family_id", fam.ID()))
		return err
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return err
	}

	// Begin transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error(ctx, "Failed to begin transaction", zap.Error(err), zap.String("family_id", fam.ID()))
		return NewRepositoryError(err, "failed to begin transaction", "SQLITE_ERROR")
	}

	// Ensure transaction is rolled back if an error occurs
	committed := false
	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Log rollback error, but don't return it as it would mask the original error
				r.logger.Error(ctx, "Error rolling back transaction",
					zap.Error(rollbackErr),
					zap.String("family_id", fam.ID()))
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
		r.logger.Error(ctx, "Failed to marshal parents to JSON",
			zap.Error(err),
			zap.String("family_id", fam.ID()))
		return NewRepositoryError(err, "failed to marshal parents to JSON", "JSON_ERROR")
	}

	childrenJSON, err := json.Marshal(childDTOs)
	if err != nil {
		r.logger.Error(ctx, "Failed to marshal children to JSON",
			zap.Error(err),
			zap.String("family_id", fam.ID()))
		return NewRepositoryError(err, "failed to marshal children to JSON", "JSON_ERROR")
	}

	// Check if family exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT 1 FROM families WHERE id = ?", fam.ID()).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Error(ctx, "Failed to check if family exists",
			zap.Error(err),
			zap.String("family_id", fam.ID()))
		return NewRepositoryError(err, "failed to check if family exists", "SQLITE_ERROR")
	}

	var query string
	var args []interface{}
	var operation string

	if err == sql.ErrNoRows {
		// Insert new family
		operation = "insert"
		query = "INSERT INTO families (id, status, parents, children) VALUES (?, ?, ?, ?)"
		args = []interface{}{fam.ID(), string(fam.Status()), parentsJSON, childrenJSON}
		r.logger.Debug(ctx, "Inserting new family",
			zap.String("family_id", fam.ID()),
			zap.String("status", string(fam.Status())))
	} else {
		// Update existing family
		operation = "update"
		query = "UPDATE families SET status = ?, parents = ?, children = ? WHERE id = ?"
		args = []interface{}{string(fam.Status()), parentsJSON, childrenJSON, fam.ID()}
		r.logger.Debug(ctx, "Updating existing family",
			zap.String("family_id", fam.ID()),
			zap.String("status", string(fam.Status())))
	}

	// Execute SQL
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error(ctx, "Failed to save family to SQLite",
			zap.Error(err),
			zap.String("family_id", fam.ID()),
			zap.String("operation", operation))
		return NewRepositoryError(err, "failed to save family to SQLite", "SQLITE_ERROR")
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		r.logger.Error(ctx, "Failed to commit transaction",
			zap.Error(err),
			zap.String("family_id", fam.ID()))
		return NewRepositoryError(err, "failed to commit transaction", "SQLITE_ERROR")
	}
	committed = true

	r.logger.Info(ctx, "Successfully saved family to SQLite",
		zap.String("family_id", fam.ID()),
		zap.String("status", string(fam.Status())),
		zap.String("operation", operation),
		zap.Int("parent_count", len(parentDTOs)),
		zap.Int("children_count", len(childDTOs)))
	return nil
}

// FindByParentID finds families that contain a specific parent
func (r *SQLiteFamilyRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error) {
	r.logger.Debug(ctx, "Finding families by parent ID in SQLite", zap.String("parent_id", parentID))

	if parentID == "" {
		r.logger.Warn(ctx, "Parent ID is required for FindByParentID")
		return nil, errors.NewValidationError("parent ID is required", "parentID", nil)
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	// SQLite doesn't have native JSON path operators like PostgreSQL,
	// so we need to fetch all families and filter in application code
	r.logger.Debug(ctx, "Querying all families to filter by parent ID", zap.String("parent_id", parentID))
	rows, err := r.DB.QueryContext(ctx, "SELECT id, status, parents, children FROM families")
	if err != nil {
		r.logger.Error(ctx, "Failed to query families", zap.Error(err))
		return nil, NewRepositoryError(err, "failed to query families", "SQLITE_ERROR")
	}
	defer rows.Close()

	var families []*entity.Family
	var matchCount int

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData string

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			r.logger.Error(ctx, "Failed to scan family row", zap.Error(err))
			return nil, NewRepositoryError(err, "failed to scan family row", "SQLITE_ERROR")
		}

		// Parse parents JSON
		var parentDTOs []entity.ParentDTO
		if err := json.Unmarshal([]byte(parentsData), &parentDTOs); err != nil {
			r.logger.Error(ctx, "Failed to unmarshal parents data",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
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

		r.logger.Debug(ctx, "Found family with matching parent",
			zap.String("family_id", famID),
			zap.String("parent_id", parentID))
		matchCount++

		// Convert parent DTOs to domain entities
		parents := make([]*entity.Parent, 0, len(parentDTOs))
		for _, dto := range parentDTOs {
			p, err := entity.ParentFromDTO(dto)
			if err != nil {
				r.logger.Error(ctx, "Failed to convert parent DTO to entity",
					zap.Error(err),
					zap.String("family_id", famID),
					zap.String("parent_id", dto.ID))
				return nil, NewRepositoryError(err, "failed to convert parent DTO to entity", "CONVERSION_ERROR")
			}
			parents = append(parents, p)
		}

		// Parse children JSON
		var childDTOs []entity.ChildDTO
		if err := json.Unmarshal([]byte(childrenData), &childDTOs); err != nil {
			r.logger.Error(ctx, "Failed to unmarshal children data",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
		}

		// Convert child DTOs to domain entities
		children := make([]*entity.Child, 0, len(childDTOs))
		for _, dto := range childDTOs {
			c, err := entity.ChildFromDTO(dto)
			if err != nil {
				r.logger.Error(ctx, "Failed to convert child DTO to entity",
					zap.Error(err),
					zap.String("family_id", famID),
					zap.String("child_id", dto.ID))
				return nil, NewRepositoryError(err, "failed to convert child DTO to entity", "CONVERSION_ERROR")
			}
			children = append(children, c)
		}

		// Create family entity
		fam, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
		if err != nil {
			r.logger.Error(ctx, "Failed to create family entity",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to create family entity", "CONVERSION_ERROR")
		}

		families = append(families, fam)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, "Error iterating over family rows", zap.Error(err))
		return nil, NewRepositoryError(err, "error iterating over family rows", "SQLITE_ERROR")
	}

	r.logger.Info(ctx, "Successfully found families by parent ID",
		zap.String("parent_id", parentID),
		zap.Int("family_count", len(families)),
		zap.Int("total_matches", matchCount))
	return families, nil
}

// GetAll retrieves all families
func (r *SQLiteFamilyRepository) GetAll(ctx context.Context) ([]*entity.Family, error) {
	r.logger.Debug(ctx, "Getting all families from SQLite")

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	// Query all families
	rows, err := r.DB.QueryContext(ctx, "SELECT id, status, parents, children FROM families")
	if err != nil {
		r.logger.Error(ctx, "Failed to query all families", zap.Error(err))
		return nil, NewRepositoryError(err, "failed to query families", "SQLITE_ERROR")
	}
	defer rows.Close()

	var families []*entity.Family

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData string

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			r.logger.Error(ctx, "Failed to scan family row", zap.Error(err))
			return nil, NewRepositoryError(err, "failed to scan family row", "SQLITE_ERROR")
		}

		// Parse parents JSON
		var parentDTOs []entity.ParentDTO
		if err := json.Unmarshal([]byte(parentsData), &parentDTOs); err != nil {
			r.logger.Error(ctx, "Failed to unmarshal parents data",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
		}

		// Convert parent DTOs to domain entities
		parents := make([]*entity.Parent, 0, len(parentDTOs))
		for _, dto := range parentDTOs {
			p, err := entity.ParentFromDTO(dto)
			if err != nil {
				r.logger.Error(ctx, "Failed to convert parent DTO to entity",
					zap.Error(err),
					zap.String("family_id", famID),
					zap.String("parent_id", dto.ID))
				return nil, NewRepositoryError(err, "failed to convert parent DTO to entity", "CONVERSION_ERROR")
			}
			parents = append(parents, p)
		}

		// Parse children JSON
		var childDTOs []entity.ChildDTO
		if err := json.Unmarshal([]byte(childrenData), &childDTOs); err != nil {
			r.logger.Error(ctx, "Failed to unmarshal children data",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
		}

		// Convert child DTOs to domain entities
		children := make([]*entity.Child, 0, len(childDTOs))
		for _, dto := range childDTOs {
			c, err := entity.ChildFromDTO(dto)
			if err != nil {
				r.logger.Error(ctx, "Failed to convert child DTO to entity",
					zap.Error(err),
					zap.String("family_id", famID),
					zap.String("child_id", dto.ID))
				return nil, NewRepositoryError(err, "failed to convert child DTO to entity", "CONVERSION_ERROR")
			}
			children = append(children, c)
		}

		// Create family entity
		fam, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
		if err != nil {
			r.logger.Error(ctx, "Failed to create family entity",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to create family entity", "CONVERSION_ERROR")
		}

		r.logger.Debug(ctx, "Retrieved family",
			zap.String("family_id", famID),
			zap.String("status", statusStr))
		families = append(families, fam)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, "Error iterating over family rows", zap.Error(err))
		return nil, NewRepositoryError(err, "error iterating over family rows", "SQLITE_ERROR")
	}

	r.logger.Info(ctx, "Successfully retrieved all families", zap.Int("family_count", len(families)))
	return families, nil
}

// FindByChildID finds the family that contains a specific child
func (r *SQLiteFamilyRepository) FindByChildID(ctx context.Context, childID string) (*entity.Family, error) {
	r.logger.Debug(ctx, "Finding family by child ID in SQLite", zap.String("child_id", childID))

	if childID == "" {
		r.logger.Warn(ctx, "Child ID is required for FindByChildID")
		return nil, errors.NewValidationError("child ID is required", "childID", nil)
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	// SQLite doesn't have native JSON path operators like PostgreSQL,
	// so we need to fetch all families and filter in application code
	r.logger.Debug(ctx, "Querying all families to filter by child ID", zap.String("child_id", childID))
	rows, err := r.DB.QueryContext(ctx, "SELECT id, status, parents, children FROM families")
	if err != nil {
		r.logger.Error(ctx, "Failed to query families", zap.Error(err))
		return nil, NewRepositoryError(err, "failed to query families", "SQLITE_ERROR")
	}
	defer rows.Close()

	var familiesChecked int

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData string

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			r.logger.Error(ctx, "Failed to scan family row", zap.Error(err))
			return nil, NewRepositoryError(err, "failed to scan family row", "SQLITE_ERROR")
		}

		familiesChecked++

		// Parse children JSON
		var childDTOs []entity.ChildDTO
		if err := json.Unmarshal([]byte(childrenData), &childDTOs); err != nil {
			r.logger.Error(ctx, "Failed to unmarshal children data",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
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

		r.logger.Debug(ctx, "Found family with matching child",
			zap.String("family_id", famID),
			zap.String("child_id", childID))

		// Parse parents JSON
		var parentDTOs []entity.ParentDTO
		if err := json.Unmarshal([]byte(parentsData), &parentDTOs); err != nil {
			r.logger.Error(ctx, "Failed to unmarshal parents data",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
		}

		// Convert parent DTOs to domain entities
		parents := make([]*entity.Parent, 0, len(parentDTOs))
		for _, dto := range parentDTOs {
			p, err := entity.ParentFromDTO(dto)
			if err != nil {
				r.logger.Error(ctx, "Failed to convert parent DTO to entity",
					zap.Error(err),
					zap.String("family_id", famID),
					zap.String("parent_id", dto.ID))
				return nil, NewRepositoryError(err, "failed to convert parent DTO to entity", "CONVERSION_ERROR")
			}
			parents = append(parents, p)
		}

		// Convert child DTOs to domain entities
		children := make([]*entity.Child, 0, len(childDTOs))
		for _, dto := range childDTOs {
			c, err := entity.ChildFromDTO(dto)
			if err != nil {
				r.logger.Error(ctx, "Failed to convert child DTO to entity",
					zap.Error(err),
					zap.String("family_id", famID),
					zap.String("child_id", dto.ID))
				return nil, NewRepositoryError(err, "failed to convert child DTO to entity", "CONVERSION_ERROR")
			}
			children = append(children, c)
		}

		// Create family entity
		fam, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
		if err != nil {
			r.logger.Error(ctx, "Failed to create family entity",
				zap.Error(err),
				zap.String("family_id", famID))
			return nil, NewRepositoryError(err, "failed to create family entity", "CONVERSION_ERROR")
		}

		r.logger.Info(ctx, "Successfully found family by child ID",
			zap.String("child_id", childID),
			zap.String("family_id", famID),
			zap.Int("families_checked", familiesChecked))
		return fam, nil
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, "Error iterating over family rows", zap.Error(err))
		return nil, NewRepositoryError(err, "error iterating over family rows", "SQLITE_ERROR")
	}

	r.logger.Info(ctx, "No family found with child ID",
		zap.String("child_id", childID),
		zap.Int("families_checked", familiesChecked))
	return nil, errors.NewNotFoundError("Family with Child", childID, nil)
}