// Copyright (c) 2025 A Bit of Help, Inc.

package mongo

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

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
	Collection *mongo.Collection
	logger     *logging.ContextLogger
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
	return &MongoFamilyRepository{
		Collection: collection,
		logger:     logger,
	}
}

// GetByID retrieves a family by its ID
func (r *MongoFamilyRepository) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	r.logger.Debug(ctx, "Getting family by ID from MongoDB", zap.String("family_id", id))

	if id == "" {
		r.logger.Warn(ctx, "Family ID is required for GetByID")
		return nil, errors.NewValidationError("id is required", "id", nil)
	}

	var doc FamilyDocument
	err := r.Collection.FindOne(ctx, bson.M{"family_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Info(ctx, "Family not found in MongoDB", zap.String("family_id", id))
			return nil, errors.NewNotFoundError("Family", id, nil)
		}
		r.logger.Error(ctx, "Failed to get family from MongoDB", zap.Error(err), zap.String("family_id", id))
		return nil, errors.NewDatabaseError("failed to get family from MongoDB", "query", "families", err)
	}

	r.logger.Debug(ctx, "Successfully retrieved family from MongoDB", zap.String("family_id", id))

	// Convert document to domain entity
	return r.documentToEntity(doc)
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

	// Convert domain entity to document
	doc := r.entityToDocument(fam)

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

	// No change needed here, as parents.id is still the same field
	cursor, err := r.Collection.Find(ctx, bson.M{"parents.id": parentID})
	if err != nil {
		r.logger.Error(ctx, "Failed to find families by parent ID in MongoDB", zap.Error(err), zap.String("parent_id", parentID))
		return nil, errors.NewDatabaseError("failed to find families by parent ID", "query", "families", err)
	}
	defer cursor.Close(ctx)

	var docs []FamilyDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error(ctx, "Failed to decode families from MongoDB", zap.Error(err))
		return nil, errors.NewDatabaseError("failed to decode families", "query", "families", err)
	}

	// Convert documents to domain entities
	families := make([]*entity.Family, 0, len(docs))
	for _, doc := range docs {
		fam, err := r.documentToEntity(doc)
		if err != nil {
			r.logger.Error(ctx, "Failed to convert document to entity", zap.Error(err), zap.String("family_id", doc.FamilyID))
			return nil, err
		}
		families = append(families, fam)
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

	// No change needed here, as children.id is still the same field
	var doc FamilyDocument
	err := r.Collection.FindOne(ctx, bson.M{"children.id": childID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Info(ctx, "Family with child not found in MongoDB", zap.String("child_id", childID))
			return nil, errors.NewNotFoundError("Family with Child", childID, nil)
		}
		r.logger.Error(ctx, "Failed to find family by child ID in MongoDB", zap.Error(err), zap.String("child_id", childID))
		return nil, errors.NewDatabaseError("failed to find family by child ID", "query", "families", err)
	}

	r.logger.Debug(ctx, "Successfully found family by child ID in MongoDB", 
		zap.String("child_id", childID), 
		zap.String("family_id", doc.FamilyID))

	// Convert document to domain entity
	return r.documentToEntity(doc)
}

// GetAll retrieves all families
func (r *MongoFamilyRepository) GetAll(ctx context.Context) ([]*entity.Family, error) {
	r.logger.Debug(ctx, "Getting all families from MongoDB")

	// Find all documents in the collection
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		r.logger.Error(ctx, "Failed to get all families from MongoDB", zap.Error(err))
		return nil, errors.NewDatabaseError("failed to get all families", "query", "families", err)
	}
	defer cursor.Close(ctx)

	var docs []FamilyDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error(ctx, "Failed to decode families from MongoDB", zap.Error(err))
		return nil, errors.NewDatabaseError("failed to decode families", "query", "families", err)
	}

	// Convert documents to domain entities
	families := make([]*entity.Family, 0, len(docs))
	for _, doc := range docs {
		fam, err := r.documentToEntity(doc)
		if err != nil {
			r.logger.Error(ctx, "Failed to convert document to entity", zap.Error(err), zap.String("family_id", doc.FamilyID))
			return nil, err
		}
		families = append(families, fam)
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
