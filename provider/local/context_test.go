package local

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform/v2"
)

func TestDefaultContext(t *testing.T) {
	InitDefaultProviders()

	ctx := platform.NewHttpContext(nil)
	assert.NotNil(t, ctx)
}
