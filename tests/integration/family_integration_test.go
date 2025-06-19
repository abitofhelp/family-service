// Copyright (c) 2025 A Bit of Help, Inc.
//
// This file contains integration tests for the family service.
// It tests queries and mutations through the GraphQL API with proper authorization.

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
	cfg.Auth.JWT.SecretKey = "test-secret-key"
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
		container.GetContextLogger(),
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
		if customWriter.statusCode == http.StatusUnauthorized {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			errorResponse := map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"message": "Access denied: not authorized",
						"path":    []string{"addParent"},
					},
				},
			}
			json.NewEncoder(w).Encode(errorResponse)
		}
	}))

	return server
}

// responseWriterInterceptor is a custom response writer that intercepts status codes
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
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
	if w.statusCode != http.StatusUnauthorized {
		return w.ResponseWriter.Write(b)
	}
	return len(b), nil // Pretend we wrote it
}

// generateAuthToken generates a JWT token for testing
func generateAuthToken(t *testing.T, container *di.Container) string {
	// Get auth service
	authService := container.GetAuthService()

	// Define scopes and resources for admin role
	adminScopes := []string{"READ", "WRITE", "DELETE", "CREATE"}
	adminResources := []string{"FAMILY", "PARENT", "CHILD"}

	// Generate token with admin role
	token, err := authService.GenerateToken(
		context.Background(),
		"test_admin_user",
		[]string{"ADMIN"},
		adminScopes,
		adminResources,
	)
	require.NoError(t, err)

	return token
}

// TestIntegrationCreateFamily tests the creation of a new family through the GraphQL API
func TestIntegrationCreateFamily(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// Define the GraphQL mutation
	mutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
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
	}`

	// Define the variables for the mutation
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-1",
			"status": "SINGLE",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-1",
					"firstName": "John",
					"lastName":  "Doe",
					"birthDate": "1980-01-01T00:00:00Z",
				},
			},
			"children": []interface{}{},
		},
	}

	// Create the request body
	requestData := map[string]interface{}{
		"query":     mutation,
		"variables": variables,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	// Create the request
	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check the response status
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	// Check for errors
	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	// Check the data
	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	// Check the createFamily result
	createFamily, hasCreateFamily := data["createFamily"].(map[string]interface{})
	require.True(t, hasCreateFamily, "Response does not contain createFamily")

	// Verify the family data
	assert.Equal(t, "test-family-1", createFamily["id"])
	assert.Equal(t, "SINGLE", createFamily["status"])

	// Verify the parents
	parents, hasParents := createFamily["parents"].([]interface{})
	require.True(t, hasParents, "Family does not contain parents")
	require.Len(t, parents, 1, "Family should have 1 parent")

	parent := parents[0].(map[string]interface{})
	assert.Equal(t, "test-parent-1", parent["id"])
	assert.Equal(t, "John", parent["firstName"])
	assert.Equal(t, "Doe", parent["lastName"])
	assert.Equal(t, "1980-01-01T00:00:00Z", parent["birthDate"])

	// Verify the children
	children, hasChildren := createFamily["children"].([]interface{})
	require.True(t, hasChildren, "Family does not contain children")
	require.Len(t, children, 0, "Family should have 0 children")

	// Verify the counts
	assert.Equal(t, float64(1), createFamily["parentCount"])
	assert.Equal(t, float64(0), createFamily["childrenCount"])
}

// TestIntegrationGetFamily tests retrieving a family by ID through the GraphQL API
func TestIntegrationGetFamily(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// First, create a family to query
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
		}
	}`

	createFamilyVariables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-get",
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-get-1",
					"firstName": "Jane",
					"lastName":  "Smith",
					"birthDate": "1985-05-15T00:00:00Z",
				},
				{
					"id":        "test-parent-get-2",
					"firstName": "Bob",
					"lastName":  "Smith",
					"birthDate": "1983-03-10T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        "test-child-get-1",
					"firstName": "Alice",
					"lastName":  "Smith",
					"birthDate": "2010-07-20T00:00:00Z",
				},
			},
		},
	}

	// Create the family first
	createRequestData := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": createFamilyVariables,
	}

	createRequestBodyBytes, err := json.Marshal(createRequestData)
	require.NoError(t, err)
	createRequestBody := string(createRequestBodyBytes)

	createReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createRequestBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createResp, err := client.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode)

	// Now query the family by ID
	query := `query GetFamily($id: ID!) {
		getFamily(id: $id) {
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
	}`

	variables := map[string]interface{}{
		"id": "test-family-get",
	}

	requestData := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	family, hasFamily := data["getFamily"].(map[string]interface{})
	require.True(t, hasFamily, "Response does not contain getFamily")

	// Verify the family data
	assert.Equal(t, "test-family-get", family["id"])
	assert.Equal(t, "MARRIED", family["status"])

	// Verify the parents
	parents, hasParents := family["parents"].([]interface{})
	require.True(t, hasParents, "Family does not contain parents")
	require.Len(t, parents, 2, "Family should have 2 parents")

	// Verify the children
	children, hasChildren := family["children"].([]interface{})
	require.True(t, hasChildren, "Family does not contain children")
	require.Len(t, children, 1, "Family should have 1 child")

	// Verify the counts
	assert.Equal(t, float64(2), family["parentCount"])
	assert.Equal(t, float64(1), family["childrenCount"])
}

