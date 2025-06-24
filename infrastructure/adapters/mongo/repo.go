// Copyright (c) 2025 A Bit of Help, Inc.

package mongo

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/servicelib/rate"
	"github.com/abitofhelp/servicelib/circuit"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/retry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// FamilyDocument represents how a family is stored in MongoDB
type FamilyDocument struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	FamilyID string             `bson:"family_id"`
	Status   string             `bson:"status"`
	Parents  []ParentDocument   `bson:"parents"`
	Children []ChildDocument    `bson:"children"`
}

// ParentDocument represents how a parent is stored in MongoDB
type ParentDocument struct {
	ID        string  `bson:"id"`
	FirstName string  `bson:"firstName"`
	LastName  string  `bson:"lastName"`
	BirthDate string  `bson:"birthDate"`
	DeathDate *string `bson:"deathDate,omitempty"`
}

// ChildDocument represents how a child is stored in MongoDB
type ChildDocument struct {
	ID        string  `bson:"id"`
	FirstName string  `bson:"firstName"`
	LastName  string  `bson:"lastName"`
	BirthDate string  `bson:"birthDate"`
	DeathDate *string `bson:"deathDate,omitempty"`
}

// MongoFamilyRepository implements the ports.FamilyRepository interface for MongoDB
type MongoFamilyRepository struct {
	Collection    *mongo.Collection
	logger        *logging.ContextLogger
	circuitBreaker *circuit.CircuitBreaker
	rateLimiter    *rate.RateLimiter
}

// Ensure MongoFamilyRepository implements ports.FamilyRepository
var _ ports.FamilyRepository = (*MongoFamilyRepository)(nil)

// NewMongoFamilyRepository creates a new MongoFamilyRepository
func NewMongoFamilyRepository(collection *mongo.Collection, logger *logging.ContextLogger) *MongoFamilyRepository {
	if collection == nil {
		panic("collection cannot be nil")
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
	zapLogger, _ := zap.NewProduction()
	contextLogger := logging.NewContextLogger(zapLogger)

	// Create circuit breaker config
	circuitBreakerConfig := circuit.DefaultConfig().
		WithEnabled(circuitConfig.Enabled).
		WithTimeout(circuitConfig.Timeout).
		WithMaxConcurrent(circuitConfig.MaxConcurrent).
		WithErrorThreshold(circuitConfig.ErrorThreshold).
		WithVolumeThreshold(circuitConfig.VolumeThreshold).
		WithSleepWindow(circuitConfig.SleepWindow)

	// Create circuit breaker options
	circuitBreakerOptions := circuit.DefaultOptions().
		WithName("mongodb").
		WithLogger(contextLogger)

	// Create circuit breaker
	cb := circuit.NewCircuitBreaker(circuitBreakerConfig, circuitBreakerOptions)

	// Create rate limiter config
	rateLimiterConfig := rate.DefaultConfig().
		WithEnabled(rateConfig.Enabled).
		WithRequestsPerSecond(rateConfig.RequestsPerSecond).
		WithBurstSize(rateConfig.BurstSize)

	// Create rate limiter options
	rateLimiterOptions := rate.DefaultOptions().
		WithName("mongodb").
		WithLogger(contextLogger)

	// Create rate limiter
	rl := rate.NewRateLimiter(rateLimiterConfig, rateLimiterOptions)

	return &MongoFamilyRepository{
		Collection:    collection,
		logger:        logger,
		circuitBreaker: cb,
		rateLimiter:    rl,
	}
}

// GetByID retrieves a family by its ID
func (r *MongoFamilyRepository) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	r.logger.Debug(ctx, "Getting family by ID from MongoDB", zap.String("family_id", id))

	if id == "" {
		r.logger.Warn(ctx, "Family ID is required for GetByID")
		return nil, errors.NewValidationError("id is required", "id", nil)
	}

	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var doc FamilyDocument
	var family *entity.Family
	var retryErr error

	// Define the operation to retry
	operation := func(ctx context.Context) error {
		err := r.Collection.FindOne(ctx, bson.M{"family_id": id}).Decode(&doc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				r.logger.Info(ctx, "Family not found in MongoDB", zap.String("family_id", id))
				return errors.NewNotFoundError("Family", id, nil)
			}
			r.logger.Error(ctx, "Failed to get family from MongoDB", zap.Error(err), zap.String("family_id", id))
			return errors.NewDatabaseError("failed to get family from MongoDB", "query", "families", err)
		}

		// Convert document to domain entity
		var convErr error
		family, convErr = r.documentToEntity(doc)
		return convErr
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
			return nil, errors.NewDatabaseError("rate limit exceeded", "query", "families", err)
		}
		// Otherwise, assume it's a circuit breaker error
		return nil, errors.NewDatabaseError("circuit breaker is open", "query", "families", err)
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
		return nil, errors.NewDatabaseError("failed to get family from MongoDB after retries", "query", "families", retryErr)
	}

	r.logger.Debug(ctx, "Successfully retrieved family from MongoDB", zap.String("family_id", id))
	return family, nil
}

