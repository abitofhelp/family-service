// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"net/mail"
	"strings"
)

// Email represents an email address value object
type Email string

// NewEmail creates a new Email with validation
func NewEmail(email string) (Email, error) {
	// Trim whitespace
	trimmedEmail := strings.TrimSpace(email)

	// Empty email is allowed (optional field)
	if trimmedEmail == "" {
		return "", nil
	}

	// Validate email format
	_, err := mail.ParseAddress(trimmedEmail)
	if err != nil {
		return "", errors.New("invalid email format")
	}

	return Email(trimmedEmail), nil
}

// String returns the string representation of the Email
func (e Email) String() string {
	return string(e)
}

// Equals checks if two Emails are equal
func (e Email) Equals(other Email) bool {
	return strings.EqualFold(string(e), string(other))
}

// IsEmpty checks if the Email is empty
func (e Email) IsEmpty() bool {
	return e == ""
}
