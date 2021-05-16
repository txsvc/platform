package platform

import (
	"context"
	htp "net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform/v2/pkg/apis/provider"
)

type (
	TestProviderImpl struct {
	}
)

func newTestProvider() interface{} {
	return &TestProviderImpl{}
}

func (c *TestProviderImpl) NewHttpContext(req *htp.Request) context.Context {
	return context.Background()
}

func TestWithProvider(t *testing.T) {
	opt := provider.WithProvider("test", provider.TypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	assert.Equal(t, "test", opt.ID)
	assert.Equal(t, provider.TypeLogger, opt.Type)
	assert.NotNil(t, opt.Impl)
}

func TestInitDefaultPlatform(t *testing.T) {
	reset()

	p := DefaultPlatform()

	assert.NotNil(t, p)
	assert.NotNil(t, p.errorReportingProvider)
	assert.NotNil(t, p.httpContextProvider)
	assert.NotNil(t, p.metricsProvdider)
}

func TestInitPlatform(t *testing.T) {
	reset()

	p, err := InitPlatform(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, p)

	assert.Equal(t, 0, len(p.instances))
	assert.Equal(t, 0, len(p.logger))
	assert.Equal(t, 0, len(p.providers))

	assert.Nil(t, p.errorReportingProvider)
	assert.Nil(t, p.httpContextProvider)
	assert.Nil(t, p.metricsProvdider)
}

func TestInitPlatformDuplicateProvider(t *testing.T) {
	reset()

	opt1 := provider.WithProvider("test1", provider.TypeLogger, newTestProvider)
	opt2 := provider.WithProvider("test2", provider.TypeLogger, newTestProvider)

	p, err := InitPlatform(context.Background(), opt1, opt2)
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestRegisterPlatform(t *testing.T) {
	reset()

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
	reset()

	p := DefaultPlatform()

	assert.Empty(t, p.logger)

	log := Logger("test")
	assert.NotNil(t, log)
	assert.NotEmpty(t, p.logger)
}

func TestRegisterProvider(t *testing.T) {
	reset()

	opt := provider.WithProvider("test", provider.TypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	p := DefaultPlatform()
	assert.NotNil(t, p)

	err := p.RegisterProviders(false, opt)
	assert.Error(t, err)

	err = p.RegisterProviders(true, opt)
	assert.NoError(t, err)

}

func TestGetDefaultProviders(t *testing.T) {
	reset()

	p1, ok := Provider(provider.TypeLogger)
	assert.True(t, ok)
	assert.NotNil(t, p1)

	logger := p1.(provider.LoggingProvider)
	assert.NotNil(t, logger)

	p2, ok := Provider(provider.TypeErrorReporter)
	assert.True(t, ok)
	assert.NotNil(t, p2)

	errorReporter := p2.(provider.ErrorReportingProvider)
	assert.NotNil(t, errorReporter)

	p3, ok := Provider(provider.TypeHttpContext)
	assert.True(t, ok)
	assert.NotNil(t, p3)

	httpContext := p3.(provider.HttpContextProvider)
	assert.NotNil(t, httpContext)

	p4, ok := Provider(provider.TypeMetrics)
	assert.True(t, ok)
	assert.NotNil(t, p4)

	metrics := p4.(provider.MetricsProvider)
	assert.NotNil(t, metrics)
}

func TestGetProviderFailure(t *testing.T) {
	p1, ok := Provider(provider.TypeTask)
	assert.False(t, ok)
	assert.Nil(t, p1)
}
