package local

import (
	"context"
	"log"
	h "net/http"

	"go.uber.org/zap"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/errorreporting"
	"github.com/txsvc/platform/v2/http"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/metrics"
)

type (
	LocalProviderImpl struct {
	}

	LocalLoggingProviderImpl struct {
		lvl logging.Severity
		log *zap.SugaredLogger
	}

	LocalErrorReportingProviderImpl struct {
		log *zap.SugaredLogger
	}
)

var (
	DefaultLoggingConfig        platform.PlatformOpts = platform.WithProvider("platform.default.logger", platform.ProviderTypeLogger, NewLocalLoggingProvider)
	DefaultErrorReportingConfig platform.PlatformOpts = platform.WithProvider("platform.default.errorreporting", platform.ProviderTypeErrorReporter, NewLocalErrorReportingProvider)
	DefaultContextConfig        platform.PlatformOpts = platform.WithProvider("platform.default.context", platform.ProviderTypeHttpContext, NewLocalProvider)
	DefaultMetricsConfig        platform.PlatformOpts = platform.WithProvider("platform.default.metrics", platform.ProviderTypeMetrics, NewLocalProvider)

	errorReportingClient *LocalErrorReportingProviderImpl

	// Interface guards
	_ platform.GenericProvider        = (*LocalProviderImpl)(nil)
	_ http.HttpRequestContextProvider = (*LocalProviderImpl)(nil)
	_ metrics.MetricsProvider         = (*LocalProviderImpl)(nil)

	_ platform.GenericProvider              = (*LocalErrorReportingProviderImpl)(nil)
	_ errorreporting.ErrorReportingProvider = (*LocalErrorReportingProviderImpl)(nil)

	_ platform.GenericProvider = (*LocalLoggingProviderImpl)(nil)
	_ logging.LoggingProvider  = (*LocalLoggingProviderImpl)(nil)
)

func init() {
	callerSkipConf := zap.AddCallerSkip(2)
	l, err := zap.NewProduction(callerSkipConf)

	if err != nil {
		log.Fatal(err)
	}
	er := LocalErrorReportingProviderImpl{
		log: l.Sugar(),
	}

	errorReportingClient = &er
}

func InitLocalProviders() {
	p, err := platform.InitPlatform(context.Background(), DefaultLoggingConfig, DefaultErrorReportingConfig, DefaultContextConfig, DefaultMetricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	platform.RegisterPlatform(p)
}

func NewLocalProvider() interface{} {
	return &LocalProviderImpl{}
}

func (c *LocalProviderImpl) Close() error {
	return nil
}

func (c *LocalProviderImpl) NewHttpContext(req *h.Request) context.Context {
	return context.Background()
}

func NewLocalLoggingProvider() interface{} {
	callerSkipConf := zap.AddCallerSkip(1)

	l, err := zap.NewProduction(callerSkipConf)
	if err != nil {
		return nil
	}

	logger := LocalLoggingProviderImpl{
		lvl: logging.Info,
		log: l.Sugar(),
	}

	return &logger
}

func (l *LocalLoggingProviderImpl) Close() error {
	return nil
}

func (l *LocalLoggingProviderImpl) Log(msg string, keyValuePairs ...string) {
	l.LogWithLevel(l.lvl, msg, keyValuePairs...)
}

func (l *LocalLoggingProviderImpl) LogWithLevel(lvl logging.Severity, msg string, keyValuePairs ...string) {

	if len(keyValuePairs) > 0 {
		params := make([]interface{}, len(keyValuePairs))
		for i := range keyValuePairs {
			params[i] = keyValuePairs[i]
		}

		switch lvl {
		case logging.Info:
			l.log.Infow(msg, params...)
		case logging.Warn:
			l.log.Warnw(msg, params...)
		case logging.Error:
			l.log.Errorw(msg, params...)
		case logging.Debug:
			l.log.Debugw(msg, params...)
		}
	} else {
		switch lvl {
		case logging.Info:
			l.log.Infow(msg)
		case logging.Warn:
			l.log.Warnw(msg)
		case logging.Error:
			l.log.Errorw(msg)
		case logging.Debug:
			l.log.Debugw(msg)
		}
	}
}

func (er *LocalErrorReportingProviderImpl) Close() error {
	return nil
}

func NewLocalErrorReportingProvider() interface{} {
	return errorReportingClient
}

func (er *LocalErrorReportingProviderImpl) ReportError(e error) {
	er.log.Error(e)
}

func (m *LocalProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
	// actually does nothing right now
}