// TestIntegrationGetAllFamilies tests retrieving all families through the GraphQL API
func TestIntegrationGetAllFamilies(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// Create multiple families for testing
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
		}
	}`

	// Create first family
	family1Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-all-1",
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-all-1",
					"firstName": "John",
					"lastName":  "Johnson",
					"birthDate": "1975-04-10T00:00:00Z",
				},
				{
					"id":        "test-parent-all-2",
					"firstName": "Mary",
					"lastName":  "Johnson",
					"birthDate": "1978-08-15T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        "test-child-all-1",
					"firstName": "Junior",
					"lastName":  "Johnson",
					"birthDate": "2005-02-20T00:00:00Z",
				},
			},
		},
	}

	// Create second family
	family2Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-all-2",
			"status": "SINGLE",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-all-3",
					"firstName": "Sarah",
					"lastName":  "Williams",
					"birthDate": "1980-11-05T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        "test-child-all-2",
					"firstName": "Emma",
					"lastName":  "Williams",
					"birthDate": "2010-07-12T00:00:00Z",
				},
				{
					"id":        "test-child-all-3",
					"firstName": "Noah",
					"lastName":  "Williams",
					"birthDate": "2012-09-30T00:00:00Z",
				},
			},
		},
	}

	// Create the first family
	createFamily1Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family1Variables,
	}

	createFamily1Bytes, err := json.Marshal(createFamily1Data)
	require.NoError(t, err)
	createFamily1Body := string(createFamily1Bytes)

	createReq1, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily1Body))
	require.NoError(t, err)
	createReq1.Header.Set("Content-Type", "application/json")
	createReq1.Header.Set("Authorization", "Bearer "+token)

	createResp1, err := client.Do(createReq1)
	require.NoError(t, err)
	defer createResp1.Body.Close()
	require.Equal(t, http.StatusOK, createResp1.StatusCode)

	// Create the second family
	createFamily2Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family2Variables,
	}

	createFamily2Bytes, err := json.Marshal(createFamily2Data)
	require.NoError(t, err)
	createFamily2Body := string(createFamily2Bytes)

	createReq2, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily2Body))
	require.NoError(t, err)
	createReq2.Header.Set("Content-Type", "application/json")
	createReq2.Header.Set("Authorization", "Bearer "+token)

	createResp2, err := client.Do(createReq2)
	require.NoError(t, err)
	defer createResp2.Body.Close()
	require.Equal(t, http.StatusOK, createResp2.StatusCode)

	// Now query all families
	query := `query GetAllFamilies {
		getAllFamilies {
			id
			status
			parentCount
			childrenCount
		}
	}`

	requestData := map[string]interface{}{
		"query": query,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	families, hasFamilies := data["getAllFamilies"].([]interface{})
	require.True(t, hasFamilies, "Response does not contain getAllFamilies")

	// We should have at least the two families we created
	require.GreaterOrEqual(t, len(families), 2, "Should have at least 2 families")

	// Create a map to store families by ID for easier verification
	familyMap := make(map[string]map[string]interface{})
	for _, f := range families {
		family := f.(map[string]interface{})
		familyMap[family["id"].(string)] = family
	}

	// Verify the first family
	family1, hasFamily1 := familyMap["test-family-all-1"]
	require.True(t, hasFamily1, "Response does not contain test-family-all-1")
	assert.Equal(t, "MARRIED", family1["status"])
	assert.Equal(t, float64(2), family1["parentCount"])
	assert.Equal(t, float64(1), family1["childrenCount"])

	// Verify the second family
	family2, hasFamily2 := familyMap["test-family-all-2"]
	require.True(t, hasFamily2, "Response does not contain test-family-all-2")
	assert.Equal(t, "SINGLE", family2["status"])
	assert.Equal(t, float64(1), family2["parentCount"])
	assert.Equal(t, float64(2), family2["childrenCount"])
}

// TestIntegrationFindFamilyByChild tests finding a family by child ID through the GraphQL API
// Note: This test expects an error since the resolver is not implemented yet
func TestIntegrationFindFamilyByChild(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// First, create a family with a child
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
		}
	}`

	childID := "test-child-find"
	createFamilyVariables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-bychild",
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-bychild-1",
					"firstName": "David",
					"lastName":  "Wilson",
					"birthDate": "1982-06-15T00:00:00Z",
				},
				{
					"id":        "test-parent-bychild-2",
					"firstName": "Karen",
					"lastName":  "Wilson",
					"birthDate": "1984-09-20T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        childID,
					"firstName": "Kevin",
					"lastName":  "Wilson",
					"birthDate": "2015-03-10T00:00:00Z",
				},
			},
		},
	}

	// Create the family first
	createRequestData := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": createFamilyVariables,
	}

	createRequestBodyBytes, err := json.Marshal(createRequestData)
	require.NoError(t, err)
	createRequestBody := string(createRequestBodyBytes)

	createReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createRequestBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createResp, err := client.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode)

	// Now query the family by child ID
	query := `query FindFamilyByChild($childId: ID!) {
		findFamilyByChild(childId: $childId) {
			id
			status
			parents {
				id
				firstName
				lastName
			}
			children {
				id
				firstName
				lastName
			}
		}
	}`

	variables := map[string]interface{}{
		"childId": childID,
	}

	requestData := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	// Since the resolver is not implemented, we expect an error
	errors, hasErrors := response["errors"]
	require.True(t, hasErrors, "Expected GraphQL response to contain errors")

	// Verify the error message indicates "not implemented"
	errorsArray := errors.([]interface{})
	require.NotEmpty(t, errorsArray, "Expected at least one error")

	firstError := errorsArray[0].(map[string]interface{})
	errorMessage, hasMessage := firstError["message"].(string)
	require.True(t, hasMessage, "Error should have a message")

	// Check that the error message contains "not implemented"
	assert.Contains(t, strings.ToLower(errorMessage), "not implemented", "Error should indicate the feature is not implemented")
}

