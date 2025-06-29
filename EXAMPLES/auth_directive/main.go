//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Package main demonstrates how to use GraphQL with authorization directives.
//
// This example shows how to:
// - Create and send GraphQL queries to a server
// - Add JWT authorization tokens to requests
// - Handle GraphQL responses and errors
// - Test both authorized and unauthorized requests
//
// The example demonstrates the @isAuthorized directive in GraphQL, which
// protects queries and mutations based on user authentication status.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GraphQLRequest represents a GraphQL request to be sent to a server.
//
// This struct contains all the necessary fields for a standard GraphQL request:
// - Query: The GraphQL query or mutation string
// - OperationName: The name of the operation to execute (optional)
// - Variables: Any variables needed for the query (optional)
//
// When serialized to JSON, this struct follows the standard GraphQL request format
// expected by GraphQL servers.
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response received from a server.
//
// This struct follows the standard GraphQL response format with:
// - Data: Contains the requested data if the query was successful
// - Errors: Contains any errors that occurred during query execution
//
// The Errors field is a slice of error objects, each containing at minimum
// a message field. When unmarshaling JSON responses, this struct captures
// both successful results and error information.
type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

func main() {
	// URL of the GraphQL server
	url := "http://localhost:8080/graphql"

	// Test with no authorization
	fmt.Println("Testing with no authorization:")
	testQuery(url, "", "GetAllFamilies", `
		query GetAllFamilies {
			getAllFamilies {
				id
				status
			}
		}
	`)

	// Generate a JWT token (this would normally come from your auth system)
	// For testing, you can use the genjwt tool in the repo
	// Or you can set a valid token here if you have one
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		fmt.Println("No AUTH_TOKEN environment variable found. Set it to a valid JWT token for testing with authorization.")
	} else {
		fmt.Println("\nTesting with authorization:")
		testQuery(url, token, "GetAllFamilies", `
			query GetAllFamilies {
				getAllFamilies {
					id
					status
				}
			}
		`)
	}
}

// testQuery sends a GraphQL query to the specified URL with optional authorization.
//
// This function demonstrates how to:
// - Create a GraphQL request with the specified query and operation name
// - Add authorization via JWT token if provided
// - Send the request to the GraphQL server
// - Parse and display the response
//
// Parameters:
//   - url: The URL of the GraphQL endpoint
//   - token: JWT token for authorization (can be empty for unauthorized requests)
//   - operationName: The name of the GraphQL operation to execute
//   - query: The GraphQL query string
//
// The function handles errors at each step of the process and prints
// appropriate messages to the console.
func testQuery(url, token, operationName, query string) {
	// Create the request
	req := GraphQLRequest{
		Query:         query,
		OperationName: operationName,
	}

	// Convert request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Parse response
	var gqlResp GraphQLResponse
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	// Print response
	fmt.Printf("Status: %s\n", resp.Status)
	if len(gqlResp.Errors) > 0 {
		fmt.Println("Errors:")
		for _, err := range gqlResp.Errors {
			fmt.Printf("  - %s\n", err.Message)
		}
	} else if gqlResp.Data != nil {
		fmt.Println("Data received successfully")
	}
}
