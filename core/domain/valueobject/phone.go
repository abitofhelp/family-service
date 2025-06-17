package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

// Phone represents a phone number value object
type Phone string

var phoneRegex = regexp.MustCompile(`^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`)

// NewPhone creates a new Phone with validation
func NewPhone(phone string) (Phone, error) {
	// Trim whitespace
	trimmedPhone := strings.TrimSpace(phone)
	
	// Empty phone is allowed (optional field)
	if trimmedPhone == "" {
		return "", nil
	}
	
	// Validate phone format with a simple regex
	// This is a basic validation and might need to be enhanced for specific requirements
	if !phoneRegex.MatchString(trimmedPhone) {
		return "", errors.New("invalid phone number format")
	}
	
	return Phone(trimmedPhone), nil
}

// String returns the string representation of the Phone
func (p Phone) String() string {
	return string(p)
}

// Equals checks if two Phones are equal
func (p Phone) Equals(other Phone) bool {
	// Remove all non-digit characters for comparison
	cleanThis := regexp.MustCompile(`\D`).ReplaceAllString(string(p), "")
	cleanOther := regexp.MustCompile(`\D`).ReplaceAllString(string(other), "")
	
	return cleanThis == cleanOther
}

// IsEmpty checks if the Phone is empty
func (p Phone) IsEmpty() bool {
	return p == ""
}