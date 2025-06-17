// Package auth provides authentication and authorization functionality for the family service.
// It includes JWT token handling, OIDC integration, and middleware for securing API endpoints.
// This package is part of the infrastructure layer and implements the authorization interfaces
// defined in the ports layer.
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// JWTConfig holds the configuration for JWT token handling.
// It contains the necessary parameters for creating and validating JWT tokens.
type JWTConfig struct {
	// SecretKey is the key used to sign and verify JWT tokens
	SecretKey string

	// TokenDuration is the validity period for generated tokens
	TokenDuration time.Duration

	// Issuer identifies the entity that issued the token
	Issuer string
}

// JWTService handles JWT token operations including generation and validation.
// It implements token-based authentication for the application.
type JWTService struct {
	// config contains the JWT configuration parameters
	config JWTConfig

	// logger is used for logging token operations and errors
	logger *zap.Logger
}

// NewJWTService creates a new JWT service with the provided configuration and logger.
// Parameters:
//   - config: The configuration for JWT token handling
//   - logger: The logger for recording operations and errors
//
// Returns:
//   - *JWTService: A new instance of the JWT service
func NewJWTService(config JWTConfig, logger *zap.Logger) *JWTService {
	return &JWTService{
		config: config,
		logger: logger,
	}
}

// Claims represents the JWT claims contained in a token.
// It extends the standard JWT registered claims with custom application-specific claims.
type Claims struct {
	// UserID is the unique identifier of the user (stored in the 'sub' claim)
	UserID string `json:"sub"`

	// Roles contains the user's assigned roles for authorization
	Roles []string `json:"roles"`

	// RegisteredClaims contains the standard JWT claims like expiration time
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user with the specified roles.
// It creates a token with claims including user ID, roles, and standard JWT claims
// like expiration time, issued at time, and issuer.
// Parameters:
//   - userID: The unique identifier of the user
//   - roles: The roles assigned to the user for authorization purposes
//
// Returns:
//   - string: The signed JWT token string if successful
//   - error: An error if token generation fails
func (s *JWTService) GenerateToken(userID string, roles []string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.TokenDuration)

	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims if valid.
// It verifies the token signature, checks if the token is expired,
// and extracts the claims from the token.
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns:
//   - *Claims: The claims from the token if validation is successful
//   - error: An error if token validation fails for any reason
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		s.logger.Debug("Failed to parse token", zap.Error(err))
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		s.logger.Debug("Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		s.logger.Debug("Failed to extract claims from token")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
