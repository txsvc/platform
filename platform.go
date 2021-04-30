package platform

import (
	"context"
	"log"
	"net/http"

	"github.com/txsvc/platform/provider/local"
)

type (
	ProviderType int

	LoggingProvider interface {
		Log(string, ...interface{})
	}

	ErrorReportingProvider interface {
		ReportError(error)
	}

	HttpRequestContextProvider interface {
		NewHttpContext(*http.Request) context.Context
	}

	InstanceProviderFunc func(string) interface{}

	PlatformOpts struct {
		ID   string
		Type ProviderType
		Impl InstanceProviderFunc
	}

	Platform struct {
		logger        map[string]LoggingProvider
		errorReporter ErrorReportingProvider
		httpContext   HttpRequestContextProvider

		providers map[ProviderType]PlatformOpts
	}
)

const (
	ProviderTypeLogger ProviderType = iota
	ProviderTypeErrorReporter
	ProviderTypeHttpContext
)

var (
	DefaultLoggerConfig         PlatformOpts = PlatformOpts{ID: "platform.default.logger", Type: ProviderTypeLogger, Impl: local.NewDefaultLogger}
	DefaultErrorReportingConfig PlatformOpts = PlatformOpts{ID: "platform.default.errorreporting", Type: ProviderTypeErrorReporter, Impl: local.NewDefaultErrorReporter}
	DefaultContextConfig        PlatformOpts = PlatformOpts{ID: "platform.default.context", Type: ProviderTypeHttpContext, Impl: local.NewDefaultContextProvider}

	// internal
	platform *Platform
)

func init() {

	p, err := InitPlatform(context.TODO(), DefaultLoggerConfig, DefaultErrorReportingConfig, DefaultContextConfig)
	if err != nil {
		log.Fatal(err)
	}
	platform = p
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
	default:
		panic("unsupported")
	}
}

func InitPlatform(ctx context.Context, opts ...PlatformOpts) (*Platform, error) {
	p := Platform{
		logger:    make(map[string]LoggingProvider),
		providers: make(map[ProviderType]PlatformOpts),
	}

	for _, opt := range opts {
		if _, ok := p.providers[opt.Type]; ok {
			log.Fatalf("provider of type '%s' already registered", opt.Type.String())
		}
		p.providers[opt.Type] = opt

		switch opt.Type {
		case ProviderTypeErrorReporter:
			p.errorReporter = opt.Impl(opt.ID).(ErrorReportingProvider)
		case ProviderTypeHttpContext:
			p.httpContext = opt.Impl(opt.ID).(HttpRequestContextProvider)
		}
	}
	return &p, nil
}

func RegisterPlatform(p *Platform) *Platform {
	if p == nil {
		return nil
	}
	old := platform
	platform = p
	return old
}

func Logger(logID string) LoggingProvider {
	l, ok := platform.logger[logID]
	if !ok {
		opt, ok := platform.providers[ProviderTypeLogger]
		if !ok {
			return nil
		}
		l = opt.Impl(logID).(LoggingProvider)
		platform.logger[logID] = l
	}
	return l
}

func ReportError(e error) {
	platform.errorReporter.ReportError(e)
}

func NewHttpContext(req *http.Request) context.Context {
	return platform.httpContext.NewHttpContext(req)
}
