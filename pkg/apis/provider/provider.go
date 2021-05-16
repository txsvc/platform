package provider

import (
	"context"
	"net/http"
)

const (
	TypeLogger ProviderType = iota
	TypeErrorReporter
	TypeHttpContext
	TypeTask
	TypeMetrics
	TypeAuthentication

	LevelInfo Severity = iota
	LevelWarn
	LevelError
	LevelDebug

	HttpMethodGet HttpMethod = iota
	HttpMethodPost
	HttpMethodPut
	HttpMethodDelete
)

type (
	ProviderType int

	Severity int

	InstanceProviderFunc func() interface{}

	ProviderConfig struct {
		ID   string
		Type ProviderType
		Impl InstanceProviderFunc
	}

	HttpMethod int

	HttpTask struct {
		Method  HttpMethod
		Request string
		Token   string
		Payload interface{}
	}

	GenericProvider interface {
		Close() error
	}

	// LoggingProvider defines a generic logging provider
	LoggingProvider interface {
		Log(string, ...string)
		LogWithLevel(Severity, string, ...string)
	}

	ErrorReportingProvider interface {
		ReportError(error)
	}

	MetricsProvider interface {
		Meter(ctx context.Context, metric string, args ...string)
	}

	HttpContextProvider interface {
		NewHttpContext(*http.Request) context.Context
	}

	HttpTaskProvider interface {
		CreateHttpTask(context.Context, HttpTask) error
	}
)

// Returns the name of a provider type
func (l ProviderType) String() string {
	switch l {
	case TypeLogger:
		return "LOGGER"
	case TypeErrorReporter:
		return "ERROR_REPORTER"
	case TypeHttpContext:
		return "HTTP_CONTEXT"
	case TypeTask:
		return "TASK"
	case TypeMetrics:
		return "METRICS"
	case TypeAuthentication:
		return "AUTHENTICATION"
	default:
		panic("unsupported")
	}
}

// WithProvider returns a populated ProviderConfig struct.
func WithProvider(ID string, providerType ProviderType, impl InstanceProviderFunc) ProviderConfig {
	return ProviderConfig{
		ID:   ID,
		Type: providerType,
		Impl: impl,
	}
}
