package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txsvc/platform/v2/pkg/account"
)

const (
	scopeProductionRead  = "production:read"
	scopeProductionWrite = "production:write"
	scopeProductionBuild = "production:build"
	scopeResourceRead    = "resource:read"
	scopeResourceWrite   = "resource:write"

	accountTestRealm = "account_test"
	accountTestUser  = "account_test_user"
)

func cleanup() {
	ctx := context.TODO()

	acc, _ := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
	if acc != nil {
		DeleteAuthorization(ctx, accountTestRealm, acc.ClientID)
		account.DeleteAccount(ctx, accountTestRealm, acc.ClientID)
	}
}

func TestScope(t *testing.T) {

	scope1 := "production:read,production:write,production:build"

	assert.False(t, hasScope("", ""))
	assert.False(t, hasScope(scope1, ""))
	assert.False(t, hasScope("", scopeResourceRead))

	assert.True(t, hasScope(scope1, scopeProductionRead))
	assert.False(t, hasScope(scope1, scopeResourceRead))
}

func TestLookupAuthorizationFail(t *testing.T) {
	cleanup()
	ctx := context.TODO()

	account, err := account.CreateAccount(ctx, accountTestRealm, accountTestUser, 7)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)

		ath, err := LookupAuthorization(ctx, accountTestRealm, account.ClientID)

		if assert.NoError(t, err) {
			assert.Nil(t, ath)
		}
	}
}

func TestNewAuthorization(t *testing.T) {
	ctx := context.TODO()

	account, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)
	}

	req := AuthorizationRequest{
		Realm:    accountTestRealm,
		UserID:   accountTestUser,
		ClientID: account.ClientID,
		Scope:    DefaultScope,
	}

	ath := NewAuthorization(&req, 10)

	if assert.NotNil(t, ath) {
		assert.Equal(t, ath.ClientID, req.ClientID)
		assert.Equal(t, ath.Realm, req.Realm)
		assert.Equal(t, ath.UserID, req.UserID)
		assert.Equal(t, ath.Scope, req.Scope)
		assert.Equal(t, ath.TokenType, DefaultTokenType)
		assert.NotEmpty(t, ath.Token)
		assert.False(t, ath.Revoked)
		assert.Greater(t, ath.Expires, int64(0))
		assert.Greater(t, ath.Created, int64(0))
		assert.Greater(t, ath.Updated, int64(0))
		assert.Greater(t, ath.Expires, ath.Created)
	}

	err = UpdateAuthorization(ctx, ath)
	assert.NoError(t, err)
}

func TestLookupAuthorization(t *testing.T) {
	ctx := context.TODO()

	account, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
	assert.NoError(t, err)
	assert.NotNil(t, account)

	ath, err := LookupAuthorization(ctx, accountTestRealm, account.ClientID)
	assert.NoError(t, err)
	assert.NotNil(t, ath)

	ath, err = LookupAuthorization(ctx, accountTestRealm, "does-not-exist")
	assert.NoError(t, err)
	assert.Nil(t, ath)

}

func TestLookupAuthorizationByToken(t *testing.T) {
	ctx := context.TODO()

	account, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
	assert.NoError(t, err)
	assert.NotNil(t, account)

	ath, err := LookupAuthorization(ctx, accountTestRealm, account.ClientID)
	if assert.NoError(t, err) {
		ath1 := *ath // dereference to avoid cache issues
		ath2, err := FindAuthorizationByToken(ctx, ath1.Token)

		assert.NoError(t, err)
		assert.NotNil(t, ath2)
		assert.True(t, ath1.Equal(ath2))
	}
}

func TestUpdateAuthorization(t *testing.T) {
	ctx := context.TODO()

	account, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
	assert.NoError(t, err)
	assert.NotNil(t, account)

	ath1, err := LookupAuthorization(ctx, accountTestRealm, account.ClientID)
	assert.NoError(t, err)
	assert.NotNil(t, ath1)

	athCopy := *ath1 // dereference to avoid cache interference

	ath1.Token = "a-new-token"
	err = UpdateAuthorization(ctx, ath1)

	assert.NoError(t, err)

	ath3, err := LookupAuthorization(ctx, accountTestRealm, account.ClientID)
	assert.NoError(t, err)
	assert.NotNil(t, ath3)

	assert.False(t, ath3.Equal(&athCopy))
}

func TestDeleteAuthorization(t *testing.T) {
	t.Cleanup(cleanup)

	ctx := context.TODO()

	account, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
	assert.NoError(t, err)
	assert.NotNil(t, account)

	ath, err := DeleteAuthorization(ctx, accountTestRealm, account.ClientID)
	assert.NoError(t, err)
	assert.NotNil(t, ath)

	ath2, err := DeleteAuthorization(ctx, accountTestRealm, account.ClientID)
	assert.Error(t, err)
	assert.Nil(t, ath2)
	assert.Equal(t, err, ErrNoSuchEntity)
}
