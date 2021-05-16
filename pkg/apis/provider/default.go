package provider

import (
	"context"
	h "net/http"
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

// IF provider.HttpRequestContextProvider

func (np *defaultProviderImpl) NewHttpContext(req *h.Request) context.Context {
	return context.TODO()
}

// IF provider.ErrorReportingProvider

func (np *defaultProviderImpl) ReportError(e error) {
}

// IF provider.LoggingProvider

func (np *defaultProviderImpl) Log(msg string, keyValuePairs ...string) {
}

func (np *defaultProviderImpl) LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
}

// IF metrics.MetricsProvider

func (np *defaultProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
}
