// Copyright (c) 2025 A Bit of Help, Inc.

package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

// generateTestUUID generates a UUID for testing
func generateTestUUID() string {
	return uuid.New().String()
}

// setupTest sets up common test dependencies
func setupTest(t *testing.T) (*SQLiteFamilyRepository, *sql.DB, *gomock.Controller) {
	logger := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)

	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	repo := NewSQLiteFamilyRepository(db, logging.NewContextLogger(logger))

	// Initialize the database schema
	err = repo.ensureTableExists(context.Background())
	require.NoError(t, err)

	return repo, db, ctrl
}

// TestNewRepositoryError tests the NewRepositoryError function
func TestNewRepositoryError(t *testing.T) {
	testCases := []struct {
		name        string
		err         error
		message     string
		code        string
		expectedMsg string
	}{
		{
			name:        "SQLite error",
			err:        &sqlite3.Error{Code: sqlite3.ErrError},
			message:     "database error",
			code:        "SQLITE_ERROR",
			expectedMsg: "database error",
		},
		{
			name:        "JSON error",
			err:         &json.SyntaxError{},
			message:     "json error",
			code:        "JSON_ERROR",
			expectedMsg: "json error",
		},
		{
			name:        "Data format error",
			err:         &time.ParseError{},
			message:     "data format error",
			code:        "DATA_FORMAT_ERROR",
			expectedMsg: "data format error",
		},
		{
			name:        "Conversion error",
			err:         &json.UnmarshalTypeError{},
			message:     "conversion error",
			code:        "CONVERSION_ERROR",
			expectedMsg: "conversion error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewRepositoryError(tc.err, tc.message, tc.code)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedMsg)
		})
	}
}

// TestSQLiteFamilyRepository_GetByID tests the GetByID method
func TestSQLiteFamilyRepository_GetByID(t *testing.T) {
	repo, db, ctrl := setupTest(t)
	defer ctrl.Finish()
	defer db.Close()

	// Create test data
	parentID := generateTestUUID()
	parentBirthDate := time.Now().AddDate(-30, 0, 0)
	parent, err := entity.NewParent(parentID, "John", "Doe", parentBirthDate, nil)
	require.NoError(t, err)

	childID := generateTestUUID()
	childBirthDate := time.Now().AddDate(-5, 0, 0)
	child, err := entity.NewChild(childID, "Jane", "Doe", childBirthDate, nil)
	require.NoError(t, err)

	familyID := generateTestUUID()
	testFamily, err := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent}, []*entity.Child{child})
	require.NoError(t, err)

	// Save test family
	err = repo.Save(context.Background(), testFamily)
	require.NoError(t, err)

	t.Run("successful retrieval", func(t *testing.T) {
		// Test
		retrieved, err := repo.GetByID(context.Background(), testFamily.ID())

		// Verify
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, testFamily.ID(), retrieved.ID())
		assert.Equal(t, testFamily.Status(), retrieved.Status())
		assert.Len(t, retrieved.Parents(), 1)
		assert.Len(t, retrieved.Children(), 1)
		assert.Equal(t, testFamily.Parents()[0].FirstName(), retrieved.Parents()[0].FirstName())
		assert.Equal(t, testFamily.Children()[0].FirstName(), retrieved.Children()[0].FirstName())
	})

	t.Run("not found", func(t *testing.T) {
		// Test
		retrieved, err := repo.GetByID(context.Background(), generateTestUUID())

		// Verify
		require.Error(t, err)
		assert.Nil(t, retrieved)
	})
}

// TestSQLiteFamilyRepository_Save tests the Save method
func TestSQLiteFamilyRepository_Save(t *testing.T) {
	repo, db, ctrl := setupTest(t)
	defer ctrl.Finish()
	defer db.Close()

	t.Run("successful save", func(t *testing.T) {
		// Create test data
		parentID := generateTestUUID()
		parentBirthDate := time.Now().AddDate(-30, 0, 0)
		parent, err := entity.NewParent(parentID, "John", "Doe", parentBirthDate, nil)
		require.NoError(t, err)

		familyID := generateTestUUID()
		family, err := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent}, []*entity.Child{})
		require.NoError(t, err)

		// Test
		err = repo.Save(context.Background(), family)

		// Verify
		require.NoError(t, err)

		// Verify the save by retrieving
		retrieved, err := repo.GetByID(context.Background(), family.ID())
		require.NoError(t, err)
		assert.Equal(t, family.ID(), retrieved.ID())
		assert.Equal(t, family.Status(), retrieved.Status())
		assert.Len(t, retrieved.Parents(), 1)
		assert.Equal(t, family.Parents()[0].FirstName(), retrieved.Parents()[0].FirstName())
	})

	t.Run("invalid family", func(t *testing.T) {
		// Test with nil family
		err := repo.Save(context.Background(), nil)

		// Verify
		require.Error(t, err)
	})
}

