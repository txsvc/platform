package local

import (
	"context"
	"log"

	"github.com/txsvc/platform/v2"
)

var (
	DefaultLoggingConfig        platform.PlatformOpts = platform.WithProvider("platform.default.logger", platform.ProviderTypeLogger, NewDefaultLoggingProvider)
	DefaultErrorReportingConfig platform.PlatformOpts = platform.WithProvider("platform.default.errorreporting", platform.ProviderTypeErrorReporter, NewDefaultErrorReportingProvider)
	DefaultContextConfig        platform.PlatformOpts = platform.WithProvider("platform.default.context", platform.ProviderTypeHttpContext, NewDefaultContextProvider)
	DefaultTaskConfig           platform.PlatformOpts = platform.WithProvider("platform.default.task", platform.ProviderTypeTask, NewDefaultTaskProvider)
	DefaultMetricsConfig        platform.PlatformOpts = platform.WithProvider("platform.default.metrics", platform.ProviderTypeMetrics, NewDefaultMetricsProvider)
)

func InitDefaultProviders() {
	p, err := platform.InitPlatform(context.Background(), DefaultLoggingConfig, DefaultErrorReportingConfig, DefaultContextConfig, DefaultTaskConfig, DefaultMetricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	platform.RegisterPlatform(p)
}