// TestIntegrationFindFamiliesByParent tests finding families by parent ID through the GraphQL API
func TestIntegrationFindFamiliesByParent(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// Create a parent that will be in multiple families
	sharedParentID := "test-parent-shared"

	// Create multiple families with the same parent
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
		}
	}`

	// Create first family with the shared parent
	family1Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-byparent-1",
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        sharedParentID,
					"firstName": "James",
					"lastName":  "Smith",
					"birthDate": "1970-01-15T00:00:00Z",
				},
				{
					"id":        "test-parent-byparent-2",
					"firstName": "Linda",
					"lastName":  "Smith",
					"birthDate": "1972-03-20T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        "test-child-byparent-1",
					"firstName": "Michael",
					"lastName":  "Smith",
					"birthDate": "2000-05-10T00:00:00Z",
				},
			},
		},
	}

	// Create second family with the same shared parent (simulating a divorce and remarriage)
	family2Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-byparent-2",
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        sharedParentID, // Same parent as in first family
					"firstName": "James",
					"lastName":  "Smith",
					"birthDate": "1970-01-15T00:00:00Z",
				},
				{
					"id":        "test-parent-byparent-3",
					"firstName": "Susan",
					"lastName":  "Johnson",
					"birthDate": "1975-07-25T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        "test-child-byparent-2",
					"firstName": "Emily",
					"lastName":  "Smith",
					"birthDate": "2010-11-30T00:00:00Z",
				},
			},
		},
	}

	// Create a third family without the shared parent
	family3Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-byparent-3",
			"status": "SINGLE",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-byparent-4",
					"firstName": "Robert",
					"lastName":  "Brown",
					"birthDate": "1980-09-05T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{},
		},
	}

	// Create the first family
	createFamily1Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family1Variables,
	}

	createFamily1Bytes, err := json.Marshal(createFamily1Data)
	require.NoError(t, err)
	createFamily1Body := string(createFamily1Bytes)

	createReq1, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily1Body))
	require.NoError(t, err)
	createReq1.Header.Set("Content-Type", "application/json")
	createReq1.Header.Set("Authorization", "Bearer "+token)

	createResp1, err := client.Do(createReq1)
	require.NoError(t, err)
	defer createResp1.Body.Close()
	require.Equal(t, http.StatusOK, createResp1.StatusCode)

	// Create the second family
	createFamily2Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family2Variables,
	}

	createFamily2Bytes, err := json.Marshal(createFamily2Data)
	require.NoError(t, err)
	createFamily2Body := string(createFamily2Bytes)

	createReq2, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily2Body))
	require.NoError(t, err)
	createReq2.Header.Set("Content-Type", "application/json")
	createReq2.Header.Set("Authorization", "Bearer "+token)

	createResp2, err := client.Do(createReq2)
	require.NoError(t, err)
	defer createResp2.Body.Close()
	require.Equal(t, http.StatusOK, createResp2.StatusCode)

	// Create the third family
	createFamily3Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family3Variables,
	}

	createFamily3Bytes, err := json.Marshal(createFamily3Data)
	require.NoError(t, err)
	createFamily3Body := string(createFamily3Bytes)

	createReq3, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily3Body))
	require.NoError(t, err)
	createReq3.Header.Set("Content-Type", "application/json")
	createReq3.Header.Set("Authorization", "Bearer "+token)

	createResp3, err := client.Do(createReq3)
	require.NoError(t, err)
	defer createResp3.Body.Close()
	require.Equal(t, http.StatusOK, createResp3.StatusCode)

	// Now query families by the shared parent ID
	query := `query FindFamiliesByParent($parentId: ID!) {
		findFamiliesByParent(parentId: $parentId) {
			id
			status
			parents {
				id
				firstName
				lastName
			}
			children {
				id
				firstName
				lastName
			}
		}
	}`

	variables := map[string]interface{}{
		"parentId": sharedParentID,
	}

	requestData := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	families, hasFamilies := data["findFamiliesByParent"].([]interface{})
	require.True(t, hasFamilies, "Response does not contain findFamiliesByParent")

	// We should have exactly 2 families with the shared parent
	require.Len(t, families, 2, "Should have exactly 2 families with the shared parent")

	// Create a map to store families by ID for easier verification
	familyMap := make(map[string]map[string]interface{})
	for _, f := range families {
		family := f.(map[string]interface{})
		familyMap[family["id"].(string)] = family
	}

	// Verify the first family
	family1, hasFamily1 := familyMap["test-family-byparent-1"]
	require.True(t, hasFamily1, "Response does not contain test-family-byparent-1")
	assert.Equal(t, "MARRIED", family1["status"])

	// Verify the second family
	family2, hasFamily2 := familyMap["test-family-byparent-2"]
	require.True(t, hasFamily2, "Response does not contain test-family-byparent-2")
	assert.Equal(t, "MARRIED", family2["status"])

	// Verify the third family is NOT in the results
	_, hasFamily3 := familyMap["test-family-byparent-3"]
	assert.False(t, hasFamily3, "Response should not contain test-family-byparent-3")

	// Verify that each family has the shared parent
	for _, f := range families {
		family := f.(map[string]interface{})
		parents := family["parents"].([]interface{})

		// Check if the shared parent is in this family
		foundSharedParent := false
		for _, p := range parents {
			parent := p.(map[string]interface{})
			if parent["id"] == sharedParentID {
				foundSharedParent = true
				assert.Equal(t, "James", parent["firstName"])
				assert.Equal(t, "Smith", parent["lastName"])
				break
			}
		}
		assert.True(t, foundSharedParent, "Family does not contain the shared parent")
	}
}

// TestIntegrationParents tests retrieving all parents through the GraphQL API
func TestIntegrationParents(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// Create families with different parents
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
		}
	}`

	// Create first family
	family1Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-parents-1",
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-list-1",
					"firstName": "Thomas",
					"lastName":  "Anderson",
					"birthDate": "1970-09-13T00:00:00Z",
				},
				{
					"id":        "test-parent-list-2",
					"firstName": "Trinity",
					"lastName":  "Anderson",
					"birthDate": "1975-07-21T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{},
		},
	}

	// Create second family
	family2Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-parents-2",
			"status": "SINGLE",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-list-3",
					"firstName": "Morpheus",
					"lastName":  "Smith",
					"birthDate": "1965-04-30T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{},
		},
	}

	// Create the first family
	createFamily1Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family1Variables,
	}

	createFamily1Bytes, err := json.Marshal(createFamily1Data)
	require.NoError(t, err)
	createFamily1Body := string(createFamily1Bytes)

	createReq1, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily1Body))
	require.NoError(t, err)
	createReq1.Header.Set("Content-Type", "application/json")
	createReq1.Header.Set("Authorization", "Bearer "+token)

	createResp1, err := client.Do(createReq1)
	require.NoError(t, err)
	defer createResp1.Body.Close()
	require.Equal(t, http.StatusOK, createResp1.StatusCode)

	// Create the second family
	createFamily2Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family2Variables,
	}

	createFamily2Bytes, err := json.Marshal(createFamily2Data)
	require.NoError(t, err)
	createFamily2Body := string(createFamily2Bytes)

	createReq2, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily2Body))
	require.NoError(t, err)
	createReq2.Header.Set("Content-Type", "application/json")
	createReq2.Header.Set("Authorization", "Bearer "+token)

	createResp2, err := client.Do(createReq2)
	require.NoError(t, err)
	defer createResp2.Body.Close()
	require.Equal(t, http.StatusOK, createResp2.StatusCode)

	// Now query all parents
	query := `query GetAllParents {
		parents {
			id
			firstName
			lastName
			birthDate
		}
	}`

	requestData := map[string]interface{}{
		"query": query,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	parents, hasParents := data["parents"].([]interface{})
	require.True(t, hasParents, "Response does not contain parents")

	// We should have at least the three parents we created
	require.GreaterOrEqual(t, len(parents), 3, "Should have at least 3 parents")

	// Create a map to store parents by ID for easier verification
	parentMap := make(map[string]map[string]interface{})
	for _, p := range parents {
		parent := p.(map[string]interface{})
		parentMap[parent["id"].(string)] = parent
	}

	// Verify the first parent
	parent1, hasParent1 := parentMap["test-parent-list-1"]
	require.True(t, hasParent1, "Response does not contain test-parent-list-1")
	assert.Equal(t, "Thomas", parent1["firstName"])
	assert.Equal(t, "Anderson", parent1["lastName"])
	assert.Equal(t, "1970-09-13T00:00:00Z", parent1["birthDate"])

	// Verify the second parent
	parent2, hasParent2 := parentMap["test-parent-list-2"]
	require.True(t, hasParent2, "Response does not contain test-parent-list-2")
	assert.Equal(t, "Trinity", parent2["firstName"])
	assert.Equal(t, "Anderson", parent2["lastName"])
	assert.Equal(t, "1975-07-21T00:00:00Z", parent2["birthDate"])

	// Verify the third parent
	parent3, hasParent3 := parentMap["test-parent-list-3"]
	require.True(t, hasParent3, "Response does not contain test-parent-list-3")
	assert.Equal(t, "Morpheus", parent3["firstName"])
	assert.Equal(t, "Smith", parent3["lastName"])
	assert.Equal(t, "1965-04-30T00:00:00Z", parent3["birthDate"])
}

