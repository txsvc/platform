package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/env"
	"github.com/txsvc/platform/pkg/tasks"
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

func TestCloudTasks(t *testing.T) {
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
