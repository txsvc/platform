package platform

import (
	"context"
	"fmt"
	h "net/http"

	"github.com/txsvc/platform/v2/authentication"
	"github.com/txsvc/platform/v2/errorreporting"
	"github.com/txsvc/platform/v2/http"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/metrics"
	"github.com/txsvc/platform/v2/pkg/account"
	"github.com/txsvc/platform/v2/state"
	"github.com/txsvc/platform/v2/tasks"
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
	_ GenericProvider                       = (*defaultProviderImpl)(nil)
	_ http.HttpRequestContextProvider       = (*defaultProviderImpl)(nil)
	_ errorreporting.ErrorReportingProvider = (*defaultProviderImpl)(nil)
	_ logging.LoggingProvider               = (*defaultProviderImpl)(nil)
	_ metrics.MetricsProvider               = (*defaultProviderImpl)(nil)
	_ tasks.HttpTaskProvider                = (*defaultProviderImpl)(nil)
	_ authentication.AuthenticationProvider = (*defaultProviderImpl)(nil)
	_ state.StateProvider                   = (*defaultProviderImpl)(nil)
)

// a NULL provider that does nothing but prevents NPEs in case someone forgets to actually initializa a 'real' platform provider
func newDefaultProvider(ID string) interface{} {
	return &defaultProviderImpl{}
}

// IF GenericProvider

func (np *defaultProviderImpl) Close() error {
	return nil
}

// IF http.HttpRequestContextProvider

func (np *defaultProviderImpl) NewHttpContext(req *h.Request) context.Context {
	return context.TODO()
}

// IF errorreporting.ErrorReportingProvider

func (np *defaultProviderImpl) ReportError(e error) {
}

// IF logging.LoggingProvider

func (np *defaultProviderImpl) Log(msg string, keyValuePairs ...string) {
}

func (np *defaultProviderImpl) LogWithLevel(lvl logging.Severity, msg string, keyValuePairs ...string) {
}

// IF metrics.MetricsProvider

func (np *defaultProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
}

// IF tasks.HttpTaskProvider

func (np *defaultProviderImpl) CreateHttpTask(ctx context.Context, task tasks.HttpTask) error {
	return fmt.Errorf("not implemented")
}

// IF auth.AuthenticationProvider

// AccountChallengeNotification sends a notification to the user promting to confirm the account
func (np *defaultProviderImpl) AccountChallengeNotification(ctx context.Context, account *account.Account) error {
	return nil
}

// ProvideAuthorizationToken sends a notification to the user with the current authentication token
func (np *defaultProviderImpl) ProvideAuthorizationToken(ctx context.Context, account *account.Account) error {
	return nil
}

func (np *defaultProviderImpl) Options() *authentication.AuthenticationProviderOpts {
	return &authentication.AuthenticationProviderOpts{
		Scope:                    authentication.DefaultScope,
		Endpoint:                 authentication.DefaultEndpoint,
		AuthenticationExpiration: authentication.DefaultAuthenticationExpiration,
		AuthorizationExpiration:  authentication.DefaultAuthorizationExpiration,
	}
}

// IF state.StateProvider

func (np *defaultProviderImpl) DecodeKey(encoded string) (*state.Key, error) {
	return nil, nil
}

func (np *defaultProviderImpl) NewKey(kind, key string) (*state.Key, error) {
	k := state.NewKey(kind, key)
	a := k.(*state.Key)
	return a, nil
}
