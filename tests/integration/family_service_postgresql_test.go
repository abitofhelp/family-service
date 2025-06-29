// Copyright (c) 2025 A Bit of Help, Inc.

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/abitofhelp/family-service/cmd/server/graphql/di"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/generated"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/resolver"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// setupPostgreSQLTestConfig creates a test configuration with PostgreSQL database
func setupPostgreSQLTestConfig() *config.Config {
	cfg := &config.Config{}

	// Set up database config
	cfg.Database.Type = "postgres"

	// Use environment variable for PostgreSQL DSN if available, otherwise use default with authentication
	postgresDSN := os.Getenv("POSTGRES_TEST_DSN")
	if postgresDSN == "" {
		// Get PostgreSQL credentials from environment variables or use defaults from secrets
		username := os.Getenv("POSTGRESQL_USERNAME")
		if username == "" {
			username = "postgres"
		}

		password := os.Getenv("POSTGRESQL_PASSWORD")
		if password == "" {
			password = "NVsHFXcxqUsMoEgiUnE7jvzXxhp3gn9nsgkXCsetAHLhcpyLRmWhKixUpfr3J7tE"
		}

		// Create PostgreSQL DSN with authentication
		postgresDSN = fmt.Sprintf("postgres://%s:%s@localhost:5432/family_service?sslmode=disable", username, password)
	}
	cfg.Database.Postgres.DSN = postgresDSN

	// Set up auth config
	cfg.Auth.JWT.SecretKey = "test-secret-key-that-is-at-least-32-characters-long"
	cfg.Auth.JWT.Issuer = "test-issuer"
	cfg.Auth.JWT.TokenDuration = 1 * time.Hour

	// Set up server config
	cfg.Server.Port = "8080"
	cfg.Server.ReadTimeout = 5 * time.Second
	cfg.Server.WriteTimeout = 10 * time.Second
	cfg.Server.IdleTimeout = 120 * time.Second
	cfg.Server.ShutdownTimeout = 5 * time.Second

	// Set up log config
	cfg.Log.Level = "debug"
	cfg.Log.Development = true

	return cfg
}

// setupPostgreSQLTestContainer creates a test container with PostgreSQL database
func setupPostgreSQLTestContainer(t *testing.T) *di.Container {
	// Create logger using zaptest
	logger := zaptest.NewLogger(t)

	// Create context
	ctx := context.Background()

	// Create config
	cfg := setupPostgreSQLTestConfig()

	// Log the PostgreSQL DSN for debugging
	t.Logf("Using PostgreSQL DSN: %s", cfg.Database.Postgres.DSN)

	// Create container
	container, err := di.NewContainer(ctx, logger, cfg)
	require.NoError(t, err)

	// Clean up the test database before running tests
	cleanupPostgreSQLTestData(t, cfg.Database.Postgres.DSN)

	return container
}

// cleanupPostgreSQLTestData removes all test data from the PostgreSQL database
// and ensures the table exists with the correct schema
func cleanupPostgreSQLTestData(t *testing.T, dsn string) {
	// Create a PostgreSQL connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// First, drop the existing table if it exists
	_, err = pool.Exec(ctx, "DROP TABLE IF EXISTS families")
	if err != nil {
		t.Fatalf("Failed to drop families table: %v", err)
	}

	// Then, create the table with the correct schema
	createTableSQL := `
	CREATE TABLE families (
		id VARCHAR(36) PRIMARY KEY,
		status VARCHAR(20) NOT NULL,
		parents JSONB NOT NULL,
		children JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX idx_families_status ON families(status);
	CREATE INDEX idx_families_parents ON families USING GIN (parents);
	CREATE INDEX idx_families_children ON families USING GIN (children);

	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER update_families_updated_at
	BEFORE UPDATE ON families
	FOR EACH ROW
	EXECUTE FUNCTION update_updated_at_column();
	`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create families table: %v", err)
	}
}

// setupPostgreSQLTestServer creates a test GraphQL server
func setupPostgreSQLTestServer(t *testing.T, container *di.Container) *httptest.Server {
	// Create resolver using the container
	resolverObj := resolver.NewResolver(
		container.GetFamilyApplicationService(),
		container.GetFamilyMapper(),
	)

	// Create GraphQL handler
	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolverObj,
		Directives: generated.DirectiveRoot{
			IsAuthorized: resolverObj.IsAuthorized,
		},
	}))

	// Create a custom response writer that intercepts 401 responses
	// and converts them to 200 responses with GraphQL errors
	authMiddleware := container.GetAuthService().Middleware()

	// Create test server with custom auth middleware wrapper
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use a custom response writer to intercept 401 responses
		customWriter := &responseWriterInterceptor{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Apply auth middleware with the custom writer
		authHandler := authMiddleware(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// If we get here, authentication was successful
			gqlHandler.ServeHTTP(rw, req)
		}))

		// Handle the request
		authHandler.ServeHTTP(customWriter, r)

		// If we got a 401, convert it to a 200 with a GraphQL error
	}))

	return server
}