// Save persists a family
func (r *MongoFamilyRepository) Save(ctx context.Context, fam *entity.Family) error {
	if fam == nil {
		r.logger.Warn(ctx, "Family cannot be nil for Save")
		return errors.NewValidationError("family cannot be nil", "family", nil)
	}

	r.logger.Debug(ctx, "Saving family to MongoDB", zap.String("family_id", fam.ID()))

	if err := fam.Validate(); err != nil {
		r.logger.Error(ctx, "Family validation failed", zap.Error(err), zap.String("family_id", fam.ID()))
		return err
	}

	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Convert domain entity to document
	doc := r.entityToDocument(fam)
	var retryErr error

	// Define the operation to retry
	operation := func(ctx context.Context) error {
		// Use ReplaceOne with upsert to handle both insert and update
		// Query by family_id instead of _id
		_, err := r.Collection.ReplaceOne(
			ctx,
			bson.M{"family_id": doc.FamilyID},
			doc,
			options.Replace().SetUpsert(true),
		)

		if err != nil {
			r.logger.Error(ctx, "Failed to save family to MongoDB", zap.Error(err), zap.String("family_id", fam.ID()))
			return errors.NewDatabaseError("failed to save family to MongoDB", "save", "families", err)
		}
		return nil
	}

	// Define which errors are retryable
	isRetryable := func(err error) bool {
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
		_, err := circuit.Execute(ctx, r.circuitBreaker, "Save", circuitOpWrapper)
		return err
	}

	// Execute with rate limiter
	// We need to wrap the rateOperation to match the generic function signature
	rateOpWrapper := func(ctx context.Context) (bool, error) {
		err := rateOperation(ctx)
		return err == nil, err
	}
	_, err := rate.Execute(ctxWithTimeout, r.rateLimiter, "Save", rateOpWrapper)

	// Check for errors from rate limiter or circuit breaker
	if err != nil && retryErr == nil {
		// Check if it's a rate limiter error
		if strings.Contains(err.Error(), "rate limit exceeded") {
			return errors.NewDatabaseError("rate limit exceeded", "save", "families", err)
		}
		// Otherwise, assume it's a circuit breaker error
		return errors.NewDatabaseError("circuit breaker is open", "save", "families", err)
	}

	// Handle retry errors
	if retryErr != nil {
		// If it's already a typed error, return it directly
		if _, ok := retryErr.(*errors.ValidationError); ok {
			return retryErr
		}
		if _, ok := retryErr.(*errors.DatabaseError); ok {
			return retryErr
		}

		// Otherwise, wrap it in a database error
		return errors.NewDatabaseError("failed to save family to MongoDB after retries", "save", "families", retryErr)
	}

	r.logger.Debug(ctx, "Successfully saved family to MongoDB", zap.String("family_id", fam.ID()))
	return nil
}

