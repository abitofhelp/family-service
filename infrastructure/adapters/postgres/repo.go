// Copyright (c) 2025 A Bit of Help, Inc.

package postgres

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/family-service/infrastructure/adapters/circuitwrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	repoerrors "github.com/abitofhelp/family-service/infrastructure/adapters/errors"
	"github.com/abitofhelp/family-service/infrastructure/adapters/ratewrapper"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/retry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	// Global configuration instance
	globalConfig     *config.Config
	globalConfigOnce sync.Once
)

// SetGlobalConfig sets the global configuration instance
func SetGlobalConfig(cfg *config.Config) {
	globalConfigOnce.Do(func() {
		globalConfig = cfg
	})
}

// getRetryConfig returns the retry configuration
func getRetryConfig() retry.Config {
	if globalConfig != nil {
		return retry.DefaultConfig().
			WithMaxRetries(globalConfig.Retry.MaxRetries).
			WithInitialBackoff(globalConfig.Retry.InitialBackoff).
			WithMaxBackoff(globalConfig.Retry.MaxBackoff)
	}

	// Fallback to default values if configuration is not available
	return retry.DefaultConfig().
		WithMaxRetries(3).
		WithInitialBackoff(100 * time.Millisecond).
		WithMaxBackoff(1 * time.Second)
}

// NewRepositoryError is a helper function that wraps the common repository error handling
// to maintain backward compatibility.
func NewRepositoryError(err error, message string, code string) error {
	// Map old codes to new codes
	newCode := code
	switch code {
	case "POSTGRES_ERROR":
		newCode = repoerrors.PostgresErrorCode
	case "JSON_ERROR":
		newCode = repoerrors.JSONErrorCode
	case "DATA_FORMAT_ERROR":
		newCode = repoerrors.DataFormatErrorCode
	case "CONVERSION_ERROR":
		newCode = repoerrors.ConversionErrorCode
	}

	return repoerrors.NewRepositoryError(err, message, newCode, "families")
}

// PostgresFamilyRepository implements the ports.FamilyRepository interface for PostgreSQL
type PostgresFamilyRepository struct {
	DB             *pgxpool.Pool
	logger         *logging.ContextLogger
	circuitBreaker *circuit.CircuitBreaker
	rateLimiter    *rate.RateLimiter
}

// Ensure PostgresFamilyRepository implements ports.FamilyRepository
var _ ports.FamilyRepository = (*PostgresFamilyRepository)(nil)

// NewPostgresFamilyRepository creates a new PostgresFamilyRepository
func NewPostgresFamilyRepository(db *pgxpool.Pool, logger *logging.ContextLogger) *PostgresFamilyRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	if logger == nil {
		panic("logger cannot be nil")
	}

	// Get circuit breaker configuration
	var circuitConfig *config.CircuitConfig
	if globalConfig != nil {
		circuitConfig = &globalConfig.Circuit
	} else {
		// Default configuration if global config is not available
		circuitConfig = &config.CircuitConfig{
			Enabled:         true,
			Timeout:         5 * time.Second,
			MaxConcurrent:   100,
			ErrorThreshold:  0.5,
			VolumeThreshold: 20,
			SleepWindow:     10 * time.Second,
		}
	}

	// Get rate limiter configuration
	var rateConfig *config.RateConfig
	if globalConfig != nil {
		rateConfig = &globalConfig.Rate
	} else {
		// Default configuration if global config is not available
		rateConfig = &config.RateConfig{
			Enabled:           true,
			RequestsPerSecond: 100,
			BurstSize:         50,
		}
	}

	// Create a new zap logger for the circuit breaker and rate limiter
	zapLogger := zap.NewExample()

	// Create circuit breaker using the family-service wrapper
	cb := circuit.NewCircuitBreaker("postgres", circuitConfig, zapLogger)

	// Create rate limiter using the family-service wrapper
	rl := rate.NewRateLimiter("postgres", rateConfig, zapLogger)

	return &PostgresFamilyRepository{
		DB:             db,
		logger:         logger,
		circuitBreaker: cb,
		rateLimiter:    rl,
	}
}