// TestSQLiteFamilyRepository_FindByParentID tests the FindByParentID method
func TestSQLiteFamilyRepository_FindByParentID(t *testing.T) {
	repo, db, ctrl := setupTest(t)
	defer ctrl.Finish()
	defer db.Close()

	parentID := generateTestUUID()
	parentBirthDate := time.Now().AddDate(-30, 0, 0)
	parent, err := entity.NewParent(parentID, "John", "Doe", parentBirthDate, nil)
	require.NoError(t, err)

	familyID := generateTestUUID()
	family, err := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent}, []*entity.Child{})
	require.NoError(t, err)

	// Save test family
	err = repo.Save(context.Background(), family)
	require.NoError(t, err)

	t.Run("successful find", func(t *testing.T) {
		// Test
		families, err := repo.FindByParentID(context.Background(), parentID)

		// Verify
		require.NoError(t, err)
		assert.Len(t, families, 1)
		assert.Equal(t, family.ID(), families[0].ID())
	})

	t.Run("parent not found", func(t *testing.T) {
		// Test
		families, err := repo.FindByParentID(context.Background(), generateTestUUID())

		// Verify
		require.NoError(t, err)
		assert.Empty(t, families)
	})
}

// TestSQLiteFamilyRepository_FindByChildID tests the FindByChildID method
func TestSQLiteFamilyRepository_FindByChildID(t *testing.T) {
	repo, db, ctrl := setupTest(t)
	defer ctrl.Finish()
	defer db.Close()

	childID := generateTestUUID()
	childBirthDate := time.Now().AddDate(-5, 0, 0)
	child, err := entity.NewChild(childID, "Jane", "Doe", childBirthDate, nil)
	require.NoError(t, err)

	// Create a parent for the family (required by validation)
	parentID := generateTestUUID()
	parentBirthDate := time.Now().AddDate(-30, 0, 0)
	parent, err := entity.NewParent(parentID, "John", "Doe", parentBirthDate, nil)
	require.NoError(t, err)

	familyID := generateTestUUID()
	family, err := entity.NewFamily(familyID, entity.Single, []*entity.Parent{parent}, []*entity.Child{child})
	require.NoError(t, err)

	// Save test family
	err = repo.Save(context.Background(), family)
	require.NoError(t, err)

	t.Run("successful find", func(t *testing.T) {
		// Test
		found, err := repo.FindByChildID(context.Background(), childID)

		// Verify
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, family.ID(), found.ID())
	})

	t.Run("child not found", func(t *testing.T) {
		// Test
		found, err := repo.FindByChildID(context.Background(), generateTestUUID())

		// Verify
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

// TestSQLiteFamilyRepository_GetAll tests the GetAll method
func TestSQLiteFamilyRepository_GetAll(t *testing.T) {
	repo, db, ctrl := setupTest(t)
	defer ctrl.Finish()
	defer db.Close()

	// Create multiple test families
	families := make([]*entity.Family, 0, 2)

	// Create first family
	parent1ID := generateTestUUID()
	parent1BirthDate := time.Now().AddDate(-30, 0, 0)
	parent1, err := entity.NewParent(parent1ID, "John", "Doe", parent1BirthDate, nil)
	require.NoError(t, err)

	family1ID := generateTestUUID()
	family1, err := entity.NewFamily(family1ID, entity.Single, []*entity.Parent{parent1}, []*entity.Child{})
	require.NoError(t, err)
	families = append(families, family1)

	// Create second family
	parent2ID := generateTestUUID()
	parent2BirthDate := time.Now().AddDate(-28, 0, 0)
	parent2, err := entity.NewParent(parent2ID, "Jane", "Smith", parent2BirthDate, nil)
	require.NoError(t, err)

	family2ID := generateTestUUID()
	family2, err := entity.NewFamily(family2ID, entity.Single, []*entity.Parent{parent2}, []*entity.Child{})
	require.NoError(t, err)
	families = append(families, family2)

	// Save test families
	for _, f := range families {
		err = repo.Save(context.Background(), f)
		require.NoError(t, err)
	}

	t.Run("successful retrieval", func(t *testing.T) {
		// Test
		retrieved, err := repo.GetAll(context.Background())

		// Verify
		require.NoError(t, err)
		assert.Len(t, retrieved, len(families))
	})
}
