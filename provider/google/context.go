package google

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
)

type (
	GAE struct {
	}
)

func NewAppEngineContextProvider(ID string) interface{} {
	return &GAE{}
}

func (c *GAE) NewHttpContext(req *http.Request) context.Context {
	return appengine.NewContext(req)
}
