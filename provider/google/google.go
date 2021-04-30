package google

import (
	"github.com/txsvc/platform"
)

var (
	GoogleErrorReportingConfig platform.PlatformOpts = platform.PlatformOpts{ID: "platform.google.errorreporting", Type: platform.ProviderTypeErrorReporter, Impl: NewErrorReporter}
	GoogleCloudTaskConfig      platform.PlatformOpts = platform.PlatformOpts{ID: "platform.google.task", Type: platform.ProviderTypeTask, Impl: NewCloudTaskProvider}

	AppEngineContextConfig platform.PlatformOpts = platform.PlatformOpts{ID: "platform.google.context", Type: platform.ProviderTypeHttpContext, Impl: NewAppEngineContextProvider}
)
