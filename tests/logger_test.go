package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform"
	"github.com/txsvc/platform/provider/local"
)

func TestRegisterPlatform(t *testing.T) {
	dl := platform.PlatformOpts{ID: "platform.logger.default", Type: platform.ProviderTypeLogger, Impl: local.NewDefaultLogger}
	p, err := platform.InitPlatform(context.TODO(), dl)

	if assert.NoError(t, err) {
		assert.NotNil(t, p)

		platform.RegisterPlatform(p)
		logger := platform.Logger("somelogger")

		assert.NotNil(t, logger)
	}
}

func TestDefaultLogger(t *testing.T) {
	logger := platform.Logger("platform-test-logs")
	assert.NotNil(t, logger)

	logger.Log("something happened")
}

func TestEntryWithParams(t *testing.T) {
	logger := platform.Logger("platform-test-logs")
	assert.NotNil(t, logger)

	logger.Log("something with parameters happened", "foo", "bar", "question", 42, "orphan", true)
}
