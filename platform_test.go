package platform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txsvc/platform/logging"
)

func TestInitDefaultPlatform(t *testing.T) {
	assert.NotNil(t, platform)
	assert.NotNil(t, platform.providers)
	assert.NotNil(t, platform.logger)
	assert.Equal(t, 2, len(platform.providers))
}

func TestInitPlatformNoProviders(t *testing.T) {
	p, err := InitPlatform(context.TODO())
	if assert.NoError(t, err) {
		assert.NotNil(t, p)
		assert.NotNil(t, p.providers)
		assert.NotNil(t, p.logger)
		assert.Equal(t, 0, len(p.providers))
	}
}

func TestRegisterPlatformNoProviders(t *testing.T) {
	p, err := InitPlatform(context.TODO())
	if assert.NoError(t, err) {
		assert.NotNil(t, p)
		assert.NotNil(t, p.providers)
		assert.NotNil(t, p.logger)
		assert.Equal(t, 0, len(p.providers))

		logger := Logger("somelogger")
		assert.NotNil(t, logger)

		old := RegisterPlatform(nil)
		assert.Nil(t, old)

		old = RegisterPlatform(p)
		assert.NotNil(t, old)

		logger = Logger("somelogger")
		assert.Nil(t, logger)
	}
}

func TestInitPlatform(t *testing.T) {
	dl := PlatformOpts{ID: "platform.logger.default", Type: ProviderTypeLogger, Impl: logging.NewDefaultLogger}

	p, err := InitPlatform(context.TODO(), dl)

	if assert.NoError(t, err) {
		assert.NotNil(t, p)
		assert.NotNil(t, p.providers)
		assert.NotNil(t, p.logger)
		assert.Equal(t, 1, len(p.providers))

		RegisterPlatform(p)

		logger := Logger("somelogger")
		assert.NotNil(t, logger)

		logger.Log("something went OK")
	}
}
