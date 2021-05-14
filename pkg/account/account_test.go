package account

import (
	"context"
	"testing"

	mcache "github.com/OrlovEvgeny/go-mcache"
	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform/v2/pkg/loader"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

const (
	accountTestRealm = "account_test"
	accountTestUser  = "account_test_user"
)

func cleanup() {
	account, _ := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if account != nil {
		DeleteAccount(context.TODO(), accountTestRealm, account.ClientID)
	}

	// reset the loader and cache
	accountLoader = loader.New(AccountLoaderFunc, loader.DefaultTTL)
	userIDCache = mcache.New()
}

func TestFindAccountByUserIDFail(t *testing.T) {
	cleanup()

	account, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.Nil(t, account)
	}
}

func TestCreateAccount(t *testing.T) {
	account, err := CreateAccount(context.TODO(), accountTestRealm, accountTestUser, 7)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)
		assert.Equal(t, int64(0), account.Confirmed)
		assert.Equal(t, AccountUnconfirmed, account.Status)
	}
}

func TestFindAccountByUserID(t *testing.T) {
	account, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)
	}
}

func TestDuplicateAccount(t *testing.T) {
	account, err := CreateAccount(context.TODO(), accountTestRealm, accountTestUser, 7)
	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Equal(t, ErrAccountExists, err)
}

func TestUpdateAccount(t *testing.T) {

	account1, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account1)

		now := timestamp.Now()
		account1.Confirmed = now
		err = UpdateAccount(context.TODO(), account1)

		if assert.NoError(t, err) {
			account2, err := LookupAccount(context.TODO(), account1.Realm, account1.ClientID)
			if assert.NoError(t, err) {
				assert.Equal(t, now, account2.Confirmed)
			}
		}
	}
}

func TestDeleteAccount(t *testing.T) {

	account, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)

		account1, err := DeleteAccount(context.TODO(), accountTestRealm, account.ClientID)
		if assert.NoError(t, err) {
			assert.NotNil(t, account1)

			assert.True(t, account1.Equal(account))

			account2, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
			if assert.NoError(t, err) {
				assert.Nil(t, account2)
			}

			_, err = DeleteAccount(context.TODO(), accountTestRealm, account.ClientID)
			assert.Error(t, err)
			assert.Equal(t, ErrNoSuchAccount, err)
		}
	}
}

func TestCleanup(t *testing.T) {
	cleanup()
}
