package platform

import (
	"context"
	"log"
	h "net/http"

	"github.com/txsvc/platform/pkg/errorreporting"
	"github.com/txsvc/platform/pkg/http"
	"github.com/txsvc/platform/pkg/logging"
	"github.com/txsvc/platform/pkg/metrics"
	"github.com/txsvc/platform/pkg/tasks"
	"github.com/txsvc/platform/provider/local"
)

const (
	ProviderTypeLogger ProviderType = iota
	ProviderTypeErrorReporter
	ProviderTypeHttpContext
	ProviderTypeTask
	ProviderTypeMetrics
)

type (
	ProviderType int

	InstanceProviderFunc func(string) interface{}

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
	}
)

var (
	DefaultLoggingConfig        PlatformOpts = PlatformOpts{ID: "platform.default.logger", Type: ProviderTypeLogger, Impl: local.NewDefaultLoggingProvider}
	DefaultErrorReportingConfig PlatformOpts = PlatformOpts{ID: "platform.default.errorreporting", Type: ProviderTypeErrorReporter, Impl: local.NewDefaultErrorReportingProvider}
	DefaultContextConfig        PlatformOpts = PlatformOpts{ID: "platform.default.context", Type: ProviderTypeHttpContext, Impl: local.NewDefaultContextProvider}
	DefaultTaskConfig           PlatformOpts = PlatformOpts{ID: "platform.default.task", Type: ProviderTypeTask, Impl: local.NewDefaultTaskProvider}
	DefaultMetricsConfig        PlatformOpts = PlatformOpts{ID: "platform.default.metrics", Type: ProviderTypeMetrics, Impl: local.NewDefaultMetricsProvider}

	// internal
	platform *Platform
)

func init() {
	InitDefaultProviders()
}

func InitDefaultProviders() {
	p, err := InitPlatform(context.Background(), DefaultLoggingConfig, DefaultErrorReportingConfig, DefaultContextConfig, DefaultTaskConfig, DefaultMetricsConfig)
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

	for _, opt := range opts {
		if _, ok := p.providers[opt.Type]; ok {
			log.Fatalf("provider of type '%s' already registered", opt.Type.String())
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
		}
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

// DefaultPlatform returns the current default platform provider
func DefaultPlatform() *Platform {
	return platform
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

// Meter reports args
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

func NewTask(task tasks.HttpTask) error {
	return platform.backgroundTaskProvider.CreateHttpTask(context.Background(), task)
}
