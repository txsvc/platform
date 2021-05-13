package auth

import (
	"context"
	"errors"

	"github.com/txsvc/platform/v2/pkg/account"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

const (
	MsgAuthenticationNotFound = "account '%s' not found"

	// FIXME
	DefaultEndpoint = "https://podops.dev"

	// AuthTypeSimpleToken constant token
	AuthTypeSimpleToken = "token"
	// AuthTypeJWT constant jwt
	AuthTypeJWT = "jwt"
	// AuthTypeSlack constant slack
	AuthTypeSlack = "slack"

	// DefaultAuthenticationExpiration in minutes. Used when sending an
	// account challenge or the temporary token.
	DefaultAuthenticationExpiration = 10
	// DefaultAuthorizationExpiration in days
	DefaultAuthorizationExpiration = 90

	// other defaults
	DefaultTokenType = "user" // other possibilities: app, bot, ...

	// default scopes
	ScopeAPIAdmin = "api:admin"
	// FIXME DefaultScope  = "api:read,api:write"
	DefaultScope = "production:read,production:write,production:build,resource:read,resource:write"
)

type (
	// Authorization represents a user, app or bot and its permissions
	Authorization struct {
		ClientID  string `json:"client_id" binding:"required"` // UNIQUE
		Realm     string `json:"realm"`
		Token     string `json:"token" binding:"required"`
		TokenType string `json:"token_type" binding:"required"` // user,app,bot
		UserID    string `json:"user_id"`                       // depends on TokenType. UserID could equal ClientID or BotUSerID in Slack
		Scope     string `json:"scope"`                         // a comma separated list of scopes, see below
		Expires   int64  `json:"expires"`                       // 0 = never
		// internal
		Revoked bool  `json:"-"`
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// AuthorizationRequest represents a login/authorization request from a user, app, or bot
	AuthorizationRequest struct {
		Realm    string `json:"realm" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
		ClientID string `json:"client_id"`
		Token    string `json:"token"`
		Scope    string `json:"scope"`
	}

	AuthorizationProvider interface {
		// CreateAuthorization creates a new Authorization that is application/service specific
		CreateAuthorization(*account.Account, *AuthorizationRequest) *Authorization
		// Send an account challenge to confirm the account
		SendAccountChallenge(context.Context, *account.Account) error
		// Send the new token
		SendAuthToken(context.Context, *account.Account) error
		// Scope returns the scope
		Scope() string
		// Endpoint returns the api endpoint url
		Endpoint() string
		// AuthenticationExpiration
		AuthenticationExpiration() int
		// AuthorizationExpiration
		AuthorizationExpiration() int
	}
)

var (
	// ErrNotAuthorized indicates that the API caller is not authorized
	ErrNotAuthorized = errors.New("not authorized")
	// ErrNoToken indicates that no bearer token was provided
	ErrNoToken = errors.New("no token provided")
	// ErrInvalidRoute indicates that the route and/or its parameters are not valid
	ErrInvalidRoute = errors.New("invalid route")

	// the default authentication provider instance
	authProvider *AuthorizationProviderImpl
)

func init() {
	authProvider = NewAuthorizationProvider("platform.null.auth").(*AuthorizationProviderImpl)
}

func NewAuthorizationProvider(ID string) interface{} {
	return &AuthorizationProviderImpl{}
}

func (auth *AuthorizationProviderImpl) Close() error {
	return nil
}

func (auth *AuthorizationProviderImpl) CreateAuthorization(account *account.Account, req *AuthorizationRequest) *Authorization {
	now := timestamp.Now()
	scope := DefaultScope
	if req.Scope != "" {
		scope = req.Scope
	}

	a := Authorization{
		ClientID:  account.ClientID,
		Realm:     req.Realm,
		Token:     CreateSimpleToken(),
		TokenType: DefaultTokenType,
		UserID:    req.UserID,
		Scope:     scope,
		Revoked:   false,
		Expires:   now + (DefaultAuthorizationExpiration * 86400),
		Created:   now,
		Updated:   now,
	}
	return &a
}

// SendAccountChallenge sends a notification to the user promting to confirm the account
func (auth *AuthorizationProviderImpl) SendAccountChallenge(ctx context.Context, account *account.Account) error {
	return nil
}

// SendAuthToken sends a notification to the user with the current authentication token
func (auth *AuthorizationProviderImpl) SendAuthToken(ctx context.Context, account *account.Account) error {
	return nil
}

func (auth *AuthorizationProviderImpl) Scope() string {
	return DefaultScope
}

func (auth *AuthorizationProviderImpl) Endpoint() string {
	return DefaultEndpoint
}

func (auth *AuthorizationProviderImpl) AuthenticationExpiration() int {
	return DefaultAuthenticationExpiration
}

func (auth *AuthorizationProviderImpl) AuthorizationExpiration() int {
	return DefaultAuthorizationExpiration
}
