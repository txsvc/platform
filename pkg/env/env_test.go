package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	e := GetString("PATH", "nix")

	assert.NotEmpty(t, e)
	assert.NotEqual(t, "nix", e)
}

func TestGetStringDefault(t *testing.T) {
	e := GetString("FOO_BAR", "nix")

	assert.NotEmpty(t, e)
	assert.Equal(t, "nix", e)
}

func TestGetIntDefault(t *testing.T) {
	e := GetInt("FOO_BAR", 42)
	assert.Equal(t, int64(42), e)
}

func TestAssert(t *testing.T) {
	assert.True(t, Assert("PATH"))
	assert.False(t, Assert("FOO_BAR"))
}