// ensureTableExists creates the families table if it doesn't exist
func (r *PostgresFamilyRepository) ensureTableExists(ctx context.Context) error {
	r.logger.Debug(ctx, "Ensuring families table exists in PostgreSQL")

	query := `
	CREATE TABLE IF NOT EXISTS families (
		id VARCHAR(36) PRIMARY KEY,
		status VARCHAR(20) NOT NULL,
		parents JSONB NOT NULL,
		children JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes if they don't exist
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_families_status') THEN
			CREATE INDEX idx_families_status ON families(status);
		END IF;

		IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_families_parents') THEN
			CREATE INDEX idx_families_parents ON families USING GIN (parents);
		END IF;

		IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_families_children') THEN
			CREATE INDEX idx_families_children ON families USING GIN (children);
		END IF;
	END
	$$;

	-- Create update trigger function if it doesn't exist
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	-- Create trigger if it doesn't exist
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_families_updated_at') THEN
			CREATE TRIGGER update_families_updated_at
			BEFORE UPDATE ON families
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
		END IF;
	END
	$$;
	`
	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		r.logger.Error(ctx, "Failed to create families table in PostgreSQL", zap.Error(err))
		return NewRepositoryError(err, "failed to create families table", "POSTGRES_ERROR")
	}

	r.logger.Debug(ctx, "Families table exists in PostgreSQL")
	return nil
}

