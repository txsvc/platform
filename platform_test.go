package platform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
