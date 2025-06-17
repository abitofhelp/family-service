package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewJWTService(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}

	// Execute
	service := NewJWTService(config, logger)

	// Verify
	assert.NotNil(t, service)
	assert.Equal(t, config, service.config)
	assert.Equal(t, logger, service.logger)
}

func TestGenerateToken(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}
	service := NewJWTService(config, logger)
	userID := "test-user-id"
	roles := []string{"admin", "user"}

	// Execute
	token, err := service.GenerateToken(userID, roles)

	// Verify
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse the token to verify its contents
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*Claims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, config.Issuer, claims.Issuer)
}

func TestValidateToken_Valid(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}
	service := NewJWTService(config, logger)
	userID := "test-user-id"
	roles := []string{"admin", "user"}

	// Generate a token
	token, err := service.GenerateToken(userID, roles)
	assert.NoError(t, err)

	// Execute
	claims, err := service.ValidateToken(token)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, config.Issuer, claims.Issuer)
}

func TestValidateToken_Invalid(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: time.Hour,
		Issuer:        "test-issuer",
	}
	service := NewJWTService(config, logger)

	// Test cases
	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Invalid format",
			token: "not.a.valid.token",
		},
		{
			name:  "Wrong signature",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute
			claims, err := service.ValidateToken(tc.token)

			// Verify
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: -time.Hour, // Negative duration to create an expired token
		Issuer:        "test-issuer",
	}
	service := NewJWTService(config, logger)
	userID := "test-user-id"
	roles := []string{"admin", "user"}

	// Generate an expired token
	token, err := service.GenerateToken(userID, roles)
	assert.NoError(t, err)

	// Execute
	claims, err := service.ValidateToken(token)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")
}