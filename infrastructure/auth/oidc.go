package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/knadh/koanf/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// getOIDCTimeout returns the timeout for OIDC operations from configuration
func getOIDCTimeout(k *koanf.Koanf) time.Duration {
	return time.Duration(k.Int("auth.oidc_timeout")) * time.Second
}

// OIDCConfig holds the configuration for OIDC
type OIDCConfig struct {
	IssuerURL     string
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	Scopes        []string
	AdminRoleName string
}

// OIDCService handles OIDC operations
type OIDCService struct {
	config      OIDCConfig
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	oauthConfig *oauth2.Config
	logger      *zap.Logger
	tracer      trace.Tracer
	koanf       *koanf.Koanf
}

// NewOIDCService creates a new OIDC service
func NewOIDCService(ctx context.Context, config OIDCConfig, logger *zap.Logger, k *koanf.Koanf) (*OIDCService, error) {
	// Validate context
	if ctx == nil {
		ctx = context.Background()
		logger.Warn("Nil context provided to NewOIDCService, using background context")
	}

	// Create a context with timeout for OIDC provider discovery
	discoveryCtx, cancel := context.WithTimeout(ctx, getOIDCTimeout(k))
	defer cancel()

	// Discover OIDC provider
	provider, err := oidc.NewProvider(discoveryCtx, config.IssuerURL)
	if err != nil {
		logger.Error("Failed to discover OIDC provider", zap.Error(err), zap.String("issuer_url", config.IssuerURL))
		return nil, fmt.Errorf("failed to discover OIDC provider: %w", err)
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	// Create OAuth2 config
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       append([]string{oidc.ScopeOpenID, "profile", "email"}, config.Scopes...),
	}

	return &OIDCService{
		config:      config,
		provider:    provider,
		verifier:    verifier,
		oauthConfig: oauthConfig,
		logger:      logger,
		tracer:      otel.Tracer("infrastructure.auth.oidc"),
		koanf:       k,
	}, nil
}

// ValidateToken validates an OIDC token and returns the claims
func (s *OIDCService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	ctx, span := s.tracer.Start(ctx, "OIDCService.ValidateToken")
	defer span.End()

	span.SetAttributes(attribute.String("token.length", fmt.Sprintf("%d", len(tokenString))))

	// Parse and verify the ID token
	idToken, err := s.verifier.Verify(ctx, tokenString)
	if err != nil {
		s.logger.Debug("Failed to verify ID token", zap.Error(err))
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims from the ID token
	var claims struct {
		Subject string   `json:"sub"`
		Email   string   `json:"email"`
		Name    string   `json:"name"`
		Roles   []string `json:"roles"`
	}
	if err := idToken.Claims(&claims); err != nil {
		s.logger.Debug("Failed to extract claims from ID token", zap.Error(err))
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	// Create JWT claims
	jwtClaims := &Claims{
		UserID: claims.Subject,
		Roles:  claims.Roles,
	}

	return jwtClaims, nil
}

// ExtractTokenFromRequest extracts the token from the Authorization header
func (s *OIDCService) ExtractTokenFromRequest(r *http.Request) (string, error) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header provided")
	}

	// Check if the header has the Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}

// IsAdmin checks if the user has the admin role
func (s *OIDCService) IsAdmin(roles []string) bool {
	for _, role := range roles {
		if role == s.config.AdminRoleName {
			return true
		}
	}
	return false
}

// GetAuthURL returns the URL for the OAuth2 authorization endpoint
func (s *OIDCService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// Exchange exchanges an authorization code for a token
func (s *OIDCService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	ctx, span := s.tracer.Start(ctx, "OIDCService.Exchange")
	defer span.End()

	// Create a context with timeout for token exchange
	exchangeCtx, cancel := context.WithTimeout(ctx, getOIDCTimeout(s.koanf))
	defer cancel()

	return s.oauthConfig.Exchange(exchangeCtx, code)
}

// GetUserInfo gets the user info from the OIDC provider
func (s *OIDCService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error) {
	ctx, span := s.tracer.Start(ctx, "OIDCService.GetUserInfo")
	defer span.End()

	// Create a context with timeout for user info request
	userInfoCtx, cancel := context.WithTimeout(ctx, getOIDCTimeout(s.koanf))
	defer cancel()

	return s.provider.UserInfo(userInfoCtx, oauth2.StaticTokenSource(token))
}

// InitOIDCServiceFromConfig initializes the OIDC service from configuration parameters
func InitOIDCServiceFromConfig(ctx context.Context, issuerURL, clientID, clientSecret, redirectURL string, scopes []string, adminRoleName string, logger *zap.Logger, k *koanf.Koanf) (*OIDCService, error) {
	// Validate context
	if ctx == nil {
		ctx = context.Background()
		logger.Warn("Nil context provided to InitOIDCServiceFromConfig, using background context")
	}

	// Validate required configuration
	if issuerURL == "" {
		logger.Warn("OIDC issuer URL not provided, OIDC authentication will be disabled")
		return nil, nil
	}

	if clientID == "" {
		logger.Warn("OIDC client ID not provided, OIDC authentication will be disabled")
		return nil, nil
	}

	// Create OIDC configuration
	config := OIDCConfig{
		IssuerURL:     issuerURL,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		RedirectURL:   redirectURL,
		Scopes:        scopes,
		AdminRoleName: adminRoleName,
	}

	// Create OIDC service
	oidcService, err := NewOIDCService(ctx, config, logger, k)
	if err != nil {
		logger.Error("Failed to initialize OIDC service", zap.Error(err))
		return nil, err
	}

	logger.Info("OIDC service initialized successfully",
		zap.String("issuer_url", issuerURL),
		zap.String("client_id", clientID),
		zap.String("redirect_url", redirectURL),
	)

	return oidcService, nil
}
