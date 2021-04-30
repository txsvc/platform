package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform"
)

func TestDefaultContext(t *testing.T) {
	platform.InitDefaultProviders()

	ctx := platform.NewHttpContext(nil)
	assert.NotNil(t, ctx)

}
