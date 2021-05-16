package provider

import (
	"context"
)

const (
	HttpMethodGet HttpMethod = iota
	HttpMethodPost
	HttpMethodPut
	HttpMethodDelete
)

type (
	HttpMethod int

	HttpTask struct {
		Method  HttpMethod
		Request string
		Token   string
		Payload interface{}
	}

	HttpTaskProvider interface {
		CreateHttpTask(context.Context, HttpTask) error
	}
)
