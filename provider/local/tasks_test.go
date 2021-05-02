package local

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/tasks"
)

func TestDefaultTasks(t *testing.T) {
	InitDefaultProviders()

	task := tasks.HttpTask{
		Method:  tasks.HttpMethodGet,
		Request: "http://get.some.stuff",
		Token:   "abc123",
		Payload: nil,
	}
	err := platform.NewTask(task)
	assert.Error(t, err)
}
