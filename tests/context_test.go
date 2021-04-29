package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform"
)

func TestDefaultContext(t *testing.T) {
	ctx := platform.NewHttpContext(nil)
	assert.NotNil(t, ctx)

}
