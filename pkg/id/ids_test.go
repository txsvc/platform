package id

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	cs1 := Checksum("check me")
	result1 := "10f335e9"
	result2 := "00000000"

	assert.NotEmpty(t, cs1)
	assert.Equal(t, result1, cs1)

	cs2 := Checksum("")
	assert.NotEmpty(t, cs2)
	assert.Equal(t, result2, cs2)
}

func TestFingerprint(t *testing.T) {
	cs1 := Fingerprint("check me")
	result1 := "6e7227eb9fb0793b0e150facda30c40b"
	result2 := "d41d8cd98f00b204e9800998ecf8427e"

	assert.NotEmpty(t, cs1)
	assert.Equal(t, result1, cs1)

	cs2 := Fingerprint("")
	assert.NotEmpty(t, cs2)
	assert.Equal(t, result2, cs2)
}

func TestUUID(t *testing.T) {
	uuid, err := UUID()

	assert.NotEmpty(t, uuid)
	assert.NoError(t, err)

	parts := strings.Split(uuid, "-")
	assert.Equal(t, 5, len(parts))
	assert.Equal(t, 8, len(parts[0]))
	assert.Equal(t, 4, len(parts[1]))
	assert.Equal(t, 4, len(parts[2]))
	assert.Equal(t, 4, len(parts[3]))
	assert.Equal(t, 12, len(parts[4]))
}

func TestRandomToken(t *testing.T) {
	prefix := "xoxo"
	token, err := RandomToken(prefix)

	assert.NotEmpty(t, token)
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(token, prefix))

	parts := strings.Split(token, "-")
	assert.Equal(t, 3, len(parts))
}

func TestShortUUID(t *testing.T) {
	uuid, err := ShortUUID()

	assert.NotEmpty(t, uuid)
	assert.NoError(t, err)
	assert.Equal(t, 12, len(uuid))
}

func TestSimpleUUID(t *testing.T) {
	uuid, err := SimpleUUID()

	assert.NotEmpty(t, uuid)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(uuid))
}