// GetByID retrieves a family by its ID
func (r *PostgresFamilyRepository) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	r.logger.Debug(ctx, "Getting family by ID from PostgreSQL", zap.String("family_id", id))

	if id == "" {
		r.logger.Warn(ctx, "Family ID is required for GetByID")
		return nil, errors.NewValidationError("id is required", "id", nil)
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var famID string
	var statusStr string
	var parentsData, childrenData []byte
	var retryErr error

	// Define the operation to retry
	operation := func(ctx context.Context) error {
		err := r.DB.QueryRow(ctx, `
			SELECT id, status, parents, children FROM families WHERE id = $1
		`, id).Scan(&famID, &statusStr, &parentsData, &childrenData)

		if err != nil {
			if err == pgx.ErrNoRows {
				r.logger.Info(ctx, "Family not found in PostgreSQL", zap.String("family_id", id))
				return errors.NewNotFoundError("Family", id, nil)
			}
			r.logger.Error(ctx, "Failed to get family from PostgreSQL", zap.Error(err), zap.String("family_id", id))
			return NewRepositoryError(err, "failed to get family from PostgreSQL", "POSTGRES_ERROR")
		}
		return nil
	}

	// Define which errors are retryable
	isRetryable := func(err error) bool {
		// Don't retry not found errors
		if _, ok := err.(*errors.NotFoundError); ok {
			return false
		}

		// Don't retry validation errors
		if _, ok := err.(*errors.ValidationError); ok {
			return false
		}

		// Retry network errors, timeouts, and transient database errors
		return retry.IsNetworkError(err) || retry.IsTimeoutError(err) || retry.IsTransientError(err)
	}

	// Configure retry with backoff
	retryConfig := getRetryConfig()

	// Wrap the retry operation with circuit breaker
	circuitOperation := func(ctx context.Context) error {
		// Execute with retry
		retryErr = retry.Do(ctx, operation, retryConfig, isRetryable)
		return retryErr
	}

	// Wrap the circuit breaker operation with rate limiter
	rateOperation := func(ctx context.Context) error {
		// Execute with circuit breaker
		// We need to wrap the circuitOperation to match the generic function signature
		circuitOpWrapper := func(ctx context.Context) (bool, error) {
			err := circuitOperation(ctx)
			return err == nil, err
		}
		_, err := circuit.Execute(ctx, r.circuitBreaker, "GetByID", circuitOpWrapper)
		return err
	}

	// Execute with rate limiter
	// We need to wrap the rateOperation to match the generic function signature
	rateOpWrapper := func(ctx context.Context) (bool, error) {
		err := rateOperation(ctx)
		return err == nil, err
	}
	_, err := rate.Execute(ctxWithTimeout, r.rateLimiter, "GetByID", rateOpWrapper)

	// Check for errors from rate limiter or circuit breaker
	if err != nil && retryErr == nil {
		// Check if it's a rate limiter error
		if strings.Contains(err.Error(), "rate limit exceeded") {
			return nil, NewRepositoryError(err, "rate limit exceeded", "POSTGRES_ERROR")
		}
		// Otherwise, assume it's a circuit breaker error
		return nil, NewRepositoryError(err, "circuit breaker is open", "POSTGRES_ERROR")
	}

	// Handle retry errors
	if retryErr != nil {
		// If it's already a typed error, return it directly
		if _, ok := retryErr.(*errors.NotFoundError); ok {
			return nil, retryErr
		}
		if _, ok := retryErr.(*errors.ValidationError); ok {
			return nil, retryErr
		}
		if _, ok := retryErr.(*errors.DatabaseError); ok {
			return nil, retryErr
		}

		// Otherwise, wrap it in a database error
		return nil, NewRepositoryError(retryErr, "failed to get family from PostgreSQL after retries", "POSTGRES_ERROR")
	}

	r.logger.Debug(ctx, "Successfully retrieved family data from PostgreSQL", zap.String("family_id", id))

	// Define custom structs for JSON unmarshaling to handle both uppercase and lowercase field names
	type jsonParent struct {
		ID        string  `json:"ID,omitempty"`
		Id        string  `json:"id,omitempty"`
		FirstName string  `json:"FirstName,omitempty"`
		FirstN    string  `json:"firstName,omitempty"`
		LastName  string  `json:"LastName,omitempty"`
		LastN     string  `json:"lastName,omitempty"`
		BirthDate string  `json:"BirthDate,omitempty"`
		BirthD    string  `json:"birthDate,omitempty"`
		DeathDate *string `json:"DeathDate,omitempty"`
		DeathD    *string `json:"deathDate,omitempty"`
	}

	type jsonChild struct {
		ID        string  `json:"ID,omitempty"`
		Id        string  `json:"id,omitempty"`
		FirstName string  `json:"FirstName,omitempty"`
		FirstN    string  `json:"firstName,omitempty"`
		LastName  string  `json:"LastName,omitempty"`
		LastN     string  `json:"lastName,omitempty"`
		BirthDate string  `json:"BirthDate,omitempty"`
		BirthD    string  `json:"birthDate,omitempty"`
		DeathDate *string `json:"DeathDate,omitempty"`
		DeathD    *string `json:"deathDate,omitempty"`
	}

	// Parse parents JSON
	var jsonParents []jsonParent
	if err := json.Unmarshal(parentsData, &jsonParents); err != nil {
		return nil, NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
	}

	// Convert JSON parents to domain entities
	parents := make([]*entity.Parent, 0, len(jsonParents))
	for _, jp := range jsonParents {
		// Use the appropriate field based on which one is populated
		id := jp.ID
		if id == "" {
			id = jp.Id
		}
		firstName := jp.FirstName
		if firstName == "" {
			firstName = jp.FirstN
		}
		lastName := jp.LastName
		if lastName == "" {
			lastName = jp.LastN
		}
		birthDateStr := jp.BirthDate
		if birthDateStr == "" {
			birthDateStr = jp.BirthD
		}

		birthDate, err := time.Parse(time.RFC3339, birthDateStr)
		if err != nil {
			return nil, NewRepositoryError(err, "invalid parent birth date format", "DATA_FORMAT_ERROR")
		}

		var deathDate *time.Time
		deathDateStr := jp.DeathDate
		if deathDateStr == nil {
			deathDateStr = jp.DeathD
		}
		if deathDateStr != nil {
			parsedDeathDate, err := time.Parse(time.RFC3339, *deathDateStr)
			if err != nil {
				return nil, NewRepositoryError(err, "invalid parent death date format", "DATA_FORMAT_ERROR")
			}
			deathDate = &parsedDeathDate
		}

		p, err := entity.NewParent(id, firstName, lastName, birthDate, deathDate)
		if err != nil {
			return nil, NewRepositoryError(err, "failed to create parent entity", "CONVERSION_ERROR")
		}
		parents = append(parents, p)
	}

	// Parse children JSON
	var jsonChildren []jsonChild
	if err := json.Unmarshal(childrenData, &jsonChildren); err != nil {
		return nil, NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
	}

	// Convert JSON children to domain entities
	children := make([]*entity.Child, 0, len(jsonChildren))
	for _, jc := range jsonChildren {
		// Use the appropriate field based on which one is populated
		id := jc.ID
		if id == "" {
			id = jc.Id
		}
		firstName := jc.FirstName
		if firstName == "" {
			firstName = jc.FirstN
		}
		lastName := jc.LastName
		if lastName == "" {
			lastName = jc.LastN
		}
		birthDateStr := jc.BirthDate
		if birthDateStr == "" {
			birthDateStr = jc.BirthD
		}

		birthDate, err := time.Parse(time.RFC3339, birthDateStr)
		if err != nil {
			return nil, NewRepositoryError(err, "invalid child birth date format", "DATA_FORMAT_ERROR")
		}

		var deathDate *time.Time
		deathDateStr := jc.DeathDate
		if deathDateStr == nil {
			deathDateStr = jc.DeathD
		}
		if deathDateStr != nil {
			parsedDeathDate, err := time.Parse(time.RFC3339, *deathDateStr)
			if err != nil {
				return nil, NewRepositoryError(err, "invalid child death date format", "DATA_FORMAT_ERROR")
			}
			deathDate = &parsedDeathDate
		}

		c, err := entity.NewChild(id, firstName, lastName, birthDate, deathDate)
		if err != nil {
			return nil, NewRepositoryError(err, "failed to create child entity", "CONVERSION_ERROR")
		}
		children = append(children, c)
	}

	// Create family entity
	return entity.NewFamily(famID, entity.Status(statusStr), parents, children)
}

