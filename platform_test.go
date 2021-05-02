package platform

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	TestProviderImpl struct {
	}
)

func newTestProvider(ID string) interface{} {
	return &TestProviderImpl{}
}

func (c *TestProviderImpl) NewHttpContext(req *http.Request) context.Context {
	return context.Background()
}

func resetPlatform() {
	// initialize the platform with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real platform provider
	nullLoggingConfig := WithProvider("platform.null.logger", ProviderTypeLogger, newNullProvider)
	nullErrorReportingConfig := WithProvider("platform.null.errorreporting", ProviderTypeErrorReporter, newNullProvider)
	nullContextConfig := WithProvider("platform.null.context", ProviderTypeHttpContext, newNullProvider)
	nullTaskConfig := WithProvider("platform.null.task", ProviderTypeTask, newNullProvider)
	nullMetricsConfig := WithProvider("platform.null.metrics", ProviderTypeMetrics, newNullProvider)

	p, _ := InitPlatform(context.Background(), nullLoggingConfig, nullErrorReportingConfig, nullContextConfig, nullTaskConfig, nullMetricsConfig)
	platform = p
}

func TestWithProvider(t *testing.T) {
	opt := WithProvider("test", ProviderTypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	assert.Equal(t, "test", opt.ID)
	assert.Equal(t, ProviderTypeLogger, opt.Type)
	assert.NotNil(t, opt.Impl)
}

func TestInitDefaultPlatform(t *testing.T) {
	resetPlatform()

	p := DefaultPlatform()

	assert.NotNil(t, p)
	assert.NotNil(t, p.errorReportingProvider)
	assert.NotNil(t, p.httpContextProvider)
	assert.NotNil(t, p.backgroundTaskProvider)
	assert.NotNil(t, p.metricsProvdider)

}

func TestInitPlatform(t *testing.T) {
	resetPlatform()

	p, err := InitPlatform(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, p)
}

func TestInitPlatformDuplicateProvider(t *testing.T) {
	resetPlatform()

	opt1 := WithProvider("test1", ProviderTypeLogger, newTestProvider)
	opt2 := WithProvider("test2", ProviderTypeLogger, newTestProvider)

	p, err := InitPlatform(context.Background(), opt1, opt2)
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestRegisterPlatform(t *testing.T) {
	resetPlatform()

	defaultPlatform := DefaultPlatform()
	p, err := InitPlatform(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, defaultPlatform)
	assert.NotNil(t, p)

	assert.NotEqual(t, defaultPlatform, p)

	old := RegisterPlatform(p)
	assert.NotNil(t, old)
	assert.Equal(t, defaultPlatform, old)

	newDefaultPlatform := DefaultPlatform()
	assert.NotNil(t, newDefaultPlatform)
	assert.Equal(t, newDefaultPlatform, p)
}

func TestInitLogger(t *testing.T) {
	resetPlatform()

	p := DefaultPlatform()

	assert.Empty(t, p.logger)

	log := Logger("test")
	assert.NotNil(t, log)
	assert.NotEmpty(t, p.logger)
}

func TestRegisterProvider(t *testing.T) {
	resetPlatform()

	opt := WithProvider("test", ProviderTypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	p := DefaultPlatform()
	assert.NotNil(t, p)

	err := p.RegisterProviders(false, opt)
	assert.Error(t, err)

	err = p.RegisterProviders(true, opt)
	assert.NoError(t, err)

}
