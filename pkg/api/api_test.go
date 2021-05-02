package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStatus(t *testing.T) {
	stat := NewStatus(200, "ok")

	assert.NotNil(t, stat)
	assert.Equal(t, 200, stat.Status)
	assert.Equal(t, "ok", stat.Message)
}

func TestNewErrorStatus(t *testing.T) {
	err := fmt.Errorf("error")
	stat := NewErrorStatus(500, err)

	assert.NotNil(t, stat)
	assert.Equal(t, 500, stat.Status)
	assert.Equal(t, "error", stat.Message)
	assert.NotEmpty(t, stat.Error())
}
