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
		ath, _ := LookupAuthorization(ctx, accountTestRealm, acc.ClientID)
		if ath != nil {
			DeleteAuthorization(ctx, accountTestRealm, acc.ClientID)
		}
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
