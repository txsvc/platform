package provider

import (
	"context"
)

type (
	AuthenticationProvider interface {
		// Send an account challenge to confirm the account
		AccountChallengeNotification(context.Context, string, string) error
		// Send the new token
		ProvideAuthorizationToken(context.Context, string, string, string) error
		// Options returns the provider configuration
		Options() *AuthenticationProviderConfig
	}

	AuthenticationProviderConfig struct {
		Scope                    string
		Endpoint                 string
		AuthenticationExpiration int
		AuthorizationExpiration  int
	}
)