// Save persists a family
func (r *PostgresFamilyRepository) Save(ctx context.Context, fam *entity.Family) error {
	if fam == nil {
		r.logger.Warn(ctx, "Family cannot be nil for Save")
		return errors.NewValidationError("family cannot be nil", "family", nil)
	}

	r.logger.Debug(ctx, "Saving family to PostgreSQL", zap.String("family_id", fam.ID()))

	if err := fam.Validate(); err != nil {
		r.logger.Error(ctx, "Family validation failed", zap.Error(err), zap.String("family_id", fam.ID()))
		return err
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return err
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		r.logger.Error(ctx, "Failed to begin transaction", zap.Error(err), zap.String("family_id", fam.ID()))
		return NewRepositoryError(err, "failed to begin transaction", "POSTGRES_ERROR")
	}

	var txErr error
	defer func() {
		if txErr != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				// Log rollback error, but don't return it as it would mask the original error
				// In a real implementation, we'd use a logger here
			}
		}
	}()

	// Create custom JSON-compatible structures for parents and children
	// to ensure proper date formatting
	type jsonParent struct {
		ID        string  `json:"id"`
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		BirthDate string  `json:"birthDate"`
		DeathDate *string `json:"deathDate,omitempty"`
	}

	type jsonChild struct {
		ID        string  `json:"id"`
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		BirthDate string  `json:"birthDate"`
		DeathDate *string `json:"deathDate,omitempty"`
	}

	// Convert parents to JSON-compatible format
	jsonParents := make([]jsonParent, 0, len(fam.Parents()))
	for _, p := range fam.Parents() {
		var deathDateStr *string
		if p.DeathDate() != nil {
			str := p.DeathDate().Format(time.RFC3339)
			deathDateStr = &str
		}

		jsonParents = append(jsonParents, jsonParent{
			ID:        p.ID(),
			FirstName: p.FirstName(),
			LastName:  p.LastName(),
			BirthDate: p.BirthDate().Format(time.RFC3339),
			DeathDate: deathDateStr,
		})
	}

	// Convert children to JSON-compatible format
	jsonChildren := make([]jsonChild, 0, len(fam.Children()))
	for _, c := range fam.Children() {
		var deathDateStr *string
		if c.DeathDate() != nil {
			str := c.DeathDate().Format(time.RFC3339)
			deathDateStr = &str
		}

		jsonChildren = append(jsonChildren, jsonChild{
			ID:        c.ID(),
			FirstName: c.FirstName(),
			LastName:  c.LastName(),
			BirthDate: c.BirthDate().Format(time.RFC3339),
			DeathDate: deathDateStr,
		})
	}

	// Marshal to JSON
	parentsJSON, err := json.Marshal(jsonParents)
	if err != nil {
		return NewRepositoryError(err, "failed to marshal parents to JSON", "JSON_ERROR")
	}

	childrenJSON, err := json.Marshal(jsonChildren)
	if err != nil {
		return NewRepositoryError(err, "failed to marshal children to JSON", "JSON_ERROR")
	}

	// Validate that the JSON is valid
	if !json.Valid(parentsJSON) {
		return NewRepositoryError(nil, "invalid parents JSON", "JSON_ERROR")
	}

	if !json.Valid(childrenJSON) {
		return NewRepositoryError(nil, "invalid children JSON", "JSON_ERROR")
	}

	// Execute SQL
	_, txErr = tx.Exec(ctx, `
        INSERT INTO families (id, status, parents, children)
        VALUES ($1, $2, $3::jsonb, $4::jsonb)
        ON CONFLICT (id) DO UPDATE SET
            status = EXCLUDED.status,
            parents = EXCLUDED.parents,
            children = EXCLUDED.children
    `, fam.ID(), string(fam.Status()), parentsJSON, childrenJSON)

	if txErr != nil {
		return NewRepositoryError(txErr, "failed to save family to PostgreSQL", "POSTGRES_ERROR")
	}

	// Commit transaction
	if txErr = tx.Commit(ctx); txErr != nil {
		return NewRepositoryError(txErr, "failed to commit transaction", "POSTGRES_ERROR")
	}

	return nil
}

