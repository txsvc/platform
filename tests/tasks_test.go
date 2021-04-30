package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/tasks"
)

func TestDefaultTasks(t *testing.T) {
	platform.InitDefaultProviders()

	task := tasks.HttpTask{
		Method:  tasks.HttpMethodGet,
		Request: "http://get.some.stuff",
		Token:   "abc123",
		Payload: nil,
	}
	err := platform.NewTask(task)
	assert.NoError(t, err)
}
