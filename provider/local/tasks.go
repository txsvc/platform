package local

import (
	"context"
	"fmt"

	"github.com/txsvc/platform/v2/pkg/tasks"
)

type (
	Tasks struct{}
)

func NewDefaultTaskProvider(ID string) interface{} {
	return &Tasks{}
}

func (t *Tasks) CreateHttpTask(ctx context.Context, task tasks.HttpTask) error {
	return fmt.Errorf("not implemented")
}
