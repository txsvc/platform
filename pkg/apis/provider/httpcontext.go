package provider

import (
	"context"
	"net/http"
)

type (
	HttpContextProvider interface {
		NewHttpContext(*http.Request) context.Context
	}
)