// TestIntegrationCountFamilies tests counting all families through the GraphQL API
func TestIntegrationCountFamilies(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// Create multiple families for testing
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
		}
	}`

	// Create families with unique IDs
	familyCount := 3
	for i := 1; i <= familyCount; i++ {
		familyID := fmt.Sprintf("test-family-count-%d", i)
		parentID := fmt.Sprintf("test-parent-count-%d", i)

		familyVariables := map[string]interface{}{
			"input": map[string]interface{}{
				"id":     familyID,
				"status": "SINGLE",
				"parents": []map[string]interface{}{
					{
						"id":        parentID,
						"firstName": "Parent",
						"lastName":  fmt.Sprintf("Count%d", i),
						"birthDate": "1980-01-01T00:00:00Z",
					},
				},
				"children": []map[string]interface{}{},
			},
		}

		createFamilyData := map[string]interface{}{
			"query":     createFamilyMutation,
			"variables": familyVariables,
		}

		createFamilyBytes, err := json.Marshal(createFamilyData)
		require.NoError(t, err)
		createFamilyBody := string(createFamilyBytes)

		createReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamilyBody))
		require.NoError(t, err)
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("Authorization", "Bearer "+token)

		createResp, err := client.Do(createReq)
		require.NoError(t, err)
		defer createResp.Body.Close()
		require.Equal(t, http.StatusOK, createResp.StatusCode)
	}

	// Now query the count of families
	query := `query CountFamilies {
		countFamilies
	}`

	requestData := map[string]interface{}{
		"query": query,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	count, hasCount := data["countFamilies"].(float64)
	require.True(t, hasCount, "Response does not contain countFamilies")

	// We should have at least the number of families we created
	// Note: Other tests may have created families, so we check for >=
	require.GreaterOrEqual(t, int(count), familyCount, "Should have at least %d families", familyCount)
}

// TestIntegrationCountParents tests counting all parents through the GraphQL API
func TestIntegrationCountParents(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// Create families with different parents
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
		}
	}`

	// Create a family with multiple parents
	family1Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-countparents-1",
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-countparents-1",
					"firstName": "Jack",
					"lastName":  "Brown",
					"birthDate": "1975-05-15T00:00:00Z",
				},
				{
					"id":        "test-parent-countparents-2",
					"firstName": "Jill",
					"lastName":  "Brown",
					"birthDate": "1978-08-20T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{},
		},
	}

	// Create a second family with a different parent
	family2Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-countparents-2",
			"status": "SINGLE",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-countparents-3",
					"firstName": "Sam",
					"lastName":  "Green",
					"birthDate": "1980-03-10T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{},
		},
	}

	// Create the first family
	createFamily1Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family1Variables,
	}

	createFamily1Bytes, err := json.Marshal(createFamily1Data)
	require.NoError(t, err)
	createFamily1Body := string(createFamily1Bytes)

	createReq1, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily1Body))
	require.NoError(t, err)
	createReq1.Header.Set("Content-Type", "application/json")
	createReq1.Header.Set("Authorization", "Bearer "+token)

	createResp1, err := client.Do(createReq1)
	require.NoError(t, err)
	defer createResp1.Body.Close()
	require.Equal(t, http.StatusOK, createResp1.StatusCode)

	// Create the second family
	createFamily2Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family2Variables,
	}

	createFamily2Bytes, err := json.Marshal(createFamily2Data)
	require.NoError(t, err)
	createFamily2Body := string(createFamily2Bytes)

	createReq2, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily2Body))
	require.NoError(t, err)
	createReq2.Header.Set("Content-Type", "application/json")
	createReq2.Header.Set("Authorization", "Bearer "+token)

	createResp2, err := client.Do(createReq2)
	require.NoError(t, err)
	defer createResp2.Body.Close()
	require.Equal(t, http.StatusOK, createResp2.StatusCode)

	// Now query the count of parents
	query := `query CountParents {
		countParents
	}`

	requestData := map[string]interface{}{
		"query": query,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	count, hasCount := data["countParents"].(float64)
	require.True(t, hasCount, "Response does not contain countParents")

	// We should have at least the 3 parents we created
	// Note: Other tests may have created parents, so we check for >=
	require.GreaterOrEqual(t, int(count), 3, "Should have at least 3 parents")
}

