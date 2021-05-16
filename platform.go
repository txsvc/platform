package platform

import (
	"context"
	"fmt"
	"log"
	h "net/http"

	"github.com/txsvc/platform/v2/pkg/apis/provider"
)

const (
	MsgMissingProvider = "provider '%s' required"
)

type (
	Platform struct {
		errorReportingProvider provider.ErrorReportingProvider
		metricsProvdider       provider.MetricsProvider
		httpContextProvider    provider.HttpContextProvider

		logger    map[string]provider.LoggingProvider
		providers map[provider.ProviderType]provider.ProviderConfig
		instances map[provider.ProviderType]provider.GenericProvider
	}
)

var (
	// internal
	platform *Platform
)

func init() {
	reset()
}

func reset() {
	// initialize the platform with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real platform provider
	loggingConfig := provider.WithProvider("platform.null.logger", provider.TypeLogger, provider.NewDefaultProvider)
	errorReportingConfig := provider.WithProvider("platform.null.errorreporting", provider.TypeErrorReporter, provider.NewDefaultProvider)
	contextConfig := provider.WithProvider("platform.null.context", provider.TypeHttpContext, provider.NewDefaultProvider)
	metricsConfig := provider.WithProvider("platform.null.metrics", provider.TypeMetrics, provider.NewDefaultProvider)
	authenticationConfig := provider.WithProvider("platform.null.authentication", provider.TypeAuthentication, provider.NewDefaultProvider)

	p, err := InitPlatform(context.Background(), loggingConfig, errorReportingConfig, contextConfig, metricsConfig, authenticationConfig)
	if err != nil {
		log.Fatal(err)
	}
	RegisterPlatform(p)
}

// InitPlatform creates a new platform instance and configures it with providers
func InitPlatform(ctx context.Context, opts ...provider.ProviderConfig) (*Platform, error) {
	p := Platform{
		logger:    make(map[string]provider.LoggingProvider),
		providers: make(map[provider.ProviderType]provider.ProviderConfig),
	}

	if err := p.RegisterProviders(false, opts...); err != nil {
		return nil, err
	}

	return &p, nil
}

// RegisterPlatform makes p the new default platform provider
func RegisterPlatform(p *Platform) *Platform {
	if p == nil {
		return nil
	}
	old := platform
	platform = p
	return old
}

// RegisterProviders registers one or more  providers.
// An existing provider will be overwritten if ignoreExists is true, otherwise the function returns an error.
func (p *Platform) RegisterProviders(ignoreExists bool, opts ...provider.ProviderConfig) error {
	for _, opt := range opts {

		if _, ok := p.providers[opt.Type]; ok {
			if !ignoreExists {
				return fmt.Errorf("provider of type '%s' already registered", opt.Type.String())
			}
		}
		p.providers[opt.Type] = opt

		switch opt.Type {
		case provider.TypeErrorReporter:
			p.errorReportingProvider = opt.Impl().(provider.ErrorReportingProvider)
		case provider.TypeHttpContext:
			p.httpContextProvider = opt.Impl().(provider.HttpContextProvider)
		case provider.TypeMetrics:
			p.metricsProvdider = opt.Impl().(provider.MetricsProvider)
		}
	}
	return nil
}

// Close iterates over all registered providers and shuts them down.
func (p *Platform) Close() error {
	hasError := false
	for _, provider := range p.instances {
		if err := provider.Close(); err != nil {
			hasError = true
		}
	}
	if hasError {
		return fmt.Errorf("error(s) closing all providers")
	}
	return nil
}

// DefaultPlatform returns the current default platform provider.
func DefaultPlatform() *Platform {
	return platform
}

// Close asks all registered providers of the current default platform instance to gracefully shutdown.
func Close() error {
	return platform.Close()
}

// Provider returns the registered provider instance if it is defined.
// The bool flag is set to true if there is a provider and false otherwise.
func Provider(providerType provider.ProviderType) (interface{}, bool) {
	opt, ok := platform.providers[providerType]
	if !ok {
		return nil, false
	}
	return opt.Impl(), true
}

// a set of convenience functions in order to avoid getting the provider impl every time

// Logger returns a logger instance identified by ID
func Logger(logID string) provider.LoggingProvider {
	l, ok := platform.logger[logID]
	if !ok {
		opt, ok := platform.providers[provider.TypeLogger]
		if !ok {
			return nil
		}
		l = opt.Impl().(provider.LoggingProvider)
		platform.logger[logID] = l
	}
	return l
}

// Meter logs args to a metrics log from where the values can be aggregated and analyzed.
func Meter(ctx context.Context, metric string, args ...string) {
	platform.metricsProvdider.Meter(ctx, metric, args...)
}

// ReportError reports error e using the current platform's error reporting provider
func ReportError(e error) {
	platform.errorReportingProvider.ReportError(e)
}

// NewHttpContext creates a new Http context for request req
func NewHttpContext(req *h.Request) context.Context {
	return platform.httpContextProvider.NewHttpContext(req)
}