// FindByParentID finds families that contain a specific parent
func (r *PostgresFamilyRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error) {
	r.logger.Debug(ctx, "Finding families by parent ID in PostgreSQL", zap.String("parent_id", parentID))

	if parentID == "" {
		r.logger.Warn(ctx, "Parent ID is required for FindByParentID")
		return nil, errors.NewValidationError("parent ID is required", "parentID", nil)
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	// Query for both uppercase and lowercase ID fields
	rows, err := r.DB.Query(ctx, `
        SELECT id, status, parents, children FROM families 
        WHERE parents @> ANY (ARRAY[jsonb_build_array(jsonb_build_object('id', $1))]) 
        OR parents @> ANY (ARRAY[jsonb_build_array(jsonb_build_object('ID', $1))])
    `, parentID)

	if err != nil {
		return nil, NewRepositoryError(err, "failed to find families by parent ID", "POSTGRES_ERROR")
	}
	defer rows.Close()

	var families []*entity.Family

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData []byte

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			return nil, NewRepositoryError(err, "failed to scan family row", "POSTGRES_ERROR")
		}

		// Define custom structs for JSON unmarshaling to handle both uppercase and lowercase field names
		type jsonParent struct {
			ID        string  `json:"ID,omitempty"`
			Id        string  `json:"id,omitempty"`
			FirstName string  `json:"FirstName,omitempty"`
			FirstN    string  `json:"firstName,omitempty"`
			LastName  string  `json:"LastName,omitempty"`
			LastN     string  `json:"lastName,omitempty"`
			BirthDate string  `json:"BirthDate,omitempty"`
			BirthD    string  `json:"birthDate,omitempty"`
			DeathDate *string `json:"DeathDate,omitempty"`
			DeathD    *string `json:"deathDate,omitempty"`
		}

		type jsonChild struct {
			ID        string  `json:"ID,omitempty"`
			Id        string  `json:"id,omitempty"`
			FirstName string  `json:"FirstName,omitempty"`
			FirstN    string  `json:"firstName,omitempty"`
			LastName  string  `json:"LastName,omitempty"`
			LastN     string  `json:"lastName,omitempty"`
			BirthDate string  `json:"BirthDate,omitempty"`
			BirthD    string  `json:"birthDate,omitempty"`
			DeathDate *string `json:"DeathDate,omitempty"`
			DeathD    *string `json:"deathDate,omitempty"`
		}

		// Parse parents JSON
		var jsonParents []jsonParent
		if err := json.Unmarshal(parentsData, &jsonParents); err != nil {
			return nil, NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
		}

		// Convert JSON parents to domain entities
		parents := make([]*entity.Parent, 0, len(jsonParents))
		for _, jp := range jsonParents {
			// Use the appropriate field based on which one is populated
			id := jp.ID
			if id == "" {
				id = jp.Id
			}
			firstName := jp.FirstName
			if firstName == "" {
				firstName = jp.FirstN
			}
			lastName := jp.LastName
			if lastName == "" {
				lastName = jp.LastN
			}
			birthDateStr := jp.BirthDate
			if birthDateStr == "" {
				birthDateStr = jp.BirthD
			}

			birthDate, err := time.Parse(time.RFC3339, birthDateStr)
			if err != nil {
				return nil, NewRepositoryError(err, "invalid parent birth date format", "DATA_FORMAT_ERROR")
			}

			var deathDate *time.Time
			deathDateStr := jp.DeathDate
			if deathDateStr == nil {
				deathDateStr = jp.DeathD
			}
			if deathDateStr != nil {
				parsedDeathDate, err := time.Parse(time.RFC3339, *deathDateStr)
				if err != nil {
					return nil, NewRepositoryError(err, "invalid parent death date format", "DATA_FORMAT_ERROR")
				}
				deathDate = &parsedDeathDate
			}

			p, err := entity.NewParent(id, firstName, lastName, birthDate, deathDate)
			if err != nil {
				return nil, NewRepositoryError(err, "failed to create parent entity", "CONVERSION_ERROR")
			}
			parents = append(parents, p)
		}

		// Parse children JSON
		var jsonChildren []jsonChild
		if err := json.Unmarshal(childrenData, &jsonChildren); err != nil {
			return nil, NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
		}

		// Convert JSON children to domain entities
		children := make([]*entity.Child, 0, len(jsonChildren))
		for _, jc := range jsonChildren {
			// Use the appropriate field based on which one is populated
			id := jc.ID
			if id == "" {
				id = jc.Id
			}
			firstName := jc.FirstName
			if firstName == "" {
				firstName = jc.FirstN
			}
			lastName := jc.LastName
			if lastName == "" {
				lastName = jc.LastN
			}
			birthDateStr := jc.BirthDate
			if birthDateStr == "" {
				birthDateStr = jc.BirthD
			}

			birthDate, err := time.Parse(time.RFC3339, birthDateStr)
			if err != nil {
				return nil, NewRepositoryError(err, "invalid child birth date format", "DATA_FORMAT_ERROR")
			}

			var deathDate *time.Time
			deathDateStr := jc.DeathDate
			if deathDateStr == nil {
				deathDateStr = jc.DeathD
			}
			if deathDateStr != nil {
				parsedDeathDate, err := time.Parse(time.RFC3339, *deathDateStr)
				if err != nil {
					return nil, NewRepositoryError(err, "invalid child death date format", "DATA_FORMAT_ERROR")
				}
				deathDate = &parsedDeathDate
			}

			c, err := entity.NewChild(id, firstName, lastName, birthDate, deathDate)
			if err != nil {
				return nil, NewRepositoryError(err, "failed to create child entity", "CONVERSION_ERROR")
			}
			children = append(children, c)
		}

		// Create family entity
		fam, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
		if err != nil {
			return nil, NewRepositoryError(err, "failed to create family entity", "CONVERSION_ERROR")
		}

		families = append(families, fam)
	}

	if err := rows.Err(); err != nil {
		return nil, NewRepositoryError(err, "error iterating over family rows", "POSTGRES_ERROR")
	}

	return families, nil
}

