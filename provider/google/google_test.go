package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/env"
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