// FindByParentID finds families that contain a specific parent
func (r *MongoFamilyRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error) {
	r.logger.Debug(ctx, "Finding families by parent ID in MongoDB", zap.String("parent_id", parentID))

	if parentID == "" {
		r.logger.Warn(ctx, "Parent ID is required for FindByParentID")
		return nil, errors.NewValidationError("parent ID is required", "parentID", nil)
	}

	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var families []*entity.Family
	var retryErr error

	// Define the operation to retry
	operation := func(ctx context.Context) error {
		// No change needed here, as parents.id is still the same field
		cursor, err := r.Collection.Find(ctx, bson.M{"parents.id": parentID})
		if err != nil {
			r.logger.Error(ctx, "Failed to find families by parent ID in MongoDB", zap.Error(err), zap.String("parent_id", parentID))
			return errors.NewDatabaseError("failed to find families by parent ID", "query", "families", err)
		}
		defer cursor.Close(ctx)

		var docs []FamilyDocument
		if err := cursor.All(ctx, &docs); err != nil {
			r.logger.Error(ctx, "Failed to decode families from MongoDB", zap.Error(err))
			return errors.NewDatabaseError("failed to decode families", "query", "families", err)
		}

		// Convert documents to domain entities
		families = make([]*entity.Family, 0, len(docs))
		for _, doc := range docs {
			fam, err := r.documentToEntity(doc)
			if err != nil {
				r.logger.Error(ctx, "Failed to convert document to entity", zap.Error(err), zap.String("family_id", doc.FamilyID))
				return err
			}
			families = append(families, fam)
		}

		return nil
	}

	// Define which errors are retryable
	isRetryable := func(err error) bool {
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
		_, err := circuit.Execute(ctx, r.circuitBreaker, "FindByParentID", circuitOpWrapper)
		return err
	}

	// Execute with rate limiter
	// We need to wrap the rateOperation to match the generic function signature
	rateOpWrapper := func(ctx context.Context) (bool, error) {
		err := rateOperation(ctx)
		return err == nil, err
	}
	_, err := rate.Execute(ctxWithTimeout, r.rateLimiter, "FindByParentID", rateOpWrapper)

	// Check for errors from rate limiter or circuit breaker
	if err != nil && retryErr == nil {
		// Check if it's a rate limiter error
		if strings.Contains(err.Error(), "rate limit exceeded") {
			return nil, errors.NewDatabaseError("rate limit exceeded", "query", "families", err)
		}
		// Otherwise, assume it's a circuit breaker error
		return nil, errors.NewDatabaseError("circuit breaker is open", "query", "families", err)
	}

	// Handle retry errors
	if retryErr != nil {
		// If it's already a typed error, return it directly
		if _, ok := retryErr.(*errors.ValidationError); ok {
			return nil, retryErr
		}
		if _, ok := retryErr.(*errors.DatabaseError); ok {
			return nil, retryErr
		}

		// Otherwise, wrap it in a database error
		return nil, errors.NewDatabaseError("failed to find families by parent ID after retries", "query", "families", retryErr)
	}

	r.logger.Debug(ctx, "Successfully found families by parent ID in MongoDB",
		zap.String("parent_id", parentID),
		zap.Int("count", len(families)))
	return families, nil
}

