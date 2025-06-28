//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
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