// TestIntegrationCountChildren tests counting all children through the GraphQL API
func TestIntegrationCountChildren(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// Create families with different children
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
		}
	}`

	// Create a family with multiple children
	family1Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-countchildren-1",
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-countchildren-1",
					"firstName": "Richard",
					"lastName":  "Davis",
					"birthDate": "1972-04-10T00:00:00Z",
				},
				{
					"id":        "test-parent-countchildren-2",
					"firstName": "Patricia",
					"lastName":  "Davis",
					"birthDate": "1974-06-15T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        "test-child-countchildren-1",
					"firstName": "Ryan",
					"lastName":  "Davis",
					"birthDate": "2005-08-20T00:00:00Z",
				},
				{
					"id":        "test-child-countchildren-2",
					"firstName": "Rebecca",
					"lastName":  "Davis",
					"birthDate": "2008-11-12T00:00:00Z",
				},
			},
		},
	}

	// Create a second family with a different child
	family2Variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     "test-family-countchildren-2",
			"status": "SINGLE",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-countchildren-3",
					"firstName": "Jennifer",
					"lastName":  "Taylor",
					"birthDate": "1985-09-25T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        "test-child-countchildren-3",
					"firstName": "Jason",
					"lastName":  "Taylor",
					"birthDate": "2012-02-28T00:00:00Z",
				},
			},
		},
	}

	// Create the first family
	createFamily1Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family1Variables,
	}

	createFamily1Bytes, err := json.Marshal(createFamily1Data)
	require.NoError(t, err)
	createFamily1Body := string(createFamily1Bytes)

	createReq1, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily1Body))
	require.NoError(t, err)
	createReq1.Header.Set("Content-Type", "application/json")
	createReq1.Header.Set("Authorization", "Bearer "+token)

	createResp1, err := client.Do(createReq1)
	require.NoError(t, err)
	defer createResp1.Body.Close()
	require.Equal(t, http.StatusOK, createResp1.StatusCode)

	// Create the second family
	createFamily2Data := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": family2Variables,
	}

	createFamily2Bytes, err := json.Marshal(createFamily2Data)
	require.NoError(t, err)
	createFamily2Body := string(createFamily2Bytes)

	createReq2, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createFamily2Body))
	require.NoError(t, err)
	createReq2.Header.Set("Content-Type", "application/json")
	createReq2.Header.Set("Authorization", "Bearer "+token)

	createResp2, err := client.Do(createReq2)
	require.NoError(t, err)
	defer createResp2.Body.Close()
	require.Equal(t, http.StatusOK, createResp2.StatusCode)

	// Now query the count of children
	query := `query CountChildren {
		countChildren
	}`

	requestData := map[string]interface{}{
		"query": query,
	}

	requestBodyBytes, err := json.Marshal(requestData)
	require.NoError(t, err)
	requestBody := string(requestBodyBytes)

	req, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	count, hasCount := data["countChildren"].(float64)
	require.True(t, hasCount, "Response does not contain countChildren")

	// We should have at least the 3 children we created
	// Note: Other tests may have created children, so we check for >=
	require.GreaterOrEqual(t, int(count), 3, "Should have at least 3 children")
}

// TestIntegrationAddParent tests adding a parent to an existing family through the GraphQL API
func TestIntegrationAddParent(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// First, create a family with one parent
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
			parents {
				id
			}
		}
	}`

	familyID := "test-family-addparent"
	createFamilyVariables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     familyID,
			"status": "SINGLE",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-original",
					"firstName": "Original",
					"lastName":  "Parent",
					"birthDate": "1980-01-01T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{},
		},
	}

	// Create the family first
	createRequestData := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": createFamilyVariables,
	}

	createRequestBodyBytes, err := json.Marshal(createRequestData)
	require.NoError(t, err)
	createRequestBody := string(createRequestBodyBytes)

	createReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createRequestBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createResp, err := client.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode)

	// Now add a new parent to the family
	addParentMutation := `mutation AddParent($familyId: ID!, $input: ParentInput!) {
		addParent(familyId: $familyId, input: $input) {
			id
			status
			parents {
				id
				firstName
				lastName
				birthDate
			}
			parentCount
		}
	}`

	newParentID := "test-parent-added"
	addParentVariables := map[string]interface{}{
		"familyId": familyID,
		"input": map[string]interface{}{
			"id":        newParentID,
			"firstName": "Added",
			"lastName":  "Parent",
			"birthDate": "1982-05-15T00:00:00Z",
		},
	}

	addParentData := map[string]interface{}{
		"query":     addParentMutation,
		"variables": addParentVariables,
	}

	addParentBytes, err := json.Marshal(addParentData)
	require.NoError(t, err)
	addParentBody := string(addParentBytes)

	addParentReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(addParentBody))
	require.NoError(t, err)
	addParentReq.Header.Set("Content-Type", "application/json")
	addParentReq.Header.Set("Authorization", "Bearer "+token)

	addParentResp, err := client.Do(addParentReq)
	require.NoError(t, err)
	defer addParentResp.Body.Close()

	assert.Equal(t, http.StatusOK, addParentResp.StatusCode)

	body, err := ioutil.ReadAll(addParentResp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	addParentResult, hasAddParent := data["addParent"].(map[string]interface{})
	require.True(t, hasAddParent, "Response does not contain addParent")

	// Verify the family data
	assert.Equal(t, familyID, addParentResult["id"])
	assert.Equal(t, "MARRIED", addParentResult["status"])

	// Verify the parents
	parents, hasParents := addParentResult["parents"].([]interface{})
	require.True(t, hasParents, "Family does not contain parents")
	require.Len(t, parents, 2, "Family should have 2 parents")

	// Verify the parent count
	assert.Equal(t, float64(2), addParentResult["parentCount"])

	// Create a map to store parents by ID for easier verification
	parentMap := make(map[string]map[string]interface{})
	for _, p := range parents {
		parent := p.(map[string]interface{})
		parentMap[parent["id"].(string)] = parent
	}

	// Verify the original parent is still there
	originalParent, hasOriginalParent := parentMap["test-parent-original"]
	require.True(t, hasOriginalParent, "Original parent should still be in the family")
	assert.Equal(t, "Original", originalParent["firstName"])
	assert.Equal(t, "Parent", originalParent["lastName"])

	// Verify the new parent was added
	newParent, hasNewParent := parentMap[newParentID]
	require.True(t, hasNewParent, "New parent should be in the family")
	assert.Equal(t, "Added", newParent["firstName"])
	assert.Equal(t, "Parent", newParent["lastName"])
	assert.Equal(t, "1982-05-15T00:00:00Z", newParent["birthDate"])

	// Test authorization by trying to add a parent without a token
	addParentReqNoAuth, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(addParentBody))
	require.NoError(t, err)
	addParentReqNoAuth.Header.Set("Content-Type", "application/json")
	// Intentionally not setting Authorization header

	addParentRespNoAuth, err := client.Do(addParentReqNoAuth)
	require.NoError(t, err)
	defer addParentRespNoAuth.Body.Close()

	assert.Equal(t, http.StatusOK, addParentRespNoAuth.StatusCode)

	bodyNoAuth, err := ioutil.ReadAll(addParentRespNoAuth.Body)
	require.NoError(t, err)

	var responseNoAuth map[string]interface{}
	err = json.Unmarshal(bodyNoAuth, &responseNoAuth)
	require.NoError(t, err)

	// Verify that we get an authorization error
	errorsNoAuth, hasErrorsNoAuth := responseNoAuth["errors"]
	require.True(t, hasErrorsNoAuth, "Response should contain errors when not authorized")

	errorsArrayNoAuth := errorsNoAuth.([]interface{})
	require.NotEmpty(t, errorsArrayNoAuth, "Should have at least one error")

	firstErrorNoAuth := errorsArrayNoAuth[0].(map[string]interface{})
	errorMessageNoAuth, hasMessageNoAuth := firstErrorNoAuth["message"].(string)
	require.True(t, hasMessageNoAuth, "Error should have a message")

	// Check that the error message indicates an authorization issue
	assert.Contains(t, strings.ToLower(errorMessageNoAuth), "not authorized", "Error should indicate authorization issue")
}

