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

func TestWithProvider(t *testing.T) {
	opt := WithProvider("test", ProviderTypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	assert.Equal(t, "test", opt.ID)
	assert.Equal(t, ProviderTypeLogger, opt.Type)
	assert.NotNil(t, opt.Impl)
}

func TestInitDefaultPlatform(t *testing.T) {
	reset()

	p := DefaultPlatform()

	assert.NotNil(t, p)
	assert.NotNil(t, p.errorReportingProvider)
	assert.NotNil(t, p.httpContextProvider)
	assert.NotNil(t, p.backgroundTaskProvider)
	assert.NotNil(t, p.metricsProvdider)
	assert.NotNil(t, p.authProvider)
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
	assert.Nil(t, p.backgroundTaskProvider)
	assert.Nil(t, p.metricsProvdider)
	assert.Nil(t, p.authProvider)

}

func TestInitPlatformDuplicateProvider(t *testing.T) {
	reset()

	opt1 := WithProvider("test1", ProviderTypeLogger, newTestProvider)
	opt2 := WithProvider("test2", ProviderTypeLogger, newTestProvider)

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

	opt := WithProvider("test", ProviderTypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	p := DefaultPlatform()
	assert.NotNil(t, p)

	err := p.RegisterProviders(false, opt)
	assert.Error(t, err)

	err = p.RegisterProviders(true, opt)
	assert.NoError(t, err)

}

func TestInitAuthorizationProvider(t *testing.T) {
	reset()

	p := DefaultPlatform()

	assert.NotNil(t, p.authProvider)
	assert.NotNil(t, AuthorizationProvider())
}
