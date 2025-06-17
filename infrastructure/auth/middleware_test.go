package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockHandler is a mock HTTP handler for testing
type MockHandler struct {
	mock.Mock
	CalledWithContext context.Context
}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	m.CalledWithContext = r.Context()
}

func TestNewAuthMiddleware(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	jwtService := &JWTService{
		config: JWTConfig{
			SecretKey:     "test-secret-key",
			Issuer:        "test-issuer",
		},
		logger: logger,
	}

	// Execute
	middleware := NewAuthMiddleware(jwtService, logger)

	// Verify
	assert.NotNil(t, middleware)
	assert.Equal(t, jwtService, middleware.jwtService)
	assert.Equal(t, logger, middleware.logger)
	assert.NotNil(t, middleware.tracer)
}

func TestNewAuthMiddlewareWithOIDC(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	jwtService := &JWTService{
		config: JWTConfig{
			SecretKey:     "test-secret-key",
			Issuer:        "test-issuer",
		},
		logger: logger,
	}
	oidcService := &OIDCService{
		logger: logger,
	}

	// Execute
	middleware := NewAuthMiddlewareWithOIDC(jwtService, oidcService, logger)

	// Verify
	assert.NotNil(t, middleware)
	assert.Equal(t, jwtService, middleware.jwtService)
	assert.Equal(t, oidcService, middleware.oidcService)
	assert.Equal(t, logger, middleware.logger)
	assert.NotNil(t, middleware.tracer)
}

func TestMiddleware_NoAuthHeader(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	jwtService := &JWTService{
		config: JWTConfig{
			SecretKey:     "test-secret-key",
			Issuer:        "test-issuer",
		},
		logger: logger,
	}
	middleware := NewAuthMiddleware(jwtService, logger)

	mockHandler := new(MockHandler)
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	// Execute
	handler := middleware.Middleware(mockHandler)
	handler.ServeHTTP(res, req)

	// Verify
	mockHandler.AssertCalled(t, "ServeHTTP", res, mock.Anything)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	jwtService := &JWTService{
		config: JWTConfig{
			SecretKey:     "test-secret-key",
			Issuer:        "test-issuer",
		},
		logger: logger,
	}
	middleware := NewAuthMiddleware(jwtService, logger)

	mockHandler := new(MockHandler)

	testCases := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "No Bearer prefix",
			authHeader: "token123",
		},
		{
			name:       "Wrong format",
			authHeader: "Bearer",
		},
		{
			name:       "Multiple spaces",
			authHeader: "Bearer  token123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tc.authHeader)
			res := httptest.NewRecorder()

			// Execute
			handler := middleware.Middleware(mockHandler)
			handler.ServeHTTP(res, req)

			// Verify
			mockHandler.AssertNotCalled(t, "ServeHTTP")
			assert.Equal(t, http.StatusUnauthorized, res.Code)
		})
	}
}

func TestMiddleware_ValidToken(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: time.Hour * 24, // Ensure token doesn't expire
		Issuer:        "test-issuer",
	}
	jwtService := NewJWTService(config, logger)
	middleware := NewAuthMiddleware(jwtService, logger)

	mockHandler := new(MockHandler)
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).Return()

	// Generate a valid token
	userID := "test-user-id"
	roles := []string{"admin", "user"}
	token, err := jwtService.GenerateToken(userID, roles)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	// Execute
	handler := middleware.Middleware(mockHandler)
	handler.ServeHTTP(res, req)

	// Verify
	mockHandler.AssertCalled(t, "ServeHTTP", res, mock.Anything)
	assert.Equal(t, http.StatusOK, res.Code)

	// Verify context values
	ctx := mockHandler.CalledWithContext
	assert.NotNil(t, ctx)

	// Get user ID from context using the helper function
	authService := NewAuthorizationService(logger)
	ctxUserID, err := authService.GetUserID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, userID, ctxUserID)

	// Get user roles from context using the helper function
	ctxRoles, err := authService.GetUserRoles(ctx)
	assert.NoError(t, err)
	assert.Equal(t, roles, ctxRoles)
}

func TestMiddleware_InvalidToken(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
	}
	jwtService := NewJWTService(config, logger)
	middleware := NewAuthMiddleware(jwtService, logger)

	mockHandler := new(MockHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	res := httptest.NewRecorder()

	// Execute
	handler := middleware.Middleware(mockHandler)
	handler.ServeHTTP(res, req)

	// Verify
	mockHandler.AssertNotCalled(t, "ServeHTTP")
	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
