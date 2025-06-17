package mongo

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/core/domain/ports"
	"github.com/abitofhelp/servicelib/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FamilyDocument represents how a family is stored in MongoDB
type FamilyDocument struct {
	ID       string          `bson:"_id"`
	Status   string          `bson:"status"`
	Parents  []ParentDocument  `bson:"parents"`
	Children []ChildDocument   `bson:"children"`
}

// ParentDocument represents how a parent is stored in MongoDB
type ParentDocument struct {
	ID        string    `bson:"id"`
	FirstName string    `bson:"firstName"`
	LastName  string    `bson:"lastName"`
	BirthDate string    `bson:"birthDate"`
	DeathDate *string   `bson:"deathDate,omitempty"`
}

// ChildDocument represents how a child is stored in MongoDB
type ChildDocument struct {
	ID        string    `bson:"id"`
	FirstName string    `bson:"firstName"`
	LastName  string    `bson:"lastName"`
	BirthDate string    `bson:"birthDate"`
	DeathDate *string   `bson:"deathDate,omitempty"`
}

// MongoFamilyRepository implements the ports.FamilyRepository interface for MongoDB
type MongoFamilyRepository struct {
	Collection *mongo.Collection
}

// Ensure MongoFamilyRepository implements ports.FamilyRepository
var _ ports.FamilyRepository = (*MongoFamilyRepository)(nil)

// NewMongoFamilyRepository creates a new MongoFamilyRepository
func NewMongoFamilyRepository(collection *mongo.Collection) *MongoFamilyRepository {
	if collection == nil {
		panic("collection cannot be nil")
	}
	return &MongoFamilyRepository{Collection: collection}
}

// GetByID retrieves a family by its ID
func (r *MongoFamilyRepository) GetByID(ctx context.Context, id string) (*entity.Family, error) {
	if id == "" {
		return nil, errors.NewValidationError("id is required")
	}

	var doc FamilyDocument
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("Family", id)
		}
		return nil, errors.NewRepositoryError(err, "failed to get family from MongoDB", "MONGO_ERROR")
	}

	// Convert document to domain entity
	return r.documentToEntity(doc)
}

// Save persists a family
func (r *MongoFamilyRepository) Save(ctx context.Context, fam *entity.Family) error {
	if fam == nil {
		return errors.NewValidationError("family cannot be nil")
	}

	if err := fam.Validate(); err != nil {
		return err
	}

	// Convert domain entity to document
	doc := r.entityToDocument(fam)

	// Use ReplaceOne with upsert to handle both insert and update
	_, err := r.Collection.ReplaceOne(
		ctx,
		bson.M{"_id": doc.ID},
		doc,
		options.Replace().SetUpsert(true),
	)

	if err != nil {
		return errors.NewRepositoryError(err, "failed to save family to MongoDB", "MONGO_ERROR")
	}

	return nil
}

// FindByParentID finds families that contain a specific parent
func (r *MongoFamilyRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error) {
	if parentID == "" {
		return nil, errors.NewValidationError("parent ID is required")
	}

	cursor, err := r.Collection.Find(ctx, bson.M{"parents.id": parentID})
	if err != nil {
		return nil, errors.NewRepositoryError(err, "failed to find families by parent ID", "MONGO_ERROR")
	}
	defer cursor.Close(ctx)

	var docs []FamilyDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.NewRepositoryError(err, "failed to decode families", "MONGO_ERROR")
	}

	// Convert documents to domain entities
	families := make([]*entity.Family, 0, len(docs))
	for _, doc := range docs {
		fam, err := r.documentToEntity(doc)
		if err != nil {
			return nil, err
		}
		families = append(families, fam)
	}

	return families, nil
}

// FindByChildID finds the family that contains a specific child
func (r *MongoFamilyRepository) FindByChildID(ctx context.Context, childID string) (*entity.Family, error) {
	if childID == "" {
		return nil, errors.NewValidationError("child ID is required")
	}

	var doc FamilyDocument
	err := r.Collection.FindOne(ctx, bson.M{"children.id": childID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("Family with Child", childID)
		}
		return nil, errors.NewRepositoryError(err, "failed to find family by child ID", "MONGO_ERROR")
	}

	// Convert document to domain entity
	return r.documentToEntity(doc)
}

// GetAll retrieves all families
func (r *MongoFamilyRepository) GetAll(ctx context.Context) ([]*entity.Family, error) {
	// Find all documents in the collection
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.NewRepositoryError(err, "failed to get all families", "MONGO_ERROR")
	}
	defer cursor.Close(ctx)

	var docs []FamilyDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.NewRepositoryError(err, "failed to decode families", "MONGO_ERROR")
	}

	// Convert documents to domain entities
	families := make([]*entity.Family, 0, len(docs))
	for _, doc := range docs {
		fam, err := r.documentToEntity(doc)
		if err != nil {
			return nil, err
		}
		families = append(families, fam)
	}

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
			return nil, errors.NewRepositoryError(err, "invalid parent birth date format", "DATA_FORMAT_ERROR")
		}

		var deathDate *time.Time
		if p.DeathDate != nil {
			parsedDeathDate, err := time.Parse(time.RFC3339, *p.DeathDate)
			if err != nil {
				return nil, errors.NewRepositoryError(err, "invalid parent death date format", "DATA_FORMAT_ERROR")
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
			return nil, errors.NewRepositoryError(err, "invalid child birth date format", "DATA_FORMAT_ERROR")
		}

		var deathDate *time.Time
		if c.DeathDate != nil {
			parsedDeathDate, err := time.Parse(time.RFC3339, *c.DeathDate)
			if err != nil {
				return nil, errors.NewRepositoryError(err, "invalid child death date format", "DATA_FORMAT_ERROR")
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
	return entity.NewFamily(doc.ID, entity.Status(doc.Status), parents, children)
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

	// Create document
	return FamilyDocument{
		ID:       fam.ID(),
		Status:   string(fam.Status()),
		Parents:  parents,
		Children: children,
	}
}
