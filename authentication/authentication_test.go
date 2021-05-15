package authentication

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txsvc/platform/v2/pkg/account"
)

func TestConfirmLoginChallenge(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	ctx := context.TODO()

	acc := createUnconfirmedUser(t, 10)
	acc2, status, err := ConfirmLoginChallenge(ctx, acc.Token)

	assert.NoError(t, err)
	assert.NotNil(t, acc2)
	assert.Equal(t, http.StatusNoContent, status)

	assert.Equal(t, int64(0), acc2.Expires)
	assert.Greater(t, acc2.Confirmed, int64(0))
	assert.Equal(t, account.AccountLoggedOut, acc2.Status)
	assert.Empty(t, acc2.Token)
}

func TestConfirmLoginChallengeNoToken(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	ctx := context.TODO()

	createUnconfirmedUser(t, 10)
	acc2, status, err := ConfirmLoginChallenge(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, acc2)
	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Equal(t, ErrNoToken, err)
}

func TestConfirmLoginChallengeInvalidToken(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	ctx := context.TODO()

	createUnconfirmedUser(t, 10)
	acc2, status, err := ConfirmLoginChallenge(ctx, "invalid-token")

	assert.NoError(t, err)
	assert.Nil(t, acc2)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestConfirmLoginChallengeTokenExpired(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	ctx := context.TODO()

	acc := createUnconfirmedUser(t, -10) // expires < 0 -> in the past
	acc2, status, err := ConfirmLoginChallenge(ctx, acc.Token)

	assert.NoError(t, err)
	assert.NotNil(t, acc2)
	assert.Equal(t, http.StatusForbidden, status)

}

func TestLogoutAccount(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	ctx := context.TODO()

	createActiveUser()

	acc, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, acc)
		assert.Equal(t, account.AccountActive, acc.Status)

		status, err := LogoutAccount(ctx, accountTestRealm, acc.ClientID)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, status)

		acc1, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)

		assert.NoError(t, err)
		assert.NotNil(t, acc1)
		assert.Equal(t, acc1.Status, account.AccountLoggedOut)

		ath, err := LookupAuthorization(ctx, accountTestRealm, acc1.ClientID)

		assert.NoError(t, err)
		assert.NotNil(t, ath)

		assert.False(t, ath.IsValid())
		assert.True(t, ath.Revoked)
	}
}

func TestBlockAccount(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	ctx := context.TODO()

	createActiveUser()

	acc, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, acc)
		assert.Equal(t, account.AccountActive, acc.Status)

		err := BlockAccount(ctx, accountTestRealm, acc.ClientID)

		assert.NoError(t, err)

		acc1, err := account.FindAccountByUserID(ctx, accountTestRealm, accountTestUser)

		assert.NoError(t, err)
		assert.NotNil(t, acc1)
		assert.Equal(t, acc1.Status, account.AccountBlocked)

		ath, err := LookupAuthorization(ctx, accountTestRealm, acc1.ClientID)

		assert.NoError(t, err)
		assert.NotNil(t, ath)

		assert.False(t, ath.IsValid())
		assert.True(t, ath.Revoked)
	}
}
