package authentication

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txsvc/platform/v2/pkg/account"
)

const (
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

func createUnconfirmedUser(t *testing.T, expires int) *account.Account {
	ctx := context.TODO()

	// create an account

	acc, err := account.CreateAccount(ctx, accountTestRealm, accountTestUser, expires)
	if err != nil {
		panic(err)
	}

	assert.NotNil(t, acc)

	// make sure the account is unconfirmed
	assert.Equal(t, int64(0), acc.Confirmed)
	assert.Greater(t, acc.Expires, int64(0))
	assert.Equal(t, account.AccountUnconfirmed, acc.Status)
	assert.NotEmpty(t, acc.Token)

	return acc
}

func createActiveUser() {
	ctx := context.TODO()

	// create an account

	acc, err := account.CreateAccount(ctx, accountTestRealm, accountTestUser, 7)
	if err != nil {
		panic(err)
	}

	req := AuthorizationRequest{
		Realm:    accountTestRealm,
		UserID:   accountTestUser,
		ClientID: acc.ClientID,
		Scope:    DefaultScope,
	}

	// create a matching authorization

	ath := NewAuthorization(&req, 10)
	err = UpdateAuthorization(ctx, ath)
	if err != nil {
		panic(err)
	}

	// login
	acc.Status = account.AccountActive
	err = account.UpdateAccount(ctx, acc)
	if err != nil {
		panic(err)
	}
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
