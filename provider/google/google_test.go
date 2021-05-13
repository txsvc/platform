package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/pkg/env"
	"github.com/txsvc/platform/v2/tasks"
)

func TestGoogleErrorReporting(t *testing.T) {
	require.True(t, env.Assert("PROJECT_ID"))
	require.True(t, env.Assert("GOOGLE_APPLICATION_CREDENTIALS"))

	p, err := platform.InitPlatform(context.TODO(), GoogleErrorReportingConfig)

	if assert.NoError(t, err) {
		assert.NotNil(t, p)

		platform.RegisterPlatform(p)

		err := fmt.Errorf("something went wrong")
		platform.ReportError(err)
	}
}

func TTestCloudTasks(t *testing.T) {
	require.True(t, env.Assert("PROJECT_ID"))
	require.True(t, env.Assert("GOOGLE_APPLICATION_CREDENTIALS"))

	p, err := platform.InitPlatform(context.TODO(), GoogleCloudTaskConfig)

	if assert.NoError(t, err) {
		assert.NotNil(t, p)

		platform.RegisterPlatform(p)

		task := tasks.HttpTask{
			Method:  tasks.HttpMethodGet,
			Request: "http://podops.dev",
			//Token:   "abc123",
			//Payload: nil,
		}

		err := platform.NewTask(task)
		assert.NoError(t, err)
	}
}

func TestCloudLogging(t *testing.T) {
	require.True(t, env.Assert("PROJECT_ID"))
	require.True(t, env.Assert("GOOGLE_APPLICATION_CREDENTIALS"))

	p, err := platform.InitPlatform(context.TODO(), GoogleCloudLoggingConfig)

	if assert.NoError(t, err) {
		assert.NotNil(t, p)

		platform.RegisterPlatform(p)
		log := platform.Logger("test")

		log.Log("something went OK")
		log.LogWithLevel(logging.Warn, "something with WARN")

		log.Log("something with PARAMS", "foo", "bar", "baz", fmt.Sprintf("%d", 42))
		log.Log("something with 1 PARAM", "foo")
		log.Log("something with odd PARAMS", "foo", "bar", "baz")
	}
}

func TestCloudMetrics(t *testing.T) {
	require.True(t, env.Assert("PROJECT_ID"))
	require.True(t, env.Assert("GOOGLE_APPLICATION_CREDENTIALS"))

	p, err := platform.InitPlatform(context.TODO(), GoogleCloudLoggingConfig, GoogleCloudMetricsConfig)

	if assert.NoError(t, err) {
		assert.NotNil(t, p)

		platform.RegisterPlatform(p)
		platform.Meter(context.Background(), "some.thing", "foo", "bar", "a", "B")
	}
}