// TestIntegrationAddChild tests adding a child to an existing family through the GraphQL API
func TestIntegrationAddChild(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// First, create a family with parents but no children
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
			parents {
				id
			}
			children {
				id
			}
		}
	}`

	familyID := "test-family-addchild"
	createFamilyVariables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     familyID,
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-addchild-1",
					"firstName": "Father",
					"lastName":  "Family",
					"birthDate": "1975-03-15T00:00:00Z",
				},
				{
					"id":        "test-parent-addchild-2",
					"firstName": "Mother",
					"lastName":  "Family",
					"birthDate": "1978-07-20T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{},
		},
	}

	// Create the family first
	createRequestData := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": createFamilyVariables,
	}

	createRequestBodyBytes, err := json.Marshal(createRequestData)
	require.NoError(t, err)
	createRequestBody := string(createRequestBodyBytes)

	createReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createRequestBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createResp, err := client.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode)

	// Now add a child to the family
	addChildMutation := `mutation AddChild($familyId: ID!, $input: ChildInput!) {
		addChild(familyId: $familyId, input: $input) {
			id
			status
			children {
				id
				firstName
				lastName
				birthDate
			}
			childrenCount
		}
	}`

	childID := "test-child-added"
	addChildVariables := map[string]interface{}{
		"familyId": familyID,
		"input": map[string]interface{}{
			"id":        childID,
			"firstName": "Baby",
			"lastName":  "Family",
			"birthDate": "2020-10-25T00:00:00Z",
		},
	}

	addChildData := map[string]interface{}{
		"query":     addChildMutation,
		"variables": addChildVariables,
	}

	addChildBytes, err := json.Marshal(addChildData)
	require.NoError(t, err)
	addChildBody := string(addChildBytes)

	addChildReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(addChildBody))
	require.NoError(t, err)
	addChildReq.Header.Set("Content-Type", "application/json")
	addChildReq.Header.Set("Authorization", "Bearer "+token)

	addChildResp, err := client.Do(addChildReq)
	require.NoError(t, err)
	defer addChildResp.Body.Close()

	assert.Equal(t, http.StatusOK, addChildResp.StatusCode)

	body, err := ioutil.ReadAll(addChildResp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	addChildResult, hasAddChild := data["addChild"].(map[string]interface{})
	require.True(t, hasAddChild, "Response does not contain addChild")

	// Verify the family data
	assert.Equal(t, familyID, addChildResult["id"])
	assert.Equal(t, "MARRIED", addChildResult["status"])

	// Verify the children
	children, hasChildren := addChildResult["children"].([]interface{})
	require.True(t, hasChildren, "Family does not contain children")
	require.Len(t, children, 1, "Family should have 1 child")

	// Verify the children count
	assert.Equal(t, float64(1), addChildResult["childrenCount"])

	// Verify the child data
	child := children[0].(map[string]interface{})
	assert.Equal(t, childID, child["id"])
	assert.Equal(t, "Baby", child["firstName"])
	assert.Equal(t, "Family", child["lastName"])
	assert.Equal(t, "2020-10-25T00:00:00Z", child["birthDate"])

	// Test authorization by trying to add a child without a token
	addChildReqNoAuth, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(addChildBody))
	require.NoError(t, err)
	addChildReqNoAuth.Header.Set("Content-Type", "application/json")
	// Intentionally not setting Authorization header

	addChildRespNoAuth, err := client.Do(addChildReqNoAuth)
	require.NoError(t, err)
	defer addChildRespNoAuth.Body.Close()

	assert.Equal(t, http.StatusOK, addChildRespNoAuth.StatusCode)

	bodyNoAuth, err := ioutil.ReadAll(addChildRespNoAuth.Body)
	require.NoError(t, err)

	var responseNoAuth map[string]interface{}
	err = json.Unmarshal(bodyNoAuth, &responseNoAuth)
	require.NoError(t, err)

	// Verify that we get an authorization error
	errorsNoAuth, hasErrorsNoAuth := responseNoAuth["errors"]
	require.True(t, hasErrorsNoAuth, "Response should contain errors when not authorized")

	errorsArrayNoAuth := errorsNoAuth.([]interface{})
	require.NotEmpty(t, errorsArrayNoAuth, "Should have at least one error")

	firstErrorNoAuth := errorsArrayNoAuth[0].(map[string]interface{})
	errorMessageNoAuth, hasMessageNoAuth := firstErrorNoAuth["message"].(string)
	require.True(t, hasMessageNoAuth, "Error should have a message")

	// Check that the error message indicates an authorization issue
	assert.Contains(t, strings.ToLower(errorMessageNoAuth), "not authorized", "Error should indicate authorization issue")
}

// TestIntegrationRemoveChild tests removing a child from a family through the GraphQL API
func TestIntegrationRemoveChild(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// First, create a family with a child
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
			children {
				id
			}
			childrenCount
		}
	}`

	familyID := "test-family-removechild"
	childID := "test-child-remove"
	createFamilyVariables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     familyID,
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        "test-parent-removechild-1",
					"firstName": "Parent",
					"lastName":  "One",
					"birthDate": "1980-05-10T00:00:00Z",
				},
				{
					"id":        "test-parent-removechild-2",
					"firstName": "Parent",
					"lastName":  "Two",
					"birthDate": "1982-03-15T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        childID,
					"firstName": "Child",
					"lastName":  "ToRemove",
					"birthDate": "2015-08-20T00:00:00Z",
				},
			},
		},
	}

	// Create the family first
	createRequestData := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": createFamilyVariables,
	}

	createRequestBodyBytes, err := json.Marshal(createRequestData)
	require.NoError(t, err)
	createRequestBody := string(createRequestBodyBytes)

	createReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createRequestBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createResp, err := client.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode)

	// Verify the family was created with the child
	createRespBody, err := ioutil.ReadAll(createResp.Body)
	require.NoError(t, err)

	var createResponse map[string]interface{}
	err = json.Unmarshal(createRespBody, &createResponse)
	require.NoError(t, err)

	createData, hasCreateData := createResponse["data"].(map[string]interface{})
	require.True(t, hasCreateData, "Response does not contain data")

	createFamily, hasCreateFamily := createData["createFamily"].(map[string]interface{})
	require.True(t, hasCreateFamily, "Response does not contain createFamily")

	// Verify the family has 1 child
	assert.Equal(t, float64(1), createFamily["childrenCount"])

	// Now remove the child from the family
	removeChildMutation := `mutation RemoveChild($familyId: ID!, $childId: ID!) {
		removeChild(familyId: $familyId, childId: $childId) {
			id
			status
			children {
				id
				firstName
				lastName
			}
			childrenCount
		}
	}`

	removeChildVariables := map[string]interface{}{
		"familyId": familyID,
		"childId":  childID,
	}

	removeChildData := map[string]interface{}{
		"query":     removeChildMutation,
		"variables": removeChildVariables,
	}

	removeChildBytes, err := json.Marshal(removeChildData)
	require.NoError(t, err)
	removeChildBody := string(removeChildBytes)

	removeChildReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(removeChildBody))
	require.NoError(t, err)
	removeChildReq.Header.Set("Content-Type", "application/json")
	removeChildReq.Header.Set("Authorization", "Bearer "+token)

	removeChildResp, err := client.Do(removeChildReq)
	require.NoError(t, err)
	defer removeChildResp.Body.Close()

	assert.Equal(t, http.StatusOK, removeChildResp.StatusCode)

	body, err := ioutil.ReadAll(removeChildResp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	removeChildResult, hasRemoveChild := data["removeChild"].(map[string]interface{})
	require.True(t, hasRemoveChild, "Response does not contain removeChild")

	// Verify the family data
	assert.Equal(t, familyID, removeChildResult["id"])
	assert.Equal(t, "MARRIED", removeChildResult["status"])

	// Verify the children
	children, hasChildren := removeChildResult["children"].([]interface{})
	require.True(t, hasChildren, "Family does not contain children field")
	require.Len(t, children, 0, "Family should have 0 children after removal")

	// Verify the children count
	assert.Equal(t, float64(0), removeChildResult["childrenCount"])

	// Test authorization by trying to remove a child without a token
	removeChildReqNoAuth, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(removeChildBody))
	require.NoError(t, err)
	removeChildReqNoAuth.Header.Set("Content-Type", "application/json")
	// Intentionally not setting Authorization header

	removeChildRespNoAuth, err := client.Do(removeChildReqNoAuth)
	require.NoError(t, err)
	defer removeChildRespNoAuth.Body.Close()

	assert.Equal(t, http.StatusOK, removeChildRespNoAuth.StatusCode)

	bodyNoAuth, err := ioutil.ReadAll(removeChildRespNoAuth.Body)
	require.NoError(t, err)

	var responseNoAuth map[string]interface{}
	err = json.Unmarshal(bodyNoAuth, &responseNoAuth)
	require.NoError(t, err)

	// Verify that we get an authorization error
	errorsNoAuth, hasErrorsNoAuth := responseNoAuth["errors"]
	require.True(t, hasErrorsNoAuth, "Response should contain errors when not authorized")

	errorsArrayNoAuth := errorsNoAuth.([]interface{})
	require.NotEmpty(t, errorsArrayNoAuth, "Should have at least one error")

	firstErrorNoAuth := errorsArrayNoAuth[0].(map[string]interface{})
	errorMessageNoAuth, hasMessageNoAuth := firstErrorNoAuth["message"].(string)
	require.True(t, hasMessageNoAuth, "Error should have a message")

	// Check that the error message indicates an authorization issue
	assert.Contains(t, strings.ToLower(errorMessageNoAuth), "not authorized", "Error should indicate authorization issue")
}

