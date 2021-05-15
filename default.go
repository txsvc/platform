package platform

import (
	"context"
	h "net/http"

	"github.com/txsvc/platform/v2/errorreporting"
	"github.com/txsvc/platform/v2/http"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/metrics"
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
)

// a NULL provider that does nothing but prevents NPEs in case someone forgets to actually initializa a 'real' platform provider
func newDefaultProvider() interface{} {
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
