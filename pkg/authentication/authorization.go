package authentication

import (
	"context"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	mcache "github.com/OrlovEvgeny/go-mcache"
	"github.com/labstack/echo/v4"

	"github.com/txsvc/platform/v2/pkg/account"
	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/loader"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	datastoreAuthorizations string = "AUTHORIZATIONS"
)

var (
	// secondary cache, maps token to authorization
	authCache = mcache.New()
)

func (ath *Authorization) Equal(a *Authorization) bool {
	if a == nil {
		return false
	}
	return ath.Token == a.Token && ath.Realm == a.Realm && ath.ClientID == a.ClientID && ath.UserID == a.UserID
}

// IsValid verifies that the Authorization is still valid, i.e. is not expired and not revoked.
func (ath *Authorization) IsValid() bool {
	if ath.Revoked {
		return false
	}
	if ath.Expires < timestamp.Now() {
		return false
	}
	return true
}

// HasAdminScope checks if the authorization includes scope 'api:admin'
func (ath *Authorization) HasAdminScope() bool {
	return strings.Contains(ath.Scope, ScopeAPIAdmin)
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

func NewAuthorization(req *AuthorizationRequest, expires int) *Authorization {
	now := timestamp.Now()

	a := Authorization{
		ClientID:  req.ClientID,
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

// UpdateAuthorization updates all data needed for the auth fu
func UpdateAuthorization(ctx context.Context, auth *Authorization) error {
	k := nativeKey(auth.Key())

	// remove from the cache
	authCache.Remove(auth.Token)

	// we simply overwrite the existing authorization. If this is no desired, use GetAuthorization first,
	// update the Authorization and then write it back.
	if _, err := ds.DataStore().Put(ctx, k, auth); err != nil {
		return err
	}

	return nil
}

// LookupAuthorization looks for an authorization
func LookupAuthorization(ctx context.Context, realm, clientID string) (*Authorization, error) {
	var auth Authorization

	k := nativeKey(namedKey(realm, clientID))

	if err := ds.DataStore().Get(ctx, k, &auth); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // Not finding one is not an error!
		}
		return nil, err
	}

	return &auth, nil
}

func DeleteAuthorization(ctx context.Context, realm, clientID string) (*Authorization, error) {
	auth, err := LookupAuthorization(ctx, realm, clientID)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, ErrNoSuchEntity
	}

	k := nativeKey(namedKey(realm, clientID))
	if err := ds.DataStore().Delete(ctx, k); err != nil {
		return nil, err
	}

	// remove from the cache
	authCache.Remove(auth.Token)

	return auth, nil
}

// FindAuthorizationByToken looks for an authorization by the token
func FindAuthorizationByToken(ctx context.Context, token string) (*Authorization, error) {
	if token == "" {
		return nil, ErrNoSuchEntity
	}
	if a, ok := authCache.Get(token); ok {
		return a.(*Authorization), nil
	}

	var auth []*Authorization

	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreAuthorizations).Filter("Token =", token), &auth); err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, nil
	}

	a := auth[0]

	// add the authorization to the cache
	authCache.Set(a.Token, a, loader.DefaultTTL)

	return a, nil
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
		req.ClientID = acc.ClientID
		auth = NewAuthorization(req, expires)
	}
	auth.Token = CreateSimpleToken()
	auth.Revoked = false
	auth.Expires = now + (int64(expires) * 86400)
	auth.Updated = now

	err = UpdateAuthorization(ctx, auth)
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

//
// keys, cache and dataloader
//

func (ath *Authorization) Key() string {
	return namedKey(ath.Realm, ath.ClientID)
}

func nativeKey(key string) *datastore.Key {
	return datastore.NameKey(datastoreAuthorizations, key, nil)
}

func namedKey(part1, part2 string) string {
	return part1 + "." + part2
}

func hasScope(scopes, scope string) bool {
	// FIXME this is a VERY simple implementation
	if scopes == "" || scope == "" {
		return false // empty inputs should never evalute to true
	}

	// FIXME this is a VERY naiv implementation
	return strings.Contains(scopes, scope)
}
