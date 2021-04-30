package local

import (
	"context"
	"fmt"

	"github.com/txsvc/platform/pkg/tasks"
)

type (
	Tasks struct{}
)

func NewDefaultTaskProvider(ID string) interface{} {
	return &Tasks{}
}

func (t *Tasks) CreateHttpTask(ctx context.Context, task tasks.HttpTask) error {
	fmt.Println(task)
	return nil
}
