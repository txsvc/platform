package auth

import (
	"context"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/labstack/echo/v4"

	"github.com/txsvc/platform/v2/pkg/account"
	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/id"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	datastoreAuthorizations string = "AUTHORIZATIONS"
)

type (
	AuthorizationProviderImpl struct {
	}
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

func CreateSimpleToken() string {
	token, _ := id.UUID()
	return token
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
