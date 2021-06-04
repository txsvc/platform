package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotEmpty(t *testing.T) {
	assert.True(t, NotEmpty("a"))
	assert.True(t, NotEmpty("a", "b"))
}

func TestEmpty(t *testing.T) {
	assert.False(t, NotEmpty())
	assert.False(t, NotEmpty(""))
	assert.False(t, NotEmpty("a", "", "c"))
}

func TestIsMemberOf(t *testing.T) {
	assert.True(t, IsMemberOf("ab", "a", "b", "ab", "aa"))

	assert.False(t, IsMemberOf("a"))
	assert.False(t, IsMemberOf("bb", "a", "b", "ab", "aa"))
}
