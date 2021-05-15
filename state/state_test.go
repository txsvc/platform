package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKey(t *testing.T) {

	k1 := NewKey("thekind", "thekey_1")
	assert.NotNil(t, k1)
	key1 := k1.(Key)
	assert.NotNil(t, key1)

	assert.NotEmpty(t, key1.Name())
	assert.NotEmpty(t, key1.Kind())
	assert.NotEmpty(t, key1.String())
	assert.NotEmpty(t, key1.Encode())

	k2 := NewKey("thekind", "thekey_2")
	assert.NotNil(t, k2)
	key2 := k2.(Key)
	assert.NotNil(t, key2)

	assert.True(t, key1.Equal(key1))
	assert.False(t, key1.Equal(key2))
}
