// Copyright (c) 2025 A Bit of Help, Inc.

package security

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	// Common validation patterns
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	uuidRegex     = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	alphaNumRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

// InputValidator provides validation functionality for input data
type InputValidator struct {
	validator *validator.Validate
	logger    *logging.ContextLogger
}

// NewInputValidator creates a new InputValidator
func NewInputValidator(logger *logging.ContextLogger) *InputValidator {
	if logger == nil {
		panic("logger cannot be nil")
	}

	v := validator.New()

	// Register custom validation functions
	v.RegisterValidation("email", validateEmail)
	v.RegisterValidation("uuid", validateUUID)
	v.RegisterValidation("alphanum", validateAlphaNum)
	v.RegisterValidation("nohtml", validateNoHTML)

	return &InputValidator{
		validator: v,
		logger:    logger,
	}
}

// Validate validates the given struct using struct tags
func (v *InputValidator) Validate(ctx context.Context, data interface{}) error {
	v.logger.Debug(ctx, "Validating input data", zap.String("type", reflect.TypeOf(data).String()))

	if err := v.validator.Struct(data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// Convert validation errors to a more user-friendly format
			errorMessages := make([]string, 0, len(validationErrors))
			for _, e := range validationErrors {
				errorMessages = append(errorMessages, fmt.Sprintf(
					"Field '%s' failed validation: %s",
					e.Field(),
					e.Tag(),
				))
			}
			errorMsg := strings.Join(errorMessages, "; ")
			v.logger.Warn(ctx, "Input validation failed", zap.String("errors", errorMsg))
			return errors.NewValidationError(errorMsg, "input", err)
		}
		v.logger.Error(ctx, "Unexpected validation error", zap.Error(err))
		return errors.NewValidationError("validation failed", "input", err)
	}

	v.logger.Debug(ctx, "Input validation successful")
	return nil
}

// ValidateField validates a single field with the given tag
func (v *InputValidator) ValidateField(ctx context.Context, fieldName string, value interface{}, tag string) error {
	v.logger.Debug(ctx, "Validating field",
		zap.String("field", fieldName),
		zap.String("tag", tag))

	err := v.validator.Var(value, tag)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make([]string, 0, len(validationErrors))
			for _, e := range validationErrors {
				errorMessages = append(errorMessages, fmt.Sprintf(
					"Field '%s' failed validation: %s",
					fieldName,
					e.Tag(),
				))
			}
			errorMsg := strings.Join(errorMessages, "; ")
			v.logger.Warn(ctx, "Field validation failed",
				zap.String("field", fieldName),
				zap.String("errors", errorMsg))
			return errors.NewValidationError(errorMsg, fieldName, err)
		}
		v.logger.Error(ctx, "Unexpected field validation error",
			zap.String("field", fieldName),
			zap.Error(err))
		return errors.NewValidationError("validation failed", fieldName, err)
	}

	v.logger.Debug(ctx, "Field validation successful", zap.String("field", fieldName))
	return nil
}

// Custom validation functions

func validateEmail(fl validator.FieldLevel) bool {
	return emailRegex.MatchString(fl.Field().String())
}

func validateUUID(fl validator.FieldLevel) bool {
	return uuidRegex.MatchString(fl.Field().String())
}

func validateAlphaNum(fl validator.FieldLevel) bool {
	return alphaNumRegex.MatchString(fl.Field().String())
}

func validateNoHTML(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return !strings.Contains(value, "<") && !strings.Contains(value, ">")
}

// SanitizeString sanitizes a string by removing potentially dangerous characters
func SanitizeString(input string) string {
	// Remove HTML tags
	noHTML := strings.ReplaceAll(input, "<", "&lt;")
	noHTML = strings.ReplaceAll(noHTML, ">", "&gt;")

	// Remove other potentially dangerous characters
	noHTML = strings.ReplaceAll(noHTML, "\"", "&quot;")
	noHTML = strings.ReplaceAll(noHTML, "'", "&#39;")
	noHTML = strings.ReplaceAll(noHTML, "`", "&#96;")

	return noHTML
}