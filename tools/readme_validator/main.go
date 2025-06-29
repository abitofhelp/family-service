// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Required sections in component README files
var componentRequiredSections = []string{
	"Overview",
	"Features",
	"Installation",
	"Quick Start",
	"Configuration",
	"API Documentation",
	"Examples",
	"Best Practices",
	"Troubleshooting",
	"Related Components",
	"Contributing",
	"License",
}

// Required sections in example README files
var exampleRequiredSections = []string{
	"Overview",
	"Features",
	"Running the Example",
	"Code Walkthrough",
	"Expected Output",
	"Related Examples",
	"Related Components",
	"License",
}

// validateReadme checks if a README.md file follows the template structure
func validateReadme(filePath string) (bool, []string) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, []string{fmt.Sprintf("Error opening file: %v", err)}
	}
	defer file.Close()

	// Skip validation for the main README.md, EXAMPLES/README.md, and COMPONENT_README_TEMPLATE.md
	if strings.ToLower(filePath) == "readme.md" || 
	   strings.ToLower(filePath) == "examples/readme.md" || 
	   strings.ToLower(filePath) == "docs/tools/scripts/readme.md" || 
	   strings.ToLower(filepath.Base(filePath)) == "component_readme_template.md" {
		return true, nil
	}

	scanner := bufio.NewScanner(file)
	var foundSections []string
	var errors []string
	lineNum := 0

	// Check for title (first line should be a level 1 heading)
	if scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if !strings.HasPrefix(line, "# ") {
			errors = append(errors, fmt.Sprintf("Line %d: Missing title (should start with '# ')", lineNum))
		}
	}

	// Scan for section headings - match any level 2 heading
	sectionRegex := regexp.MustCompile(`^## (.*)$`)
	examplePathRegex := regexp.MustCompile(`\[.*\]\((\.\.)?/[Ee][Xx][Aa][Mm][Pp][Ll][Ee][Ss]/.*\)`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for section headings
		matches := sectionRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			foundSections = append(foundSections, matches[1])
		}

		// Check for inconsistent example paths
		if examplePathRegex.MatchString(line) {
			// Standardize to lowercase /examples/
			if strings.Contains(line, "/EXAMPLES/") {
				errors = append(errors, fmt.Sprintf("Line %d: Inconsistent example path (use lowercase '/examples/')", lineNum))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, []string{fmt.Sprintf("Error reading file: %v", err)}
	}

	// Determine which set of required sections to use based on the file path
	var requiredSections []string
	if strings.Contains(strings.ToLower(filePath), "examples/") && !strings.EqualFold(filePath, "examples/readme.md") {
		// Use example required sections for README.md files in EXAMPLES subdirectories
		requiredSections = exampleRequiredSections
	} else {
		// Use component required sections for all other README.md files
		requiredSections = componentRequiredSections
	}

	// Check for missing required sections
	for _, section := range requiredSections {
		found := false
		for _, foundSection := range foundSections {
			// Case-insensitive check for section name
			if strings.Contains(strings.ToLower(foundSection), strings.ToLower(section)) {
				found = true
				break
			}
		}
		if !found {
			errors = append(errors, fmt.Sprintf("Missing required section: %s", section))
		}
	}

	return len(errors) == 0, errors
}

// findReadmeFiles finds all README.md files in the project
func findReadmeFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(info.Name()) == "readme.md" {
			// Convert to relative path
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}

func main() {
	// Get the project root directory
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	// Find all README.md files
	readmeFiles, err := findReadmeFiles(root)
	if err != nil {
		fmt.Printf("Error finding README.md files: %v\n", err)
		os.Exit(1)
	}

	// Validate each README.md file
	valid := true
	for _, file := range readmeFiles {
		fileValid, errors := validateReadme(file)
		if !fileValid {
			valid = false
			fmt.Printf("Validation failed for %s:\n", file)
			for _, err := range errors {
				fmt.Printf("  - %s\n", err)
			}
			fmt.Println()
		}
	}

	if !valid {
		fmt.Println("README validation failed. Please fix the issues above.")
		os.Exit(1)
	}

	fmt.Println("All README.md files are valid!")
}
