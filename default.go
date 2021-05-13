package platform

import (
	"context"
	"fmt"
	h "net/http"

	"github.com/txsvc/platform/v2/auth"
	"github.com/txsvc/platform/v2/errorreporting"
	"github.com/txsvc/platform/v2/http"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/metrics"
	"github.com/txsvc/platform/v2/pkg/account"
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
	_ auth.AuthorizationProvider            = (*defaultProviderImpl)(nil)
)

// a NULL provider that does nothing but prevents NPEs in case someone forgets to actually initializa a 'real' platform provider
func newDefaultProvider(ID string) interface{} {
	return &defaultProviderImpl{}
}

func (np *defaultProviderImpl) Close() error {
	return nil
}

func (np *defaultProviderImpl) NewHttpContext(req *h.Request) context.Context {
	return context.TODO()
}

func (np *defaultProviderImpl) ReportError(e error) {
}

func (np *defaultProviderImpl) Log(msg string, keyValuePairs ...string) {
}

func (np *defaultProviderImpl) LogWithLevel(lvl logging.Severity, msg string, keyValuePairs ...string) {
}

func (np *defaultProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
}

func (np *defaultProviderImpl) CreateHttpTask(ctx context.Context, task tasks.HttpTask) error {
	return fmt.Errorf("not implemented")
}

// SendAccountChallenge sends a notification to the user promting to confirm the account
func (np *defaultProviderImpl) SendAccountChallenge(ctx context.Context, account *account.Account) error {
	return nil
}

// SendAuthToken sends a notification to the user with the current authentication token
func (np *defaultProviderImpl) SendAuthToken(ctx context.Context, account *account.Account) error {
	return nil
}

func (np *defaultProviderImpl) Scope() string {
	return auth.DefaultScope
}

func (np *defaultProviderImpl) Endpoint() string {
	return auth.DefaultEndpoint
}

func (np *defaultProviderImpl) AuthenticationExpiration() int {
	return auth.DefaultAuthenticationExpiration
}

func (np *defaultProviderImpl) AuthorizationExpiration() int {
	return auth.DefaultAuthorizationExpiration
}
