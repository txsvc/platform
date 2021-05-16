package authentication

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/txsvc/platform/v2/pkg/account"
	"github.com/txsvc/platform/v2/pkg/id"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

const (
	// AuthTypeBearerToken constant token
	AuthTypeBearerToken = "token"
	// AuthTypeJWT constant jwt
	AuthTypeJWT = "jwt"
	// AuthTypeSlack constant slack
	AuthTypeSlack = "slack"

	// other defaults
	UserTokenType    = "user"
	AppTokenType     = "app"
	APITokenType     = "api"
	BotTokenType     = "bot"
	DefaultTokenType = UserTokenType

	// default scopes
	DefaultScope  = "api:read,api:write"
	ScopeAPIAdmin = "api:admin"

	// DefaultAuthenticationExpiration in minutes. Used when sending an account challenge or the temporary token.
	DefaultAuthenticationExpiration = 10
	// DefaultAuthorizationExpiration in days
	DefaultAuthorizationExpiration = 90

	// DefaultEndpoint is used to build the urls in the notifications
	DefaultEndpoint = "http://localhost:8080"

	// error messages
	MsgAuthenticationNotFound = "account '%s' not found"
)

type (
	// Authorization represents a user, app or bot and its permissions
	Authorization struct {
		ClientID  string `json:"client_id" binding:"required"` // UNIQUE
		Realm     string `json:"realm"`
		Token     string `json:"token" binding:"required"`
		TokenType string `json:"token_type" binding:"required"` // e.g. user,app,api,bot
		UserID    string `json:"user_id"`                       // depends on TokenType. UserID could equal ClientID or BotUserID in Slack
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
)

var (
	// ErrNotAuthorized indicates that the API caller is not authorized
	ErrNotAuthorized     = errors.New("not authorized")
	ErrAlreadyAuthorized = errors.New("already authorized")

	// ErrNoSuchEntity indicates that the authorization does not exist
	ErrNoSuchEntity = errors.New("entity does not exist")

	// ErrNoToken indicates that no bearer token was provided
	ErrNoToken = errors.New("no token provided")
	// ErrNoScope indicates that no scope was provided
	ErrNoScope = errors.New("no scope provided")
	// ErrInvalidRoute indicates that the route and/or its parameters are not valid
	ErrInvalidRoute = errors.New("invalid route")
)

func CreateSimpleToken() string {
	token, _ := id.UUID()
	return token
}

// GetBearerToken extracts the bearer token
func GetBearerToken(r *http.Request) (string, error) {

	// FIXME optimize this !!

	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", ErrNoToken
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		return "", ErrNoToken
	}
	if parts[0] == "Bearer" {
		return parts[1], nil
	}

	return "", ErrNoToken
}

// GetClientID extracts the ClientID from the token
func GetClientID(ctx context.Context, r *http.Request) (string, error) {
	token, err := GetBearerToken(r)
	if err != nil {
		return "", err
	}

	// FIXME optimize this, e.g. implement caching

	auth, err := FindAuthorizationByToken(ctx, token)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", ErrNotAuthorized
	}

	return auth.ClientID, nil
}

// ConfirmLoginChallenge confirms the account
func ConfirmLoginChallenge(ctx context.Context, token string) (*account.Account, int, error) {
	if token == "" {
		return nil, http.StatusUnauthorized, ErrNoToken
	}
	acc, err := account.FindAccountByToken(ctx, token)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if acc == nil {
		return nil, http.StatusUnauthorized, nil
	}
	now := timestamp.Now()
	if acc.Expires < now {
		return acc, http.StatusForbidden, nil
	}

	acc.Confirmed = now
	acc.Expires = 0
	acc.Status = account.AccountLoggedOut
	acc.Token = ""

	err = account.UpdateAccount(ctx, acc)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return acc, http.StatusNoContent, nil
}

func LogoutAccount(ctx context.Context, realm, clientID string) (int, error) {
	// find the account
	acc, err := account.LookupAccount(ctx, realm, clientID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if acc == nil {
		return http.StatusBadRequest, ErrNoSuchEntity
	}

	if acc.Status < 0 {
		return http.StatusForbidden, nil // account is blocked or deactivated etc ...
	}

	acc.Status = account.AccountLoggedOut
	if err := account.UpdateAccount(ctx, acc); err != nil {
		return http.StatusInternalServerError, err
	}

	// find the matching authorization and revoke it
	auth, err := LookupAuthorization(ctx, acc.Realm, acc.ClientID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if auth != nil {
		auth.Revoked = true
		err = UpdateAuthorization(ctx, auth)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusNoContent, nil
}

func BlockAccount(ctx context.Context, realm, clientID string) error {
	acc, err := account.LookupAccount(ctx, realm, clientID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrNoSuchEntity
	}

	auth, err := LookupAuthorization(ctx, acc.Realm, acc.ClientID)
	if err != nil {
		return err
	}
	if auth != nil {
		auth.Revoked = true
		err = UpdateAuthorization(ctx, auth)
		if err != nil {
			return err
		}
	}

	acc.Status = account.AccountBlocked
	return account.UpdateAccount(ctx, acc)
}
