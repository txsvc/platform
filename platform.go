package platform

import (
	"context"
	"fmt"
	"log"
	h "net/http"

	//"github.com/txsvc/platform/v2/auth"
	"github.com/txsvc/platform/v2/errorreporting"
	"github.com/txsvc/platform/v2/http"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/metrics"
	"github.com/txsvc/platform/v2/tasks"
)

const (
	ProviderTypeLogger ProviderType = iota
	ProviderTypeErrorReporter
	ProviderTypeHttpContext
	ProviderTypeTask
	ProviderTypeMetrics
	ProviderTypeAuth
)

type (
	ProviderType int

	InstanceProviderFunc func(string) interface{}

	GenericProvider interface {
		Close() error
	}

	PlatformOpts struct {
		ID   string
		Type ProviderType
		Impl InstanceProviderFunc
	}

	Platform struct {
		errorReportingProvider errorreporting.ErrorReportingProvider
		httpContextProvider    http.HttpRequestContextProvider
		backgroundTaskProvider tasks.HttpTaskProvider
		metricsProvdider       metrics.MetricsProvider

		logger    map[string]logging.LoggingProvider
		providers map[ProviderType]PlatformOpts
		instances map[ProviderType]GenericProvider
	}
)

var (
	// internal
	platform *Platform
)

func init() {
	// initialize the platform with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real platform provider
	nullLoggingConfig := WithProvider("platform.null.logger", ProviderTypeLogger, newDefaultProvider)
	nullErrorReportingConfig := WithProvider("platform.null.errorreporting", ProviderTypeErrorReporter, newDefaultProvider)
	nullContextConfig := WithProvider("platform.null.context", ProviderTypeHttpContext, newDefaultProvider)
	nullTaskConfig := WithProvider("platform.null.task", ProviderTypeTask, newDefaultProvider)
	nullMetricsConfig := WithProvider("platform.null.metrics", ProviderTypeMetrics, newDefaultProvider)

	p, err := InitPlatform(context.Background(), nullLoggingConfig, nullErrorReportingConfig, nullContextConfig, nullTaskConfig, nullMetricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	RegisterPlatform(p)
}

// Returns the name of a provider type
func (l ProviderType) String() string {
	switch l {
	case ProviderTypeLogger:
		return "LOGGER"
	case ProviderTypeErrorReporter:
		return "ERROR_REPORTER"
	case ProviderTypeHttpContext:
		return "HTTP_CONTEXT"
	case ProviderTypeTask:
		return "TASK"
	case ProviderTypeMetrics:
		return "METRICS"
	default:
		panic("unsupported")
	}
}

// InitPlatform creates a new platform instance and configures it with providers
func InitPlatform(ctx context.Context, opts ...PlatformOpts) (*Platform, error) {
	p := Platform{
		logger:    make(map[string]logging.LoggingProvider),
		providers: make(map[ProviderType]PlatformOpts),
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
func (p *Platform) RegisterProviders(ignoreExists bool, opts ...PlatformOpts) error {
	for _, opt := range opts {

		if _, ok := p.providers[opt.Type]; ok {
			if !ignoreExists {
				return fmt.Errorf("provider of type '%s' already registered", opt.Type.String())
			}
		}
		p.providers[opt.Type] = opt

		switch opt.Type {
		case ProviderTypeErrorReporter:
			p.errorReportingProvider = opt.Impl(opt.ID).(errorreporting.ErrorReportingProvider)
		case ProviderTypeHttpContext:
			p.httpContextProvider = opt.Impl(opt.ID).(http.HttpRequestContextProvider)
		case ProviderTypeTask:
			p.backgroundTaskProvider = opt.Impl(opt.ID).(tasks.HttpTaskProvider)
		case ProviderTypeMetrics:
			p.metricsProvdider = opt.Impl(opt.ID).(metrics.MetricsProvider)
			//case ProviderTypeAuth:
			//	opt.Impl(opt.ID).(auth.AuthorizationProvider)
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

// WithProvider returns a populated PlatformOption struct.
func WithProvider(ID string, providerType ProviderType, impl InstanceProviderFunc) PlatformOpts {
	return PlatformOpts{
		ID:   ID,
		Type: providerType,
		Impl: impl,
	}
}

// Logger returns a logger instance identified by ID
func Logger(logID string) logging.LoggingProvider {
	l, ok := platform.logger[logID]
	if !ok {
		opt, ok := platform.providers[ProviderTypeLogger]
		if !ok {
			return nil
		}
		l = opt.Impl(logID).(logging.LoggingProvider)
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

// NewTask schedules a new http background task
func NewTask(task tasks.HttpTask) error {
	return platform.backgroundTaskProvider.CreateHttpTask(context.Background(), task)
}
