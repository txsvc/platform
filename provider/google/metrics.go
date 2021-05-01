package google

import (
	"context"
	"fmt"

	"github.com/txsvc/platform/pkg/env"
	"github.com/txsvc/platform/pkg/logging"
)

// the metrics implementation is basically a logger.

func NewGoogleCloudMetricsProvider(ID string) interface{} {
	metrics := fmt.Sprintf("%s-metrics", env.GetString("SERVICE_NAME", "default"))
	return &GoogleCloudLogger{
		logger: client.Logger(metrics),
	}
}

func (l *GoogleCloudLogger) Meter(ctx context.Context, metric string, args ...string) {
	l.LogWithLevel(logging.Info, metric, args...)
}
