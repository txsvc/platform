package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecondaryCacheEntryFail(t *testing.T) {
	cleanup()

	account, err := CreateAccount(context.TODO(), accountTestRealm, accountTestUser, 7)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)

		key := namedKey(accountTestRealm, accountTestUser)
		_, ok := userIDCache.Get(key)
		assert.False(t, ok)
	}
}

func TestAccountLoaderFunc(t *testing.T) {

	account1, err := FindAccountByUserID(context.TODO(), accountTestRealm, accountTestUser)
	if assert.NoError(t, err) {
		assert.NotNil(t, account1)
	}

	k := nativeKey(account1.Key())

	account2, err := AccountLoaderFunc(context.TODO(), k.Encode())
	if assert.NoError(t, err) {
		assert.NotNil(t, account2)

		account3 := account2.(*Account)
		assert.NotNil(t, account3)

		assert.True(t, account1.Equal(account3))
	}
}

func TestLookupAccountCacheHitMiss(t *testing.T) {
	cleanup()

	account, err := CreateAccount(context.TODO(), accountTestRealm, accountTestUser, 7)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)

		hits := accountLoader.Hits()
		misses := accountLoader.Misses()

		LookupAccount(context.TODO(), accountTestRealm, account.ClientID)

		assert.Equal(t, int64(0), hits)
		assert.Equal(t, int64(1), misses)

		key := namedKey(accountTestRealm, accountTestUser)
		_, ok := userIDCache.Get(key)
		assert.True(t, ok)

		hits = accountLoader.Hits()
		misses = accountLoader.Misses()

		LookupAccount(context.TODO(), accountTestRealm, account.ClientID)

		assert.Equal(t, hits+1, accountLoader.Hits())
		assert.Equal(t, misses, accountLoader.Misses())
	}
}

func TestUpdateAccountCacheInvalidation(t *testing.T) {
	cleanup()

	account, err := CreateAccount(context.TODO(), accountTestRealm, accountTestUser, 7)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)

		hits := accountLoader.Hits()
		misses := accountLoader.Misses()

		LookupAccount(context.TODO(), accountTestRealm, account.ClientID)

		assert.Equal(t, int64(0), hits)
		assert.Equal(t, int64(1), misses)

		UpdateAccount(context.TODO(), account)

		hits = accountLoader.Hits()
		misses = accountLoader.Misses()

		// UpdateAccount should invalidate the loader, i.e. we expect a miss
		LookupAccount(context.TODO(), accountTestRealm, account.ClientID)

		assert.Equal(t, hits, accountLoader.Hits())
		assert.Equal(t, misses+1, accountLoader.Misses())
	}
}

func TestFindByUserID(t *testing.T) {
	cleanup()
	ctx := context.TODO()

	account, err := CreateAccount(ctx, accountTestRealm, accountTestUser, 7)
	if assert.NoError(t, err) {
		assert.NotNil(t, account)

		// asume secondary cache is empty
		key := namedKey(accountTestRealm, accountTestUser)
		_, ok := userIDCache.Get(key)
		assert.False(t, ok)

		// assume accountLoader is empty
		k := nativeKey(namedKey(accountTestRealm, account.ClientID))
		assert.False(t, accountLoader.Contains(ctx, k.Encode()))

		// first time, nothing is in the caches
		account1, err := FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
		if assert.NoError(t, err) {
			assert.NotNil(t, account1)
		}

		// check again
		_, ok = userIDCache.Get(key)
		assert.True(t, ok)
		assert.True(t, accountLoader.Contains(ctx, k.Encode()))

		// second find should come from the cache
		hits := accountLoader.Hits()
		misses := accountLoader.Misses()

		account2, err := FindAccountByUserID(ctx, accountTestRealm, accountTestUser)
		if assert.NoError(t, err) {
			assert.NotNil(t, account2)
		}

		assert.Equal(t, hits+1, accountLoader.Hits())
		assert.Equal(t, misses, accountLoader.Misses())

	}
}
