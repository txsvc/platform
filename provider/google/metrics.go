package google

import (
	"context"

	"github.com/txsvc/platform/v2/pkg/env"
	"github.com/txsvc/platform/v2/pkg/logging"
)

// the metrics implementation is basically a logger.

// see https://pkg.go.dev/go.opentelemetry.io/otel/metric for inspiration

func NewGoogleCloudMetricsProvider(ID string) interface{} {
	metrics := env.GetString("METRICS_LOG_NAME", "metrics")
	return &GoogleCloudLogger{
		logger: client.Logger(metrics),
	}
}

func (l *GoogleCloudLogger) Meter(ctx context.Context, metric string, args ...string) {
	l.LogWithLevel(logging.Info, metric, args...)
}
