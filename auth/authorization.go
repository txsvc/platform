package auth

import (
	"context"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/labstack/echo/v4"

	"github.com/txsvc/platform/v2/pkg/account"
	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	datastoreAuthorizations string = "AUTHORIZATIONS"
)

// IsValid verifies that the Authorization is still valid, i.e. is not expired and not revoked.
func (a *Authorization) IsValid() bool {
	if a.Revoked {
		return false
	}
	if a.Expires < timestamp.Now() {
		return false
	}
	return true
}

// HasAdminScope checks if the authorization includes scope 'api:admin'
func (a *Authorization) HasAdminScope() bool {
	return strings.Contains(a.Scope, ScopeAPIAdmin)
}

// CheckAuthorization relies on the presence of a bearer token and validates the
// matching authorization against a list of requested scopes. If everything checks
// out, the function returns the authorization or an error otherwise.
func CheckAuthorization(ctx context.Context, c echo.Context, scope string) (*Authorization, error) {
	token, err := GetBearerToken(c.Request())
	if err != nil {
		return nil, err
	}

	auth, err := FindAuthorizationByToken(ctx, token)
	if err != nil || auth == nil || !auth.IsValid() {
		return nil, ErrNotAuthorized
	}

	acc, err := account.FindAccountByUserID(ctx, auth.Realm, auth.UserID)
	if err != nil {
		return nil, err
	}

	if acc.Status != account.AccountActive {
		return nil, ErrNotAuthorized // not logged-in
	}

	if !hasScope(auth.Scope, scope) {
		return nil, ErrNotAuthorized
	}

	return auth, nil
}

func NewAuthorization(account *account.Account, req *AuthorizationRequest, expires int) *Authorization {
	now := timestamp.Now()

	a := Authorization{
		ClientID:  account.ClientID,
		Realm:     req.Realm,
		Token:     CreateSimpleToken(),
		TokenType: DefaultTokenType,
		UserID:    req.UserID,
		Scope:     req.Scope,
		Revoked:   false,
		Expires:   now + int64(expires*86400),
		Created:   now,
		Updated:   now,
	}
	return &a
}

// CreateAuthorization creates all data needed for the auth fu
func CreateAuthorization(ctx context.Context, auth *Authorization) error {
	k := authorizationKey(auth.Realm, auth.ClientID)

	// FIXME add a cache ?

	// we simply overwrite the existing authorization. If this is no desired, use GetAuthorization first,
	// update the Authorization and then write it back.
	_, err := ds.DataStore().Put(ctx, k, auth)
	return err
}

// UpdateAuthorization updates all data needed for the auth fu
func UpdateAuthorization(ctx context.Context, auth *Authorization) error {
	k := authorizationKey(auth.Realm, auth.ClientID)
	// FIXME add a cache ?

	// we simply overwrite the existing authorization. If this is no desired, use GetAuthorization first,
	// update the Authorization and then write it back.
	_, err := ds.DataStore().Put(ctx, k, auth)
	return err
}

// LookupAuthorization looks for an authorization
func LookupAuthorization(ctx context.Context, realm, clientID string) (*Authorization, error) {
	var auth Authorization
	k := authorizationKey(realm, clientID)

	// FIXME add a cache ?

	if err := ds.DataStore().Get(ctx, k, &auth); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // Not finding one is not an error!
		}
		return nil, err
	}
	return &auth, nil
}

// FindAuthorizationByToken looks for an authorization by the token
func FindAuthorizationByToken(ctx context.Context, token string) (*Authorization, error) {
	var auth []*Authorization

	// FIXME add a cache ?

	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreAuthorizations).Filter("Token =", token), &auth); err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, nil
	}
	return auth[0], nil
}

// ExchangeToken confirms the temporary auth token and creates the permanent one
func ExchangeToken(ctx context.Context, req *AuthorizationRequest, expires int, loginFrom string) (*Authorization, int, error) {
	var auth *Authorization

	acc, err := account.FindAccountByUserID(ctx, req.Realm, req.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if acc == nil {
		return nil, http.StatusNotFound, nil
	}
	now := timestamp.Now()
	if acc.Expires < now || acc.Token != req.Token {
		return nil, http.StatusUnauthorized, nil
	}

	// all OK, create or update the authorization
	auth, err = LookupAuthorization(ctx, acc.Realm, acc.ClientID)
	if err != nil {
		return nil, http.StatusInternalServerError, err // FIXME maybe use a different code here
	}
	if auth == nil {
		if req.Scope == "" {
			return nil, http.StatusBadRequest, ErrNoScope
		}
		auth = NewAuthorization(acc, req, expires)
	}
	auth.Token = CreateSimpleToken()
	auth.Expires = now + (int64(expires) * 86400)
	auth.Updated = now

	err = CreateAuthorization(ctx, auth)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// update the account
	acc.Status = account.AccountActive
	acc.LastLogin = now
	acc.LoginCount = acc.LoginCount + 1
	acc.LoginFrom = loginFrom
	acc.Token = ""
	acc.Expires = 0 // never

	err = account.UpdateAccount(ctx, acc)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return auth, http.StatusOK, nil
}

// authorizationKey creates a datastore key for a workspace authorization based on the team_id.
func authorizationKey(realm, client string) *datastore.Key {
	return datastore.NameKey(datastoreAuthorizations, namedKey(realm, client), nil)
}

func namedKey(part1, part2 string) string {
	return part1 + "." + part2
}

func hasScope(scopes, scope string) bool {
	if scopes == "" || scope == "" {
		return false // empty inputs should never evalute to true
	}

	// FIXME this is a VERY naiv implementation
	return strings.Contains(scopes, scope)
}
