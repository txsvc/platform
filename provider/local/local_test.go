package local

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/authentication"
	"github.com/txsvc/platform/v2/http"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/metrics"
)

func TestRegisterPlatform(t *testing.T) {
	dl := platform.PlatformOpts{ID: "platform.logger.default", Type: platform.ProviderTypeLogger, Impl: LocalLoggingProvider}
	p, err := platform.InitPlatform(context.Background(), dl)

	if assert.NoError(t, err) {
		assert.NotNil(t, p)

		platform.RegisterPlatform(p)
		logger := platform.Logger("somelogger")

		assert.NotNil(t, logger)
	}
}

func TestDefaultContext(t *testing.T) {
	InitLocalProviders()

	ctx := platform.NewHttpContext(nil)
	assert.NotNil(t, ctx)
}

func TestLoggingProvider(t *testing.T) {
	InitLocalProviders()

	logger := platform.Logger("platform-test-logs")
	assert.NotNil(t, logger)

	logger.Log("something happened")
}

func TestDefaultLoggerWithLevel(t *testing.T) {
	InitLocalProviders()

	logger := platform.Logger("platform-test-logs")
	assert.NotNil(t, logger)

	logger.LogWithLevel(logging.Info, "something happened with level INFO")
	logger.LogWithLevel(logging.Warn, "something happened with level WARN")
	logger.LogWithLevel(logging.Error, "something happened with level ERROR")
	logger.LogWithLevel(logging.Debug, "something happened with level DEBUG")
}

func TestLoggingWithParams(t *testing.T) {
	InitLocalProviders()

	logger := platform.Logger("platform-test-logs")
	assert.NotNil(t, logger)

	logger.LogWithLevel(logging.Info, "something with parameters happened", "foo", "bar", "question", fmt.Sprintf("%d", 42), "orphan", fmt.Sprintf("%v", true))
}

func TestErrorReportingProvider(t *testing.T) {
	InitLocalProviders()

	err := fmt.Errorf("something went wrong")

	platform.ReportError(err)
}

func TestHttpContextProvider(t *testing.T) {
	InitLocalProviders()

	p, ok := platform.Provider(platform.ProviderTypeHttpContext)
	assert.True(t, ok)
	assert.NotNil(t, p)

	httpContext := p.(http.HttpRequestContextProvider)
	assert.NotNil(t, httpContext)

	ctx := httpContext.NewHttpContext(nil)
	assert.NotNil(t, ctx)
}

func TestMetricsProvider(t *testing.T) {
	InitLocalProviders()

	p, ok := platform.Provider(platform.ProviderTypeMetrics)
	assert.True(t, ok)
	assert.NotNil(t, p)

	metrics := p.(metrics.MetricsProvider)
	assert.NotNil(t, metrics)
}

func TestAuthenticationProvider(t *testing.T) {
	InitLocalProviders()
	platform.DefaultPlatform().RegisterProviders(false, authenticationConfig)

	p, ok := platform.Provider(platform.ProviderTypeAuthentication)
	assert.True(t, ok)
	assert.NotNil(t, p)

	auth := p.(authentication.AuthenticationProvider)
	assert.NotNil(t, auth)

	assert.NotEmpty(t, auth.Options())
}
