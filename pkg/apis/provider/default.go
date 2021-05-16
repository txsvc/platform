package provider

import (
	"context"
	h "net/http"

	"github.com/txsvc/platform/v2/pkg/account"
)

type (
	// defaultProviderImpl provides a default implementation in the absence of any other configuration.
	defaultProviderImpl struct {
	}
)

var (
	// Interface guards.

	// This enforces a compile-time check of the provider implmentation,
	// making sure all the methods defined in the provider interfaces are implemented.

	_ GenericProvider        = (*defaultProviderImpl)(nil)
	_ HttpContextProvider    = (*defaultProviderImpl)(nil)
	_ ErrorReportingProvider = (*defaultProviderImpl)(nil)
	_ LoggingProvider        = (*defaultProviderImpl)(nil)
	_ MetricsProvider        = (*defaultProviderImpl)(nil)
)

// a NULL provider that does nothing but prevents NPEs in case someone forgets to actually initializa a 'real' platform provider
func NewDefaultProvider() interface{} {
	return &defaultProviderImpl{}
}

// IF GenericProvider

func (np *defaultProviderImpl) Close() error {
	return nil
}

// IF HttpRequestContextProvider

func (np *defaultProviderImpl) NewHttpContext(req *h.Request) context.Context {
	return context.TODO()
}

// IF ErrorReportingProvider

func (np *defaultProviderImpl) ReportError(e error) {
}

// IF LoggingProvider

func (np *defaultProviderImpl) Log(msg string, keyValuePairs ...string) {
}

func (np *defaultProviderImpl) LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
}

// IF MetricsProvider

func (np *defaultProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
}

// IF AuthenticationProvider

// AccountChallengeNotification sends a notification to the user promting to confirm the account
func (a *defaultProviderImpl) AccountChallengeNotification(ctx context.Context, account *account.Account) error {
	return nil
}

// ProvideAuthorizationToken sends a notification to the user with the current authentication token
func (a *defaultProviderImpl) ProvideAuthorizationToken(ctx context.Context, account *account.Account) error {
	return nil
}

func (a *defaultProviderImpl) Options() *AuthenticationProviderOpts {
	return &AuthenticationProviderOpts{
		Scope:                    "",
		Endpoint:                 "http://localhost:8080",
		AuthenticationExpiration: 10, // min
		AuthorizationExpiration:  90, // days
	}
}
