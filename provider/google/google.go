package google

import (
	"github.com/txsvc/platform"
)

var (
	// Google Cloud Platform
	GoogleErrorReportingConfig platform.PlatformOpts = platform.PlatformOpts{ID: "platform.google.errorreporting", Type: platform.ProviderTypeErrorReporter, Impl: NewErrorReporter}
	GoogleCloudTaskConfig      platform.PlatformOpts = platform.PlatformOpts{ID: "platform.google.task", Type: platform.ProviderTypeTask, Impl: NewCloudTaskProvider}
	GoogleCloudLoggingConfig   platform.PlatformOpts = platform.PlatformOpts{ID: "platform.google.logger", Type: platform.ProviderTypeLogger, Impl: NewGoogleCloudLoggingProvider}
	GoogleCloudMetricsConfig   platform.PlatformOpts = platform.PlatformOpts{ID: "platform.google.metrics", Type: platform.ProviderTypeMetrics, Impl: NewGoogleCloudMetricsProvider}

	// Google AppEngine
	AppEngineContextConfig platform.PlatformOpts = platform.PlatformOpts{ID: "platform.google.context", Type: platform.ProviderTypeHttpContext, Impl: NewAppEngineContextProvider}
)
