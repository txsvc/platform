package provider

import (
	"context"

	"github.com/txsvc/platform/v2/pkg/account"
)

type (
	AuthenticationProvider interface {
		// Send an account challenge to confirm the account
		AccountChallengeNotification(context.Context, *account.Account) error
		// Send the new token
		ProvideAuthorizationToken(context.Context, *account.Account) error
		// Options returns the provider configuration
		Options() *AuthenticationProviderOpts
	}

	AuthenticationProviderOpts struct {
		Scope                    string
		Endpoint                 string
		AuthenticationExpiration int
		AuthorizationExpiration  int
	}
)
