package local

import (
	"context"
	"log"
	h "net/http"

	"go.uber.org/zap"

	"github.com/txsvc/platform/v2"

	"github.com/txsvc/platform/v2/pkg/apis/provider"
)

type (
	LocalProviderImpl struct {
	}

	LocalLoggingProviderImpl struct {
		lvl provider.Severity
		log *zap.SugaredLogger
	}

	LocalErrorReportingProviderImpl struct {
		log *zap.SugaredLogger
	}
)

var (
	loggingConfig        provider.ProviderConfig = provider.WithProvider("platform.default.logger", provider.TypeLogger, LocalLoggingProvider)
	errorReportingConfig provider.ProviderConfig = provider.WithProvider("platform.default.errorreporting", provider.TypeErrorReporter, LocalErrorReportingProvider)
	contextConfig        provider.ProviderConfig = provider.WithProvider("platform.default.context", provider.TypeHttpContext, LocalHttpContextProvider)
	metricsConfig        provider.ProviderConfig = provider.WithProvider("platform.default.metrics", provider.TypeMetrics, LocalMetricsProvider)

	errorReportingClient *LocalErrorReportingProviderImpl

	// Interface guards
	_ provider.GenericProvider     = (*LocalProviderImpl)(nil)
	_ provider.HttpContextProvider = (*LocalProviderImpl)(nil)
	_ provider.MetricsProvider     = (*LocalProviderImpl)(nil)

	_ provider.GenericProvider        = (*LocalErrorReportingProviderImpl)(nil)
	_ provider.ErrorReportingProvider = (*LocalErrorReportingProviderImpl)(nil)

	_ provider.GenericProvider = (*LocalLoggingProviderImpl)(nil)
	_ provider.LoggingProvider = (*LocalLoggingProviderImpl)(nil)
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
	p, err := platform.InitPlatform(context.Background(), loggingConfig, errorReportingConfig, contextConfig, metricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	platform.RegisterPlatform(p)
}

func LocalHttpContextProvider() interface{} {
	return &LocalProviderImpl{}
}

func LocalMetricsProvider() interface{} {
	return &LocalProviderImpl{}
}

func (c *LocalProviderImpl) Close() error {
	return nil
}

func (c *LocalProviderImpl) NewHttpContext(req *h.Request) context.Context {
	return context.Background()
}

func LocalLoggingProvider() interface{} {
	callerSkipConf := zap.AddCallerSkip(1)

	l, err := zap.NewProduction(callerSkipConf)
	if err != nil {
		return nil
	}

	logger := LocalLoggingProviderImpl{
		lvl: provider.LevelInfo,
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

func (l *LocalLoggingProviderImpl) LogWithLevel(lvl provider.Severity, msg string, keyValuePairs ...string) {

	if len(keyValuePairs) > 0 {
		params := make([]interface{}, len(keyValuePairs))
		for i := range keyValuePairs {
			params[i] = keyValuePairs[i]
		}

		switch lvl {
		case provider.LevelInfo:
			l.log.Infow(msg, params...)
		case provider.LevelWarn:
			l.log.Warnw(msg, params...)
		case provider.LevelError:
			l.log.Errorw(msg, params...)
		case provider.LevelDebug:
			l.log.Debugw(msg, params...)
		}
	} else {
		switch lvl {
		case provider.LevelInfo:
			l.log.Infow(msg)
		case provider.LevelWarn:
			l.log.Warnw(msg)
		case provider.LevelError:
			l.log.Errorw(msg)
		case provider.LevelDebug:
			l.log.Debugw(msg)
		}
	}
}

func (er *LocalErrorReportingProviderImpl) Close() error {
	return nil
}

func LocalErrorReportingProvider() interface{} {
	return errorReportingClient
}

func (er *LocalErrorReportingProviderImpl) ReportError(e error) {
	er.log.Error(e)
}

func (m *LocalProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
	// actually does nothing right now
}
