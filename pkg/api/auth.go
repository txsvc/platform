package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform/v2"

	"github.com/txsvc/platform/v2/auth"
	"github.com/txsvc/platform/v2/pkg/account"
)

// LoginRequestEndpoint initiates the login process.
//
// It creates a new account if the user does not exist and sends
// confirmation request. Once the account is conformed, it will send the
// confirmation token that can be swapped for a real login token.
//
// POST /login
// status 201: new account, account confirmation sent
// status 204: existing account, email with auth token sent
// status 400: invalid request data
// status 403: only logged-out and confirmed users can proceed
func LoginRequestEndpoint(c echo.Context) error {
	var req *auth.AuthorizationRequest = new(auth.AuthorizationRequest)
	ctx := platform.NewHttpContext(c.Request())

	err := c.Bind(req)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if req.Realm == "" || req.UserID == "" {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	acc, err := account.FindAccountByUserID(ctx, req.Realm, req.UserID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// new account
	if acc == nil {
		// #1: create a new account
		acc, err = account.CreateAccount(ctx, req.Realm, req.UserID, platform.AuthorizationProvider().AuthenticationExpiration())
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, err)
		}
		// #2: send the confirmation link
		err = platform.AuthorizationProvider().AccountChallengeNotification(ctx, acc)
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, err)
		}
		// status 201: new account
		return c.NoContent(http.StatusCreated)
	}

	// existing account but check some stuff first ...
	if acc.Confirmed == 0 {
		// #1: update the expiration timestamp
		acc, err = account.ResetAccountChallenge(ctx, acc, platform.AuthorizationProvider().AuthenticationExpiration())
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, err)
		}
		// #2: send the account confirmation link
		err = platform.AuthorizationProvider().AccountChallengeNotification(ctx, acc)
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, err)
		}
		// status 201: new account
		return c.NoContent(http.StatusCreated)
	}
	if acc.Status != 0 {
		// status 403: only logged-out and confirmed users can proceed, do nothing otherwise
		return ErrorResponse(c, http.StatusForbidden, auth.ErrAlreadyAuthorized)
	}

	// create and send the auth token
	acc, err = account.ResetTemporaryToken(ctx, acc, platform.AuthorizationProvider().AuthenticationExpiration())
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	err = platform.AuthorizationProvider().ProvideAuthorizationToken(ctx, acc)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// status 204: existing account, email with token sent
	return c.NoContent(http.StatusNoContent)
}

// LogoutRequestEndpoint removes the session data and state
//
// POST /logout

func LogoutRequestEndpoint(c echo.Context) error {
	var req *auth.AuthorizationRequest = new(auth.AuthorizationRequest)
	ctx := platform.NewHttpContext(c.Request())

	err := c.Bind(req)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}
	if req.Realm == "" || req.UserID == "" {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	token, err := auth.GetBearerToken(c.Request())
	if err != nil {
		return auth.ErrNoToken
	}
	ath, err := auth.FindAuthorizationByToken(ctx, token)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}
	if ath.UserID != req.UserID || ath.Realm != req.Realm {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	// logout starts here
	status, err := auth.LogoutAccount(ctx, ath.Realm, ath.ClientID)
	if err != nil {
		return ErrorResponse(c, status, err)
	}
	return c.NoContent(status)
}

// LoginConfirmationEndpoint validates an email.
//
// GET /login/:token
// status 307: account is confirmed, redirect to podops.dev/confirmed
// status 400: the request could not be understood by the server due to malformed syntax
// status 401: token is wrong
// status 403: token is expired or has already been used
// status 404: token was not found
func LoginConfirmationEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	token := c.Param("token")
	if token == "" {
		return ErrorResponse(c, http.StatusBadRequest, auth.ErrInvalidRoute)
	}

	acc, status, err := auth.ConfirmLoginChallenge(ctx, token)
	if status != http.StatusNoContent {
		return ErrorResponse(c, status, err)
	}

	acc, err = account.ResetTemporaryToken(ctx, acc, platform.AuthorizationProvider().AuthenticationExpiration())
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	err = platform.AuthorizationProvider().ProvideAuthorizationToken(ctx, acc)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	// status 307: account is confirmed, email with auth token sent, redirect now
	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/confirmed", platform.AuthorizationProvider().Endpoint()))
}

// GetAuthorizationEndpoint exchanges a temporary confirmation token for a 'real' token.
//
// POST /auth
// status 200: success, the real token is in the response
// status 401: token is expired or has already been used, token and user_id do not match
// status 404: token was not found
func GetAuthorizationEndpoint(c echo.Context) error {
	var req *auth.AuthorizationRequest = new(auth.AuthorizationRequest)
	ctx := platform.NewHttpContext(c.Request())

	err := c.Bind(req)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	if req.Token == "" || req.Realm == "" || req.UserID == "" {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	// make sure we have a known default scope and no one sneaks something in
	req.Scope = platform.AuthorizationProvider().Scope()

	ath, status, err := auth.ExchangeToken(ctx, req, platform.AuthorizationProvider().AuthorizationExpiration(), c.Request().RemoteAddr)
	if status != http.StatusOK {
		return ErrorResponse(c, status, err)
	}

	req.Token = ath.Token
	req.ClientID = ath.ClientID

	return StandardResponse(c, status, req)
}