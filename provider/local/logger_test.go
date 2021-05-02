package local

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/logging"
)

func TestRegisterPlatform(t *testing.T) {
	dl := platform.PlatformOpts{ID: "platform.logger.default", Type: platform.ProviderTypeLogger, Impl: NewDefaultLoggingProvider}
	p, err := platform.InitPlatform(context.Background(), dl)

	if assert.NoError(t, err) {
		assert.NotNil(t, p)

		platform.RegisterPlatform(p)
		logger := platform.Logger("somelogger")

		assert.NotNil(t, logger)
	}
}

func TestDefaultLogger(t *testing.T) {
	InitDefaultProviders()

	logger := platform.Logger("platform-test-logs")
	assert.NotNil(t, logger)

	logger.Log("something happened")
}

func TestDefaultLoggerWithLevel(t *testing.T) {
	InitDefaultProviders()

	logger := platform.Logger("platform-test-logs")
	assert.NotNil(t, logger)

	logger.LogWithLevel(logging.Info, "something happened with level INFO")
	logger.LogWithLevel(logging.Warn, "something happened with level WARN")
	logger.LogWithLevel(logging.Error, "something happened with level ERROR")
	logger.LogWithLevel(logging.Debug, "something happened with level DEBUG")
}

func TestLoggingWithParams(t *testing.T) {
	InitDefaultProviders()

	logger := platform.Logger("platform-test-logs")
	assert.NotNil(t, logger)

	logger.LogWithLevel(logging.Info, "something with parameters happened", "foo", "bar", "question", fmt.Sprintf("%d", 42), "orphan", fmt.Sprintf("%v", true))
}
