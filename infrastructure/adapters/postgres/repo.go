package postgres

import (
	"context"
	"encoding/json"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresFamilyRepository implements the ports.FamilyRepository interface for PostgreSQL
type PostgresFamilyRepository struct {
	DB *pgxpool.Pool
}

// Ensure PostgresFamilyRepository implements ports.FamilyRepository
var _ ports.FamilyRepository = (*PostgresFamilyRepository)(nil)

// NewPostgresFamilyRepository creates a new PostgresFamilyRepository
func NewPostgresFamilyRepository(db *pgxpool.Pool) *PostgresFamilyRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &PostgresFamilyRepository{DB: db}
}

// GetByID retrieves a family by its ID
func (r *PostgresFamilyRepository) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	if id == "" {
		return nil, errors.NewValidationError("id is required")
	}

	var famID string
	var statusStr string
	var parentsData, childrenData []byte

	err := r.DB.QueryRow(ctx, `
        SELECT id, status, parents, children FROM families WHERE id = $1
    `, id).Scan(&famID, &statusStr, &parentsData, &childrenData)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NewNotFoundError("Family", id)
		}
		return nil, errors.NewRepositoryError(err, "failed to get family from PostgreSQL", "POSTGRES_ERROR")
	}

	// Parse parents JSON
	var parentDTOs []entity.ParentDTO
	if err := json.Unmarshal(parentsData, &parentDTOs); err != nil {
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
	if err := json.Unmarshal(childrenData, &childDTOs); err != nil {
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
func (r *PostgresFamilyRepository) Save(ctx context.Context, fam *entity.Family) error {
	if fam == nil {
		return errors.NewValidationError("family cannot be nil")
	}

	if err := fam.Validate(); err != nil {
		return err
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return errors.NewRepositoryError(err, "failed to begin transaction", "POSTGRES_ERROR")
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				// Log rollback error, but don't return it as it would mask the original error
				// In a real implementation, we'd use a logger here
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

	// Execute SQL
	_, err = tx.Exec(ctx, `
        INSERT INTO families (id, status, parents, children)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET
            status = EXCLUDED.status,
            parents = EXCLUDED.parents,
            children = EXCLUDED.children
    `, fam.ID(), string(fam.Status()), parentsJSON, childrenJSON)

	if err != nil {
		return errors.NewRepositoryError(err, "failed to save family to PostgreSQL", "POSTGRES_ERROR")
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return errors.NewRepositoryError(err, "failed to commit transaction", "POSTGRES_ERROR")
	}

	return nil
}

// FindByParentID finds families that contain a specific parent
func (r *PostgresFamilyRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error) {
	if parentID == "" {
		return nil, errors.NewValidationError("parent ID is required")
	}

	rows, err := r.DB.Query(ctx, `
        SELECT id, status, parents, children FROM families 
        WHERE parents @> ANY (ARRAY[jsonb_build_array(jsonb_build_object('id', $1))])
    `, parentID)

	if err != nil {
		return nil, errors.NewRepositoryError(err, "failed to find families by parent ID", "POSTGRES_ERROR")
	}
	defer rows.Close()

	var families []*entity.Family

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData []byte

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to scan family row", "POSTGRES_ERROR")
		}

		// Parse parents JSON
		var parentDTOs []entity.ParentDTO
		if err := json.Unmarshal(parentsData, &parentDTOs); err != nil {
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
		if err := json.Unmarshal(childrenData, &childDTOs); err != nil {
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
		return nil, errors.NewRepositoryError(err, "error iterating over family rows", "POSTGRES_ERROR")
	}

	return families, nil
}

// FindByChildID finds the family that contains a specific child
func (r *PostgresFamilyRepository) FindByChildID(ctx context.Context, childID string) (*entity.Family, error) {
	if childID == "" {
		return nil, errors.NewValidationError("child ID is required")
	}

	var famID string
	var statusStr string
	var parentsData, childrenData []byte

	err := r.DB.QueryRow(ctx, `
        SELECT id, status, parents, children FROM families 
        WHERE children @> ANY (ARRAY[jsonb_build_array(jsonb_build_object('id', $1))])
    `, childID).Scan(&famID, &statusStr, &parentsData, &childrenData)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NewNotFoundError("Family with Child", childID)
		}
		return nil, errors.NewRepositoryError(err, "failed to find family by child ID", "POSTGRES_ERROR")
	}

	// Parse parents JSON
	var parentDTOs []entity.ParentDTO
	if err := json.Unmarshal(parentsData, &parentDTOs); err != nil {
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
	if err := json.Unmarshal(childrenData, &childDTOs); err != nil {
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

// GetAll retrieves all families
func (r *PostgresFamilyRepository) GetAll(ctx context.Context) ([]*entity.Family, error) {
	rows, err := r.DB.Query(ctx, `
        SELECT id, status, parents, children FROM families
    `)

	if err != nil {
		return nil, errors.NewRepositoryError(err, "failed to get all families", "POSTGRES_ERROR")
	}
	defer rows.Close()

	var families []*entity.Family

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData []byte

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			return nil, errors.NewRepositoryError(err, "failed to scan family row", "POSTGRES_ERROR")
		}

		// Parse parents JSON
		var parentDTOs []entity.ParentDTO
		if err := json.Unmarshal(parentsData, &parentDTOs); err != nil {
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
		if err := json.Unmarshal(childrenData, &childDTOs); err != nil {
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
		return nil, errors.NewRepositoryError(err, "error iterating over family rows", "POSTGRES_ERROR")
	}

	return families, nil
}