// FindByChildID finds the family that contains a specific child
func (r *PostgresFamilyRepository) FindByChildID(ctx context.Context, childID string) (*entity.Family, error) {
	r.logger.Debug(ctx, "Finding family by child ID in PostgreSQL", zap.String("child_id", childID))

	if childID == "" {
		r.logger.Warn(ctx, "Child ID is required for FindByChildID")
		return nil, errors.NewValidationError("child ID is required", "childID", nil)
	}

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	var famID string
	var statusStr string
	var parentsData, childrenData []byte

	// Query for both uppercase and lowercase ID fields
	err := r.DB.QueryRow(ctx, `
        SELECT id, status, parents, children FROM families 
        WHERE children @> ANY (ARRAY[jsonb_build_array(jsonb_build_object('id', $1))])
        OR children @> ANY (ARRAY[jsonb_build_array(jsonb_build_object('ID', $1))])
    `, childID).Scan(&famID, &statusStr, &parentsData, &childrenData)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NewNotFoundError("Family with Child", childID, nil)
		}
		return nil, NewRepositoryError(err, "failed to find family by child ID", "POSTGRES_ERROR")
	}

	// Define custom structs for JSON unmarshaling to handle both uppercase and lowercase field names
	type jsonParent struct {
		ID        string  `json:"ID,omitempty"`
		Id        string  `json:"id,omitempty"`
		FirstName string  `json:"FirstName,omitempty"`
		FirstN    string  `json:"firstName,omitempty"`
		LastName  string  `json:"LastName,omitempty"`
		LastN     string  `json:"lastName,omitempty"`
		BirthDate string  `json:"BirthDate,omitempty"`
		BirthD    string  `json:"birthDate,omitempty"`
		DeathDate *string `json:"DeathDate,omitempty"`
		DeathD    *string `json:"deathDate,omitempty"`
	}

	type jsonChild struct {
		ID        string  `json:"ID,omitempty"`
		Id        string  `json:"id,omitempty"`
		FirstName string  `json:"FirstName,omitempty"`
		FirstN    string  `json:"firstName,omitempty"`
		LastName  string  `json:"LastName,omitempty"`
		LastN     string  `json:"lastName,omitempty"`
		BirthDate string  `json:"BirthDate,omitempty"`
		BirthD    string  `json:"birthDate,omitempty"`
		DeathDate *string `json:"DeathDate,omitempty"`
		DeathD    *string `json:"deathDate,omitempty"`
	}

	// Parse parents JSON
	var jsonParents []jsonParent
	if err := json.Unmarshal(parentsData, &jsonParents); err != nil {
		return nil, NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
	}

	// Convert JSON parents to domain entities
	parents := make([]*entity.Parent, 0, len(jsonParents))
	for _, jp := range jsonParents {
		// Use the appropriate field based on which one is populated
		id := jp.ID
		if id == "" {
			id = jp.Id
		}
		firstName := jp.FirstName
		if firstName == "" {
			firstName = jp.FirstN
		}
		lastName := jp.LastName
		if lastName == "" {
			lastName = jp.LastN
		}
		birthDateStr := jp.BirthDate
		if birthDateStr == "" {
			birthDateStr = jp.BirthD
		}

		birthDate, err := time.Parse(time.RFC3339, birthDateStr)
		if err != nil {
			return nil, NewRepositoryError(err, "invalid parent birth date format", "DATA_FORMAT_ERROR")
		}

		var deathDate *time.Time
		deathDateStr := jp.DeathDate
		if deathDateStr == nil {
			deathDateStr = jp.DeathD
		}
		if deathDateStr != nil {
			parsedDeathDate, err := time.Parse(time.RFC3339, *deathDateStr)
			if err != nil {
				return nil, NewRepositoryError(err, "invalid parent death date format", "DATA_FORMAT_ERROR")
			}
			deathDate = &parsedDeathDate
		}

		p, err := entity.NewParent(id, firstName, lastName, birthDate, deathDate)
		if err != nil {
			return nil, NewRepositoryError(err, "failed to create parent entity", "CONVERSION_ERROR")
		}
		parents = append(parents, p)
	}

	// Parse children JSON
	var jsonChildren []jsonChild
	if err := json.Unmarshal(childrenData, &jsonChildren); err != nil {
		return nil, NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
	}

	// Convert JSON children to domain entities
	children := make([]*entity.Child, 0, len(jsonChildren))
	for _, jc := range jsonChildren {
		// Use the appropriate field based on which one is populated
		id := jc.ID
		if id == "" {
			id = jc.Id
		}
		firstName := jc.FirstName
		if firstName == "" {
			firstName = jc.FirstN
		}
		lastName := jc.LastName
		if lastName == "" {
			lastName = jc.LastN
		}
		birthDateStr := jc.BirthDate
		if birthDateStr == "" {
			birthDateStr = jc.BirthD
		}

		birthDate, err := time.Parse(time.RFC3339, birthDateStr)
		if err != nil {
			return nil, NewRepositoryError(err, "invalid child birth date format", "DATA_FORMAT_ERROR")
		}

		var deathDate *time.Time
		deathDateStr := jc.DeathDate
		if deathDateStr == nil {
			deathDateStr = jc.DeathD
		}
		if deathDateStr != nil {
			parsedDeathDate, err := time.Parse(time.RFC3339, *deathDateStr)
			if err != nil {
				return nil, NewRepositoryError(err, "invalid child death date format", "DATA_FORMAT_ERROR")
			}
			deathDate = &parsedDeathDate
		}

		c, err := entity.NewChild(id, firstName, lastName, birthDate, deathDate)
		if err != nil {
			return nil, NewRepositoryError(err, "failed to create child entity", "CONVERSION_ERROR")
		}
		children = append(children, c)
	}

	// Create family entity
	return entity.NewFamily(famID, entity.Status(statusStr), parents, children)
}