// FindByChildID finds the family that contains a specific child
func (r *MongoFamilyRepository) FindByChildID(ctx context.Context, childID string) (*entity.Family, error) {
	r.logger.Debug(ctx, "Finding family by child ID in MongoDB", zap.String("child_id", childID))

	if childID == "" {
		r.logger.Warn(ctx, "Child ID is required for FindByChildID")
		return nil, errors.NewValidationError("child ID is required", "childID", nil)
	}

	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var doc FamilyDocument
	var family *entity.Family
	var retryErr error

	// Define the operation to retry
	operation := func(ctx context.Context) error {
		// No change needed here, as children.id is still the same field
		err := r.Collection.FindOne(ctx, bson.M{"children.id": childID}).Decode(&doc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				r.logger.Info(ctx, "Family with child not found in MongoDB", zap.String("child_id", childID))
				return errors.NewNotFoundError("Family with Child", childID, nil)
			}
			r.logger.Error(ctx, "Failed to find family by child ID in MongoDB", zap.Error(err), zap.String("child_id", childID))
			return errors.NewDatabaseError("failed to find family by child ID", "query", "families", err)
		}

		// Convert document to domain entity
		var convErr error
		family, convErr = r.documentToEntity(doc)
		return convErr
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
		_, err := circuit.Execute(ctx, r.circuitBreaker, "FindByChildID", circuitOpWrapper)
		return err
	}

	// Execute with rate limiter
	// We need to wrap the rateOperation to match the generic function signature
	rateOpWrapper := func(ctx context.Context) (bool, error) {
		err := rateOperation(ctx)
		return err == nil, err
	}
	_, err := rate.Execute(ctxWithTimeout, r.rateLimiter, "FindByChildID", rateOpWrapper)

	// Check for errors from rate limiter or circuit breaker
	if err != nil && retryErr == nil {
		// Check if it's a rate limiter error
		if strings.Contains(err.Error(), "rate limit exceeded") {
			return nil, errors.NewDatabaseError("rate limit exceeded", "query", "families", err)
		}
		// Otherwise, assume it's a circuit breaker error
		return nil, errors.NewDatabaseError("circuit breaker is open", "query", "families", err)
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
		return nil, errors.NewDatabaseError("failed to find family by child ID after retries", "query", "families", retryErr)
	}

	r.logger.Debug(ctx, "Successfully found family by child ID in MongoDB",
		zap.String("child_id", childID),
		zap.String("family_id", doc.FamilyID))

	return family, nil
}

// GetAll retrieves all families
func (r *MongoFamilyRepository) GetAll(ctx context.Context) ([]*entity.Family, error) {
	r.logger.Debug(ctx, "Getting all families from MongoDB")

	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var families []*entity.Family
	var retryErr error

	// Define the operation to retry
	operation := func(ctx context.Context) error {
		// Find all documents in the collection
		cursor, err := r.Collection.Find(ctx, bson.M{})
		if err != nil {
			r.logger.Error(ctx, "Failed to get all families from MongoDB", zap.Error(err))
			return errors.NewDatabaseError("failed to get all families", "query", "families", err)
		}
		defer cursor.Close(ctx)

		var docs []FamilyDocument
		if err := cursor.All(ctx, &docs); err != nil {
			r.logger.Error(ctx, "Failed to decode families from MongoDB", zap.Error(err))
			return errors.NewDatabaseError("failed to decode families", "query", "families", err)
		}

		// Convert documents to domain entities
		families = make([]*entity.Family, 0, len(docs))
		for _, doc := range docs {
			fam, err := r.documentToEntity(doc)
			if err != nil {
				r.logger.Error(ctx, "Failed to convert document to entity", zap.Error(err), zap.String("family_id", doc.FamilyID))
				return err
			}
			families = append(families, fam)
		}

		return nil
	}

	// Define which errors are retryable
	isRetryable := func(err error) bool {
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
		_, err := circuit.Execute(ctx, r.circuitBreaker, "GetAll", circuitOpWrapper)
		return err
	}

	// Execute with rate limiter
	// We need to wrap the rateOperation to match the generic function signature
	rateOpWrapper := func(ctx context.Context) (bool, error) {
		err := rateOperation(ctx)
		return err == nil, err
	}
	_, err := rate.Execute(ctxWithTimeout, r.rateLimiter, "GetAll", rateOpWrapper)

	// Check for errors from rate limiter or circuit breaker
	if err != nil && retryErr == nil {
		// Check if it's a rate limiter error
		if strings.Contains(err.Error(), "rate limit exceeded") {
			return nil, errors.NewDatabaseError("rate limit exceeded", "query", "families", err)
		}
		// Otherwise, assume it's a circuit breaker error
		return nil, errors.NewDatabaseError("circuit breaker is open", "query", "families", err)
	}

	// Handle retry errors
	if retryErr != nil {
		// If it's already a typed error, return it directly
		if _, ok := retryErr.(*errors.ValidationError); ok {
			return nil, retryErr
		}
		if _, ok := retryErr.(*errors.DatabaseError); ok {
			return nil, retryErr
		}

		// Otherwise, wrap it in a database error
		return nil, errors.NewDatabaseError("failed to get all families after retries", "query", "families", retryErr)
	}

	r.logger.Debug(ctx, "Successfully retrieved all families from MongoDB", zap.Int("count", len(families)))
	return families, nil
}

