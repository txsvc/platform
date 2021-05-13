package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/txsvc/platform/v2/pkg/account"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

func LogoutAccount(ctx context.Context, realm, clientID string) (int, error) {
	// find the account
	acc, err := account.LookupAccount(ctx, realm, clientID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if acc == nil {
		return http.StatusBadRequest, fmt.Errorf(MsgAuthenticationNotFound, fmt.Sprintf("%s.%s", realm, clientID))
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
		return fmt.Errorf(MsgAuthenticationNotFound, fmt.Sprintf("%s.%s", realm, clientID))
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
	acc.Status = account.AccountLoggedOut
	acc.Token = ""

	err = account.UpdateAccount(ctx, acc)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return acc, http.StatusNoContent, nil
}

// exchangeToken confirms the temporary auth token and creates the permanent one
func exchangeToken(ctx context.Context, req *AuthorizationRequest, loginFrom string) (*Authorization, int, error) {
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
			req.Scope = authProvider.Scope()
		}
		auth = authProvider.CreateAuthorization(acc, req)
	}
	auth.Token = CreateSimpleToken()
	auth.Expires = now + (int64(authProvider.AuthorizationExpiration()) * 86400)
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
