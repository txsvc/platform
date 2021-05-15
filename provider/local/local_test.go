package local

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/tasks"
)

func TestRegisterPlatform(t *testing.T) {
	dl := platform.PlatformOpts{ID: "platform.logger.default", Type: platform.ProviderTypeLogger, Impl: NewLocalLoggingProvider}
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

func TestDefaultLogger(t *testing.T) {
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

func TestErrorReporter(t *testing.T) {
	InitLocalProviders()

	err := fmt.Errorf("something went wrong")

	platform.ReportError(err)
}

func TestDefaultTasks(t *testing.T) {
	InitLocalProviders()

	task := tasks.HttpTask{
		Method:  tasks.HttpMethodGet,
		Request: "http://get.some.stuff",
		Token:   "abc123",
		Payload: nil,
	}
	err := platform.NewTask(task)
	assert.Error(t, err)
}

func TestDefaultAuthorizationProvider(t *testing.T) {
	InitLocalProviders()

	ap := platform.AuthenticationProvider()

	assert.NotNil(t, ap)

	opts := ap.Options()

	assert.NotNil(t, opts)
	assert.NotEqual(t, 0, opts.AuthenticationExpiration)
	assert.NotEqual(t, 0, opts.AuthorizationExpiration)
	assert.NotEmpty(t, opts.Endpoint)
	assert.NotEmpty(t, opts.Scope)
}