// TestIntegrationCreateFamilyPostgreSQL tests the createFamily mutation with PostgreSQL
func TestIntegrationCreateFamilyPostgreSQL(t *testing.T) {
	// Skip test if PostgreSQL is not available
	if os.Getenv("SKIP_POSTGRESQL_TESTS") == "true" {
		t.Skip("Skipping PostgreSQL tests")
	}

	// Setup
	container := setupPostgreSQLTestContainer(t)
	server := setupPostgreSQLTestServer(t, container)
	defer server.Close()
	defer container.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// GraphQL mutation
	mutation := `
		mutation {
			createFamily(input: {
				id: "00000000-0000-0000-0000-000000000001"
				status: MARRIED
				parents: [
					{
						id: "00000000-0000-0000-0000-000000000002"
						firstName: "John"
						lastName: "Smith"
						birthDate: "1980-01-01T00:00:00Z"
					},
					{
						id: "00000000-0000-0000-0000-000000000003"
						firstName: "Jane"
						lastName: "Smith"
						birthDate: "1982-02-02T00:00:00Z"
					}
				]
				children: [
					{
						id: "00000000-0000-0000-0000-000000000004"
						firstName: "Jimmy"
						lastName: "Smith"
						birthDate: "2010-03-03T00:00:00Z"
					},
					{
						id: "00000000-0000-0000-0000-000000000005"
						firstName: "Sally"
						lastName: "Smith"
						birthDate: "2012-04-04T00:00:00Z"
					}
				]
			}) {
				id
				status
				parents {
					id
					firstName
					lastName
					birthDate
				}
				children {
					id
					firstName
					lastName
					birthDate
				}
				parentCount
				childrenCount
			}
		}
	`

	// Create request
	req, err := http.NewRequest("POST", server.URL, strings.NewReader(fmt.Sprintf(`{"query": %q}`, mutation)))
	require.NoError(t, err)

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	// Parse response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	// Check for errors
	errors, hasErrors := result["errors"]
	if hasErrors {
		t.Fatalf("GraphQL errors: %v", errors)
	}

	// Check data
	data, hasData := result["data"].(map[string]interface{})
	require.True(t, hasData, "No data in response")

	// Check createFamily
	createFamily, hasCreateFamily := data["createFamily"].(map[string]interface{})
	require.True(t, hasCreateFamily, "No createFamily in response")

	// Check ID
	id, hasID := createFamily["id"].(string)
	require.True(t, hasID, "No ID in response")
	assert.NotEmpty(t, id, "ID should not be empty")

	// Check status
	status, hasStatus := createFamily["status"].(string)
	require.True(t, hasStatus, "No status in response")
	assert.Equal(t, "MARRIED", status, "Status should be MARRIED")

	// Check parents
	parents, hasParents := createFamily["parents"].([]interface{})
	require.True(t, hasParents, "No parents in response")
	assert.Len(t, parents, 2, "Should have 2 parents")

	// Check first parent
	parent1 := parents[0].(map[string]interface{})
	assert.NotEmpty(t, parent1["id"], "Parent ID should not be empty")
	assert.Equal(t, "John", parent1["firstName"], "First name should be John")
	assert.Equal(t, "Smith", parent1["lastName"], "Last name should be Smith")
	assert.Equal(t, "1980-01-01T00:00:00Z", parent1["birthDate"], "Birth date should be 1980-01-01T00:00:00Z")

	// Check second parent
	parent2 := parents[1].(map[string]interface{})
	assert.NotEmpty(t, parent2["id"], "Parent ID should not be empty")
	assert.Equal(t, "Jane", parent2["firstName"], "First name should be Jane")
	assert.Equal(t, "Smith", parent2["lastName"], "Last name should be Smith")
	assert.Equal(t, "1982-02-02T00:00:00Z", parent2["birthDate"], "Birth date should be 1982-02-02T00:00:00Z")

	// Check children
	children, hasChildren := createFamily["children"].([]interface{})
	require.True(t, hasChildren, "No children in response")
	assert.Len(t, children, 2, "Should have 2 children")

	// Check first child
	child1 := children[0].(map[string]interface{})
	assert.NotEmpty(t, child1["id"], "Child ID should not be empty")
	assert.Equal(t, "Jimmy", child1["firstName"], "First name should be Jimmy")
	assert.Equal(t, "Smith", child1["lastName"], "Last name should be Smith")
	assert.Equal(t, "2010-03-03T00:00:00Z", child1["birthDate"], "Birth date should be 2010-03-03T00:00:00Z")

	// Check second child
	child2 := children[1].(map[string]interface{})
	assert.NotEmpty(t, child2["id"], "Child ID should not be empty")
	assert.Equal(t, "Sally", child2["firstName"], "First name should be Sally")
	assert.Equal(t, "Smith", child2["lastName"], "Last name should be Smith")
	assert.Equal(t, "2012-04-04T00:00:00Z", child2["birthDate"], "Birth date should be 2012-04-04T00:00:00Z")

	// Check counts
	parentCount, hasParentCount := createFamily["parentCount"].(float64)
	require.True(t, hasParentCount, "No parentCount in response")
	assert.Equal(t, float64(2), parentCount, "Parent count should be 2")

	childrenCount, hasChildrenCount := createFamily["childrenCount"].(float64)
	require.True(t, hasChildrenCount, "No childrenCount in response")
	assert.Equal(t, float64(2), childrenCount, "Children count should be 2")
}
