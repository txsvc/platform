package platform

import (
	"context"
	"log"

	"github.com/txsvc/platform/errorreporting"
	"github.com/txsvc/platform/logging"
)

type (
	/*
		HttpRequestContext interface {
			NewHttpContext(echo.Context) context.Context
		}
	*/

	ProviderType     int
	InstanceProvider func(string) interface{}

	Platform struct {
		logger        map[string]logging.LoggingProvider
		errorReporter errorreporting.ErrorReportingProvider

		providers map[ProviderType]PlatformOpts
	}

	PlatformOpts struct {
		ID   string
		Type ProviderType
		Impl InstanceProvider
	}
)

const (
	ProviderTypeLogger ProviderType = iota
	ProviderTypeErrorReporter
	ProviderTypeHttpContext
)

var (
	platform *Platform
)

func init() {
	dl := PlatformOpts{ID: "platform.logger.default", Type: ProviderTypeLogger, Impl: logging.NewDefaultLogger}
	er := PlatformOpts{ID: "platform.errorreporting.default", Type: ProviderTypeErrorReporter, Impl: errorreporting.NewDefaultErrorReporter}

	p, err := InitPlatform(context.TODO(), dl, er)
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
			p.errorReporter = opt.Impl(opt.ID).(*errorreporting.ErrorReporter)
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

func Logger(logID string) logging.LoggingProvider {
	l, ok := platform.logger[logID]
	if !ok {
		opt, ok := platform.providers[ProviderTypeLogger]
		if !ok {
			return nil
		}
		l = opt.Impl(logID).(*logging.Logger)
		platform.logger[logID] = l
	}
	return l
}

func ReportError(e error) {
	platform.errorReporter.ReportError(e)
}
