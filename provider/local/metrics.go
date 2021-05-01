package local

import "context"

type (
	MeterImpl struct {
		Metric string
	}
)

// the default logger implementation

func NewDefaultMetricsProvider(ID string) interface{} {
	return &MeterImpl{
		Metric: ID,
	}
}

func (m *MeterImpl) Meter(ctx context.Context, metric string, args ...string) {
	// actually does nothing right now
}