// GetAll retrieves all families
func (r *PostgresFamilyRepository) GetAll(ctx context.Context) ([]*entity.Family, error) {
	r.logger.Debug(ctx, "Getting all families from PostgreSQL")

	// Ensure table exists
	if err := r.ensureTableExists(ctx); err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, `
        SELECT id, status, parents, children FROM families
    `)

	if err != nil {
		return nil, NewRepositoryError(err, "failed to get all families", "POSTGRES_ERROR")
	}
	defer rows.Close()

	var families []*entity.Family

	for rows.Next() {
		var famID string
		var statusStr string
		var parentsData, childrenData []byte

		if err := rows.Scan(&famID, &statusStr, &parentsData, &childrenData); err != nil {
			return nil, NewRepositoryError(err, "failed to scan family row", "POSTGRES_ERROR")
		}

		// Define custom structs for JSON unmarshaling to handle both uppercase and lowercase field names
		type jsonParent struct {
			ID        string  `json:"ID,omitempty"`
			Id        string  `json:"id,omitempty"`
			FirstName string  `json:"FirstName,omitempty"`
			FirstN    string  `json:"firstName,omitempty"`
			LastName  string  `json:"LastName,omitempty"`
			LastN     string  `json:"lastName,omitempty"`
			BirthDate string  `json:"BirthDate,omitempty"`
			BirthD    string  `json:"birthDate,omitempty"`
			DeathDate *string `json:"DeathDate,omitempty"`
			DeathD    *string `json:"deathDate,omitempty"`
		}

		type jsonChild struct {
			ID        string  `json:"ID,omitempty"`
			Id        string  `json:"id,omitempty"`
			FirstName string  `json:"FirstName,omitempty"`
			FirstN    string  `json:"firstName,omitempty"`
			LastName  string  `json:"LastName,omitempty"`
			LastN     string  `json:"lastName,omitempty"`
			BirthDate string  `json:"BirthDate,omitempty"`
			BirthD    string  `json:"birthDate,omitempty"`
			DeathDate *string `json:"DeathDate,omitempty"`
			DeathD    *string `json:"deathDate,omitempty"`
		}

		// Parse parents JSON
		var jsonParents []jsonParent
		if err := json.Unmarshal(parentsData, &jsonParents); err != nil {
			return nil, NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")
		}

		// Convert JSON parents to domain entities
		parents := make([]*entity.Parent, 0, len(jsonParents))
		for _, jp := range jsonParents {
			// Use the appropriate field based on which one is populated
			id := jp.ID
			if id == "" {
				id = jp.Id
			}
			firstName := jp.FirstName
			if firstName == "" {
				firstName = jp.FirstN
			}
			lastName := jp.LastName
			if lastName == "" {
				lastName = jp.LastN
			}
			birthDateStr := jp.BirthDate
			if birthDateStr == "" {
				birthDateStr = jp.BirthD
			}

			birthDate, err := time.Parse(time.RFC3339, birthDateStr)
			if err != nil {
				return nil, NewRepositoryError(err, "invalid parent birth date format", "DATA_FORMAT_ERROR")
			}

			var deathDate *time.Time
			deathDateStr := jp.DeathDate
			if deathDateStr == nil {
				deathDateStr = jp.DeathD
			}
			if deathDateStr != nil {
				parsedDeathDate, err := time.Parse(time.RFC3339, *deathDateStr)
				if err != nil {
					return nil, NewRepositoryError(err, "invalid parent death date format", "DATA_FORMAT_ERROR")
				}
				deathDate = &parsedDeathDate
			}

			p, err := entity.NewParent(id, firstName, lastName, birthDate, deathDate)
			if err != nil {
				return nil, NewRepositoryError(err, "failed to create parent entity", "CONVERSION_ERROR")
			}
			parents = append(parents, p)
		}

		// Parse children JSON
		var jsonChildren []jsonChild
		if err := json.Unmarshal(childrenData, &jsonChildren); err != nil {
			return nil, NewRepositoryError(err, "failed to unmarshal children data", "JSON_ERROR")
		}

		// Convert JSON children to domain entities
		children := make([]*entity.Child, 0, len(jsonChildren))
		for _, jc := range jsonChildren {
			// Use the appropriate field based on which one is populated
			id := jc.ID
			if id == "" {
				id = jc.Id
			}
			firstName := jc.FirstName
			if firstName == "" {
				firstName = jc.FirstN
			}
			lastName := jc.LastName
			if lastName == "" {
				lastName = jc.LastN
			}
			birthDateStr := jc.BirthDate
			if birthDateStr == "" {
				birthDateStr = jc.BirthD
			}

			birthDate, err := time.Parse(time.RFC3339, birthDateStr)
			if err != nil {
				return nil, NewRepositoryError(err, "invalid child birth date format", "DATA_FORMAT_ERROR")
			}

			var deathDate *time.Time
			deathDateStr := jc.DeathDate
			if deathDateStr == nil {
				deathDateStr = jc.DeathD
			}
			if deathDateStr != nil {
				parsedDeathDate, err := time.Parse(time.RFC3339, *deathDateStr)
				if err != nil {
					return nil, NewRepositoryError(err, "invalid child death date format", "DATA_FORMAT_ERROR")
				}
				deathDate = &parsedDeathDate
			}

			c, err := entity.NewChild(id, firstName, lastName, birthDate, deathDate)
			if err != nil {
				return nil, NewRepositoryError(err, "failed to create child entity", "CONVERSION_ERROR")
			}
			children = append(children, c)
		}

		// Create family entity
		fam, err := entity.NewFamily(famID, entity.Status(statusStr), parents, children)
		if err != nil {
			return nil, NewRepositoryError(err, "failed to create family entity", "CONVERSION_ERROR")
		}

		families = append(families, fam)
	}

	if err := rows.Err(); err != nil {
		return nil, NewRepositoryError(err, "error iterating over family rows", "POSTGRES_ERROR")
	}

	return families, nil
}
