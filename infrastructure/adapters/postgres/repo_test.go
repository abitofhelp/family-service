// Copyright (c) 2025 A Bit of Help, Inc.

package postgres

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// generateTestUUID generates a UUID for testing
func generateTestUUID() string {
	return uuid.New().String()
}

// setupTest sets up common test dependencies
func setupTest(t *testing.T) *PostgresFamilyRepository {
	logger := zaptest.NewLogger(t)
	return &PostgresFamilyRepository{
		logger: logging.NewContextLogger(logger),
	}
}

// TestNewRepositoryError tests the NewRepositoryError function
func TestNewRepositoryError(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		err         error
		message     string
		code        string
		expectedMsg string
	}{
		{
			name:        "PostgreSQL error",
			err:         pgx.ErrNoRows,
			message:     "database error",
			code:        "POSTGRES_ERROR",
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

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test
			err := NewRepositoryError(tc.err, tc.message, tc.code)

			// Verify
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedMsg)
		})
	}
}

// TestJSONParsingAndFormatting tests the JSON parsing and formatting functions
func TestJSONParsingAndFormatting(t *testing.T) {
	// Test JSON marshaling and unmarshaling
	type jsonParent struct {
		ID        string  `json:"id"`
		FirstName string  `json:"firstName"`
		LastName  string  `json:"lastName"`
		BirthDate string  `json:"birthDate"`
		DeathDate *string `json:"deathDate,omitempty"`
	}

	// Create test data
	birthDateStr := time.Now().AddDate(-30, 0, 0).Format(time.RFC3339)
	deathDateStr := time.Now().AddDate(-1, 0, 0).Format(time.RFC3339)

	parent := jsonParent{
		ID:        generateTestUUID(),
		FirstName: "John",
		LastName:  "Doe",
		BirthDate: birthDateStr,
		DeathDate: &deathDateStr,
	}

	// Test marshaling
	data, err := json.Marshal(parent)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test unmarshaling
	var parsedParent jsonParent
	err = json.Unmarshal(data, &parsedParent)
	require.NoError(t, err)
	assert.Equal(t, parent.ID, parsedParent.ID)
	assert.Equal(t, parent.FirstName, parsedParent.FirstName)
	assert.Equal(t, parent.LastName, parsedParent.LastName)
	assert.Equal(t, parent.BirthDate, parsedParent.BirthDate)
	assert.Equal(t, *parent.DeathDate, *parsedParent.DeathDate)
}

// TestDateParsing tests the date parsing functions
func TestDateParsing(t *testing.T) {
	// Test valid date parsing
	validDateStr := time.Now().Format(time.RFC3339)
	parsedDate, err := time.Parse(time.RFC3339, validDateStr)
	require.NoError(t, err)
	assert.Equal(t, validDateStr, parsedDate.Format(time.RFC3339))

	// Test invalid date parsing
	invalidDateStr := "invalid-date"
	_, err = time.Parse(time.RFC3339, invalidDateStr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing time")
}

// TestFieldCaseHandling tests the handling of different field name cases
func TestFieldCaseHandling(t *testing.T) {
	// Test handling of uppercase and lowercase field names
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

	// Test with uppercase fields
	upperJSON := `{
		"ID": "123",
		"FirstName": "John",
		"LastName": "Doe",
		"BirthDate": "2000-01-01T00:00:00Z",
		"DeathDate": "2020-01-01T00:00:00Z"
	}`

	var upperParent jsonParent
	err := json.Unmarshal([]byte(upperJSON), &upperParent)
	require.NoError(t, err)
	assert.Equal(t, "123", upperParent.ID)
	assert.Equal(t, "John", upperParent.FirstName)
	assert.Equal(t, "Doe", upperParent.LastName)
	assert.Equal(t, "2000-01-01T00:00:00Z", upperParent.BirthDate)
	assert.Equal(t, "2020-01-01T00:00:00Z", *upperParent.DeathDate)

	// Test with lowercase fields
	lowerJSON := `{
		"id": "456",
		"firstName": "Jane",
		"lastName": "Doe",
		"birthDate": "2001-01-01T00:00:00Z",
		"deathDate": "2021-01-01T00:00:00Z"
	}`

	var lowerParent jsonParent
	err = json.Unmarshal([]byte(lowerJSON), &lowerParent)
	require.NoError(t, err)
	assert.Equal(t, "456", lowerParent.Id)
	assert.Equal(t, "Jane", lowerParent.FirstN)
	assert.Equal(t, "Doe", lowerParent.LastN)
	assert.Equal(t, "2001-01-01T00:00:00Z", lowerParent.BirthD)
	assert.Equal(t, "2021-01-01T00:00:00Z", *lowerParent.DeathD)

	// Test field selection logic
	id := upperParent.ID
	if id == "" {
		id = upperParent.Id
	}
	assert.Equal(t, "123", id)

	firstName := upperParent.FirstName
	if firstName == "" {
		firstName = upperParent.FirstN
	}
	assert.Equal(t, "John", firstName)
}
