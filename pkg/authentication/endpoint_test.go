package authentication

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform/v2/pkg/account"
	ds "github.com/txsvc/platform/v2/pkg/datastore"
)

const (
	realm  = "podops"
	userID = "me@podops.dev"

	loginRequestRoute      = "/login"
	loginConfirmationRoute = "/login/:token"
	logoutRequestRoute     = "/logout"
	getAuthorizationRoute  = "/auth"
)

// Scenario 1: new account, login, account confirmation, token swap
func TestLoginScenario1(t *testing.T) {
	t.Cleanup(cleaner)
	cleaner()

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	acc := getAccount(t)
	loginStep2(t, acc.Token, http.StatusTemporaryRedirect, true) // confirm the new account, send auth token

	acc = getAccount(t)
	loginStep3(t, realm, userID, acc.ClientID, acc.Token, account.AccountActive, http.StatusOK, true) // exchange auth token for a permanent token

	verifyAccountAndAuth(t)

	auth, _ := LookupAuthorization(context.TODO(), acc.Realm, acc.ClientID)
	logoutStep(t, realm, userID, acc.ClientID, auth.Token, http.StatusNoContent, true)
}

// Scenario 2: new account, login, duplicate login request
func TestLoginScenario2(t *testing.T) {
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account
	account1 := getAccount(t)
	token1 := account1.Token

	loginStep1(t, http.StatusCreated) // existing account, request login again, create the account
	account2 := getAccount(t)
	token2 := account2.Token

	// requires a new token
	assert.NotEqual(t, token1, token2)

	loginStep2(t, account2.Token, http.StatusTemporaryRedirect, true) // confirm the new account, send auth token

	account3 := getAccount(t)
	loginStep3(t, realm, userID, account3.ClientID, account3.Token, account.AccountActive, http.StatusOK, true) // exchange auth token for a permanent token

	verifyAccountAndAuth(t)
}

// Scenario 3: new account, login, duplicate account confirmation
func TestLoginScenario3(t *testing.T) {
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	acc := getAccount(t)
	token := acc.Token

	loginStep2(t, token, http.StatusTemporaryRedirect, true) // confirm the new acc, send auth token
	loginStep2(t, token, http.StatusUnauthorized, true)      // confirm again

	acc = getAccount(t)
	loginStep3(t, realm, userID, acc.ClientID, acc.Token, account.AccountActive, http.StatusOK, true) // exchange auth token for a permanent token

	verifyAccountAndAuth(t)
}

// Scenario 4: new account, login, account confirmation, duplicate token swap
func TestLoginScenario4(t *testing.T) {
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	acc := getAccount(t)
	loginStep2(t, acc.Token, http.StatusTemporaryRedirect, true) // confirm the new account, send auth token

	acc = getAccount(t)
	token := acc.Token

	loginStep3(t, realm, userID, acc.ClientID, token, account.AccountActive, http.StatusOK, true) // exchange auth token for a permanent token

	loginStep3(t, realm, userID, acc.ClientID, token, account.AccountActive, http.StatusUnauthorized, true)

	verifyAccountAndAuth(t)
}

// Scenario 5: new account, login, invalid confirmation
func TestLoginScenario5(t *testing.T) {
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	loginStep2(t, "this_is_not_valid", http.StatusUnauthorized, false)
}

// Scenario 6: new account, login, account confirmation, various invalid token swaps
func TestLoginScenario6(t *testing.T) {
	t.Cleanup(cleaner)

	loginStep1(t, http.StatusCreated) // new account, request login, create the account

	acc := getAccount(t)
	loginStep2(t, acc.Token, http.StatusTemporaryRedirect, true) // confirm the new account, send auth token

	acc = getAccount(t)
	loginStep3(t, "", "", "", "", account.AccountLoggedOut, http.StatusBadRequest, false)
	loginStep3(t, "wrong_realm", "wrong_user", acc.ClientID, acc.Token, account.AccountLoggedOut, http.StatusNotFound, false)
	loginStep3(t, realm, userID, acc.ClientID, "wrong_auth_token", account.AccountLoggedOut, http.StatusUnauthorized, false)
}

// FIXME test account confirmation timeout

// FIXME test auth swap timeout

func loginStep1(t *testing.T, status int) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, loginRequestRoute, strings.NewReader(createAuthRequestJSON(realm, userID, "", "")))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := LoginRequestEndpoint(c)

	if assert.NoError(t, err) {
		acc := getAccount(t)
		assert.NotEqual(t, int64(0), acc.Token)
		assert.Equal(t, status, rec.Result().StatusCode)
	}
}

func loginStep2(t *testing.T, token string, status int, validate bool) {

	url := fmt.Sprintf("/login/%s", token)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	r := e.Router()
	r.Add(http.MethodGet, loginConfirmationRoute, LoginConfirmationEndpoint)

	c := e.NewContext(req, rec)
	r.Find(http.MethodGet, url, c)
	err := LoginConfirmationEndpoint(c)

	if assert.NoError(t, err) {
		assert.Equal(t, status, rec.Result().StatusCode)
		if validate {
			acc := getAccount(t)
			assert.NotEqual(t, int64(0), acc.Confirmed)
			assert.Equal(t, account.AccountLoggedOut, acc.Status)
			assert.NotEqual(t, int64(0), acc.Token)
		}
	}
}

func loginStep3(t *testing.T, testRealm, testUser, testClient, testToken string, state, status int, validate bool) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, getAuthorizationRoute, strings.NewReader(createAuthRequestJSON(testRealm, testUser, "", testToken)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	err := GetAuthorizationEndpoint(c)

	if assert.NoError(t, err) {
		assert.Equal(t, status, rec.Result().StatusCode)
		if validate {
			acc := getAccount(t)
			assert.Equal(t, state, acc.Status)
			assert.Equal(t, "", acc.Token)
		}
	}
}

func logoutStep(t *testing.T, testRealm, testUser, testClient, testToken string, status int, validate bool) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, logoutRequestRoute, strings.NewReader(createAuthRequestJSON(testRealm, testUser, testClient, "")))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+testToken)
	c := e.NewContext(req, rec)

	err := LogoutRequestEndpoint(c)

	if assert.NoError(t, err) {
		assert.Equal(t, status, rec.Result().StatusCode)
		if validate {
			acc := getAccount(t)
			assert.Equal(t, account.AccountLoggedOut, acc.Status)

		}
	}
}

func createAuthRequestJSON(real, user, client, token string) string {
	return fmt.Sprintf(`{"realm":"%s","user_id":"%s","client_id":"%s","token":"%s"}`, realm, user, client, token)
}

func cleaner() {
	acc, err := account.FindAccountByUserID(context.TODO(), realm, userID)
	if err == nil && acc != nil {
		account.DeleteAccount(context.TODO(), acc.Realm, acc.ClientID)

		k := nativeKey(namedKey(realm, acc.ClientID))
		ds.DataStore().Delete(context.TODO(), k)
	}
}

func getAccount(t *testing.T) *account.Account {
	acc, err := account.FindAccountByUserID(context.TODO(), realm, userID)
	if assert.NoError(t, err) {
		if assert.NotNil(t, acc) {
			return acc
		}
	}
	t.FailNow()
	return nil
}

func verifyAccountAndAuth(t *testing.T) {
	acc, err := account.FindAccountByUserID(context.TODO(), realm, userID)
	if err == nil && acc != nil {
		auth, err := LookupAuthorization(context.TODO(), acc.Realm, acc.ClientID)
		if err == nil && auth != nil {
			assert.Equal(t, acc.ClientID, auth.ClientID)
		}
	}
}
