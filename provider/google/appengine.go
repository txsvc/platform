package google

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
)

type (
	GoogleAppEngine struct {
	}
)

func NewAppEngineContextProvider(ID string) interface{} {
	return &GoogleAppEngine{}
}

func (c *GoogleAppEngine) NewHttpContext(req *http.Request) context.Context {
	return appengine.NewContext(req)
}
