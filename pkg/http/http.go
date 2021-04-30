package http

import (
	"context"
	"net/http"
)

type (
	HttpRequestContextProvider interface {
		NewHttpContext(*http.Request) context.Context
	}
)