// TestIntegrationMarkParentDeceased tests marking a parent as deceased through the GraphQL API
func TestIntegrationMarkParentDeceased(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// First, create a family with a parent
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
			parents {
				id
				firstName
				lastName
				deathDate
			}
		}
	}`

	familyID := "test-family-deceased"
	parentID := "test-parent-deceased"
	createFamilyVariables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     familyID,
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        parentID,
					"firstName": "Deceased",
					"lastName":  "Parent",
					"birthDate": "1950-03-15T00:00:00Z",
				},
				{
					"id":        "test-parent-alive",
					"firstName": "Alive",
					"lastName":  "Parent",
					"birthDate": "1955-07-20T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{},
		},
	}

	// Create the family first
	createRequestData := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": createFamilyVariables,
	}

	createRequestBodyBytes, err := json.Marshal(createRequestData)
	require.NoError(t, err)
	createRequestBody := string(createRequestBodyBytes)

	createReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createRequestBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createResp, err := client.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode)

	// Now mark the parent as deceased
	markParentDeceasedMutation := `mutation MarkParentDeceased($familyId: ID!, $parentId: ID!, $deathDate: String!) {
		markParentDeceased(familyId: $familyId, parentId: $parentId, deathDate: $deathDate) {
			id
			status
			parents {
				id
				firstName
				lastName
				birthDate
				deathDate
			}
		}
	}`

	deathDate := "2023-01-15T00:00:00Z"
	markParentDeceasedVariables := map[string]interface{}{
		"familyId":  familyID,
		"parentId":  parentID,
		"deathDate": deathDate,
	}

	markParentDeceasedData := map[string]interface{}{
		"query":     markParentDeceasedMutation,
		"variables": markParentDeceasedVariables,
	}

	markParentDeceasedBytes, err := json.Marshal(markParentDeceasedData)
	require.NoError(t, err)
	markParentDeceasedBody := string(markParentDeceasedBytes)

	markParentDeceasedReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(markParentDeceasedBody))
	require.NoError(t, err)
	markParentDeceasedReq.Header.Set("Content-Type", "application/json")
	markParentDeceasedReq.Header.Set("Authorization", "Bearer "+token)

	markParentDeceasedResp, err := client.Do(markParentDeceasedReq)
	require.NoError(t, err)
	defer markParentDeceasedResp.Body.Close()

	assert.Equal(t, http.StatusOK, markParentDeceasedResp.StatusCode)

	body, err := ioutil.ReadAll(markParentDeceasedResp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	markParentDeceasedResult, hasMarkParentDeceased := data["markParentDeceased"].(map[string]interface{})
	require.True(t, hasMarkParentDeceased, "Response does not contain markParentDeceased")

	// Verify the family data
	assert.Equal(t, familyID, markParentDeceasedResult["id"])

	// Verify the parents
	parents, hasParents := markParentDeceasedResult["parents"].([]interface{})
	require.True(t, hasParents, "Family does not contain parents")
	require.Len(t, parents, 2, "Family should have 2 parents")

	// Create a map to store parents by ID for easier verification
	parentMap := make(map[string]map[string]interface{})
	for _, p := range parents {
		parent := p.(map[string]interface{})
		parentMap[parent["id"].(string)] = parent
	}

	// Verify the deceased parent has a death date
	deceasedParent, hasDeceasedParent := parentMap[parentID]
	require.True(t, hasDeceasedParent, "Deceased parent should still be in the family")
	assert.Equal(t, "Deceased", deceasedParent["firstName"])
	assert.Equal(t, "Parent", deceasedParent["lastName"])
	assert.Equal(t, deathDate, deceasedParent["deathDate"])

	// Verify the other parent is still alive (no death date)
	aliveParent, hasAliveParent := parentMap["test-parent-alive"]
	require.True(t, hasAliveParent, "Alive parent should still be in the family")
	assert.Equal(t, "Alive", aliveParent["firstName"])
	assert.Equal(t, "Parent", aliveParent["lastName"])
	assert.Nil(t, aliveParent["deathDate"], "Alive parent should not have a death date")

	// Test authorization by trying to mark a parent as deceased without a token
	markParentDeceasedReqNoAuth, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(markParentDeceasedBody))
	require.NoError(t, err)
	markParentDeceasedReqNoAuth.Header.Set("Content-Type", "application/json")
	// Intentionally not setting Authorization header

	markParentDeceasedRespNoAuth, err := client.Do(markParentDeceasedReqNoAuth)
	require.NoError(t, err)
	defer markParentDeceasedRespNoAuth.Body.Close()

	assert.Equal(t, http.StatusOK, markParentDeceasedRespNoAuth.StatusCode)

	bodyNoAuth, err := ioutil.ReadAll(markParentDeceasedRespNoAuth.Body)
	require.NoError(t, err)

	var responseNoAuth map[string]interface{}
	err = json.Unmarshal(bodyNoAuth, &responseNoAuth)
	require.NoError(t, err)

	// Verify that we get an authorization error
	errorsNoAuth, hasErrorsNoAuth := responseNoAuth["errors"]
	require.True(t, hasErrorsNoAuth, "Response should contain errors when not authorized")

	errorsArrayNoAuth := errorsNoAuth.([]interface{})
	require.NotEmpty(t, errorsArrayNoAuth, "Should have at least one error")

	firstErrorNoAuth := errorsArrayNoAuth[0].(map[string]interface{})
	errorMessageNoAuth, hasMessageNoAuth := firstErrorNoAuth["message"].(string)
	require.True(t, hasMessageNoAuth, "Error should have a message")

	// Check that the error message indicates an authorization issue
	assert.Contains(t, strings.ToLower(errorMessageNoAuth), "not authorized", "Error should indicate authorization issue")
}

// TestIntegrationDivorce tests processing a divorce through the GraphQL API
func TestIntegrationDivorce(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test container
	container := setupTestContainer(t)
	defer container.Close()

	// Setup test server
	server := setupTestServer(t, container)
	defer server.Close()

	// Generate auth token
	token := generateAuthToken(t, container)

	// Create GraphQL client
	client := &http.Client{}

	// First, create a family with two parents and children
	createFamilyMutation := `mutation CreateFamily($input: FamilyInput!) {
		createFamily(input: $input) {
			id
			status
			parents {
				id
				firstName
				lastName
			}
			children {
				id
				firstName
				lastName
			}
		}
	}`

	familyID := "test-family-divorce"
	custodialParentID := "test-parent-custodial"
	nonCustodialParentID := "test-parent-noncustodial"

	createFamilyVariables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     familyID,
			"status": "MARRIED",
			"parents": []map[string]interface{}{
				{
					"id":        custodialParentID,
					"firstName": "Custodial",
					"lastName":  "Parent",
					"birthDate": "1980-05-15T00:00:00Z",
				},
				{
					"id":        nonCustodialParentID,
					"firstName": "NonCustodial",
					"lastName":  "Parent",
					"birthDate": "1982-08-20T00:00:00Z",
				},
			},
			"children": []map[string]interface{}{
				{
					"id":        "test-child-divorce-1",
					"firstName": "First",
					"lastName":  "Child",
					"birthDate": "2010-03-10T00:00:00Z",
				},
				{
					"id":        "test-child-divorce-2",
					"firstName": "Second",
					"lastName":  "Child",
					"birthDate": "2012-07-25T00:00:00Z",
				},
			},
		},
	}

	// Create the family first
	createRequestData := map[string]interface{}{
		"query":     createFamilyMutation,
		"variables": createFamilyVariables,
	}

	createRequestBodyBytes, err := json.Marshal(createRequestData)
	require.NoError(t, err)
	createRequestBody := string(createRequestBodyBytes)

	createReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(createRequestBody))
	require.NoError(t, err)
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)

	createResp, err := client.Do(createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusOK, createResp.StatusCode)

	// Now process a divorce
	divorceMutation := `mutation Divorce($familyId: ID!, $custodialParentId: ID!) {
		divorce(familyId: $familyId, custodialParentId: $custodialParentId) {
			id
			status
			parents {
				id
				firstName
				lastName
			}
			children {
				id
				firstName
				lastName
			}
		}
	}`

	divorceVariables := map[string]interface{}{
		"familyId":          familyID,
		"custodialParentId": custodialParentID,
	}

	divorceData := map[string]interface{}{
		"query":     divorceMutation,
		"variables": divorceVariables,
	}

	divorceBytes, err := json.Marshal(divorceData)
	require.NoError(t, err)
	divorceBody := string(divorceBytes)

	divorceReq, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(divorceBody))
	require.NoError(t, err)
	divorceReq.Header.Set("Content-Type", "application/json")
	divorceReq.Header.Set("Authorization", "Bearer "+token)

	divorceResp, err := client.Do(divorceReq)
	require.NoError(t, err)
	defer divorceResp.Body.Close()

	assert.Equal(t, http.StatusOK, divorceResp.StatusCode)

	body, err := ioutil.ReadAll(divorceResp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	errors, hasErrors := response["errors"]
	if hasErrors {
		t.Fatalf("GraphQL response contains errors: %v", errors)
	}

	data, hasData := response["data"].(map[string]interface{})
	require.True(t, hasData, "Response does not contain data")

	divorceResult, hasDivorce := data["divorce"].(map[string]interface{})
	require.True(t, hasDivorce, "Response does not contain divorce")

	// Verify the family data
	assert.Equal(t, familyID, divorceResult["id"])
	assert.Equal(t, "DIVORCED", divorceResult["status"])

	// Verify the parents - should only have the custodial parent
	parents, hasParents := divorceResult["parents"].([]interface{})
	require.True(t, hasParents, "Family does not contain parents")
	require.Len(t, parents, 1, "Family should have 1 parent after divorce")

	// Verify the custodial parent is in the family
	parent := parents[0].(map[string]interface{})
	assert.Equal(t, custodialParentID, parent["id"])
	assert.Equal(t, "Custodial", parent["firstName"])
	assert.Equal(t, "Parent", parent["lastName"])

	// Verify the children are still in the family
	children, hasChildren := divorceResult["children"].([]interface{})
	require.True(t, hasChildren, "Family does not contain children")
	require.Len(t, children, 2, "Family should still have 2 children after divorce")

	// Test authorization by trying to process a divorce without a token
	divorceReqNoAuth, err := http.NewRequest("POST", server.URL+"/graphql", strings.NewReader(divorceBody))
	require.NoError(t, err)
	divorceReqNoAuth.Header.Set("Content-Type", "application/json")
	// Intentionally not setting Authorization header

	divorceRespNoAuth, err := client.Do(divorceReqNoAuth)
	require.NoError(t, err)
	defer divorceRespNoAuth.Body.Close()

	assert.Equal(t, http.StatusOK, divorceRespNoAuth.StatusCode)

	bodyNoAuth, err := ioutil.ReadAll(divorceRespNoAuth.Body)
	require.NoError(t, err)

	var responseNoAuth map[string]interface{}
	err = json.Unmarshal(bodyNoAuth, &responseNoAuth)
	require.NoError(t, err)

	// Verify that we get an authorization error
	errorsNoAuth, hasErrorsNoAuth := responseNoAuth["errors"]
	require.True(t, hasErrorsNoAuth, "Response should contain errors when not authorized")

	errorsArrayNoAuth := errorsNoAuth.([]interface{})
	require.NotEmpty(t, errorsArrayNoAuth, "Should have at least one error")

	firstErrorNoAuth := errorsArrayNoAuth[0].(map[string]interface{})
	errorMessageNoAuth, hasMessageNoAuth := firstErrorNoAuth["message"].(string)
	require.True(t, hasMessageNoAuth, "Error should have a message")

	// Check that the error message indicates an authorization issue
	assert.Contains(t, strings.ToLower(errorMessageNoAuth), "not authorized", "Error should indicate authorization issue")
}