// documentToEntity converts a FamilyDocument to a Family entity
func (r *MongoFamilyRepository) documentToEntity(doc FamilyDocument) (*entity.Family, error) {
	// Convert parents
	parents := make([]*entity.Parent, 0, len(doc.Parents))
	for _, p := range doc.Parents {
		// Parse dates
		birthDate, err := time.Parse(time.RFC3339, p.BirthDate)
		if err != nil {
			return nil, errors.NewDatabaseError("invalid parent birth date format", "parse", "families", err)
		}

		var deathDate *time.Time
		if p.DeathDate != nil {
			parsedDeathDate, err := time.Parse(time.RFC3339, *p.DeathDate)
			if err != nil {
				return nil, errors.NewDatabaseError("invalid parent death date format", "parse", "families", err)
			}
			deathDate = &parsedDeathDate
		}

		// Create parent entity
		parentEntity, err := entity.NewParent(p.ID, p.FirstName, p.LastName, birthDate, deathDate)
		if err != nil {
			return nil, err
		}
		parents = append(parents, parentEntity)
	}

	// Convert children
	children := make([]*entity.Child, 0, len(doc.Children))
	for _, c := range doc.Children {
		// Parse dates
		birthDate, err := time.Parse(time.RFC3339, c.BirthDate)
		if err != nil {
			return nil, errors.NewDatabaseError("invalid child birth date format", "parse", "families", err)
		}

		var deathDate *time.Time
		if c.DeathDate != nil {
			parsedDeathDate, err := time.Parse(time.RFC3339, *c.DeathDate)
			if err != nil {
				return nil, errors.NewDatabaseError("invalid child death date format", "parse", "families", err)
			}
			deathDate = &parsedDeathDate
		}

		// Create child entity
		childEntity, err := entity.NewChild(c.ID, c.FirstName, c.LastName, birthDate, deathDate)
		if err != nil {
			return nil, err
		}
		children = append(children, childEntity)
	}

	// Create family entity
	// Use FamilyID field which contains the string ID
	return entity.NewFamily(doc.FamilyID, entity.Status(doc.Status), parents, children)
}

// entityToDocument converts a Family entity to a FamilyDocument
func (r *MongoFamilyRepository) entityToDocument(fam *entity.Family) FamilyDocument {
	// Convert parents
	parents := make([]ParentDocument, 0, len(fam.Parents()))
	for _, p := range fam.Parents() {
		var deathDateStr *string
		if p.DeathDate() != nil {
			str := p.DeathDate().Format(time.RFC3339)
			deathDateStr = &str
		}

		parents = append(parents, ParentDocument{
			ID:        p.ID(),
			FirstName: p.FirstName(),
			LastName:  p.LastName(),
			BirthDate: p.BirthDate().Format(time.RFC3339),
			DeathDate: deathDateStr,
		})
	}

	// Convert children
	children := make([]ChildDocument, 0, len(fam.Children()))
	for _, c := range fam.Children() {
		var deathDateStr *string
		if c.DeathDate() != nil {
			str := c.DeathDate().Format(time.RFC3339)
			deathDateStr = &str
		}

		children = append(children, ChildDocument{
			ID:        c.ID(),
			FirstName: c.FirstName(),
			LastName:  c.LastName(),
			BirthDate: c.BirthDate().Format(time.RFC3339),
			DeathDate: deathDateStr,
		})
	}

	// Create a new ObjectID for MongoDB
	objectID := primitive.NewObjectID()

	// Create document
	return FamilyDocument{
		ID:       objectID,
		FamilyID: fam.ID(),
		Status:   string(fam.Status()),
		Parents:  parents,
		Children: children,
	}
}
