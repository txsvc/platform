package google

import (
	"context"
	"log"

	"github.com/txsvc/platform"
)

var (
	// Google Cloud Platform
	GoogleErrorReportingConfig platform.PlatformOpts = platform.WithProvider("platform.google.errorreporting", platform.ProviderTypeErrorReporter, NewErrorReporter)
	GoogleCloudTaskConfig      platform.PlatformOpts = platform.WithProvider("platform.google.task", platform.ProviderTypeTask, NewCloudTaskProvider)
	GoogleCloudLoggingConfig   platform.PlatformOpts = platform.WithProvider("platform.google.logger", platform.ProviderTypeLogger, NewGoogleCloudLoggingProvider)
	GoogleCloudMetricsConfig   platform.PlatformOpts = platform.WithProvider("platform.google.metrics", platform.ProviderTypeMetrics, NewGoogleCloudMetricsProvider)
	// AppEngine
	AppEngineContextConfig platform.PlatformOpts = platform.WithProvider("platform.google.context", platform.ProviderTypeHttpContext, NewAppEngineContextProvider)
)

func InitDefaultGoogleProviders() {
	p, err := platform.InitPlatform(context.Background(), GoogleErrorReportingConfig, GoogleCloudTaskConfig, GoogleCloudLoggingConfig, GoogleCloudMetricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	platform.RegisterPlatform(p)
}
