package provider

import (
	"context"
)

type (
	MetricsProvider interface {
		Meter(ctx context.Context, metric string, args ...string)
	}
)
