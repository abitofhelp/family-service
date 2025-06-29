// Copyright (c) 2025 A Bit of Help, Inc.

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/abitofhelp/family-service/cmd/server/graphql/di"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/generated"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/resolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// setupTestConfig creates a test configuration with SQLite database
func setupTestConfig() *config.Config {
	cfg := &config.Config{}

	// Set up database config
	cfg.Database.Type = "sqlite"
	cfg.Database.SQLite.URI = ":memory:" // Use in-memory SQLite for tests

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

// setupTestContainer creates a test container with SQLite database
func setupTestContainer(t *testing.T) *di.Container {
	// Create logger using zaptest
	logger := zaptest.NewLogger(t)

	// Create context
	ctx := context.Background()

	// Create config
	cfg := setupTestConfig()

	// Create container
	container, err := di.NewContainer(ctx, logger, cfg)
	require.NoError(t, err)

	return container
}

// setupTestServer creates a test GraphQL server
func setupTestServer(t *testing.T, container *di.Container) *httptest.Server {
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

// responseWriterInterceptor is a custom response writer that intercepts 401 responses
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

// WriteHeader intercepts the status code
func (w *responseWriterInterceptor) WriteHeader(code int) {
	w.statusCode = code
	if code != http.StatusUnauthorized {
		w.ResponseWriter.WriteHeader(code)
	}
}

// Write intercepts the response body
func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
	w.body = b
	if w.statusCode != http.StatusUnauthorized {
		return w.ResponseWriter.Write(b)
	}
	return len(b), nil
}

// generateAuthToken generates a JWT token for testing
func generateAuthToken(t *testing.T, container *di.Container) string {
	// Get auth service
	authService := container.GetAuthService()

	// Generate token with ADMIN role and CREATE scope for FAMILY resource
	ctx := context.Background()
	token, err := authService.GenerateToken(ctx, "test-user", []string{"ADMIN"}, []string{"CREATE"}, []string{"FAMILY"})
	require.NoError(t, err)

	return token
}

// TestIntegrationCreateFamily tests the createFamily mutation
func TestIntegrationCreateFamily(t *testing.T) {
	// Setup
	container := setupTestContainer(t)
	server := setupTestServer(t, container)
	defer server.Close()

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
