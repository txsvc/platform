package provider

const (
	TypeLogger ProviderType = iota
	TypeErrorReporter
	TypeHttpContext
	TypeTask
	TypeMetrics
	TypeAuthentication
)

type (
	ProviderType int

	InstanceProviderFunc func() interface{}

	ProviderConfig struct {
		ID   string
		Type ProviderType
		Impl InstanceProviderFunc
	}

	GenericProvider interface {
		Close() error
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
