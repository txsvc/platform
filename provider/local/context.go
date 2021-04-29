package local

import (
	"context"
	"net/http"
)

type (
	DefaultContext struct {
	}
)

func NewDefaultContextProvider(ID string) interface{} {
	return &DefaultContext{}
}

func (c *DefaultContext) NewHttpContext(req *http.Request) context.Context {
	return context.Background()
}
