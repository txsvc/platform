package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	h "net/http"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	stackdriver_error "cloud.google.com/go/errorreporting"
	stackdriver_logging "cloud.google.com/go/logging"
	"google.golang.org/appengine"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/errorreporting"
	"github.com/txsvc/platform/v2/http"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/metrics"
	"github.com/txsvc/platform/v2/pkg/env"
	"github.com/txsvc/platform/v2/tasks"
)

type (
	AppEngineContextImpl struct {
	}

	GoogleErrorReportingProviderImpl struct {
		client *stackdriver_error.Client
	}

	StackdriverLoggingProviderImpl struct {
		logger *stackdriver_logging.Logger
	}

	CloudTaskProviderImpl struct {
		client *cloudtasks.Client
	}
)

var (
	// Google Cloud Platform
	GoogleErrorReportingConfig platform.PlatformOpts = platform.WithProvider("platform.google.errorreporting", platform.ProviderTypeErrorReporter, NewStackdriverErrorReportingProvider)
	GoogleCloudTaskConfig      platform.PlatformOpts = platform.WithProvider("platform.google.task", platform.ProviderTypeTask, NewCloudTasksProvider)
	GoogleCloudLoggingConfig   platform.PlatformOpts = platform.WithProvider("platform.google.logger", platform.ProviderTypeLogger, NewStackdriverLoggingProvider)
	GoogleCloudMetricsConfig   platform.PlatformOpts = platform.WithProvider("platform.google.metrics", platform.ProviderTypeMetrics, NewStackdriverMetricsProvider)
	// AppEngine
	AppEngineContextConfig platform.PlatformOpts = platform.WithProvider("platform.google.context", platform.ProviderTypeHttpContext, NewAppEngineContextProvider)

	client *stackdriver_logging.Client

	// UserAgentString identifies any http request podops makes
	userAgentString string = "txsvc/platform 1.0.0"
	// workerQueue is the main worker queue for all the background tasks
	workerQueue string = fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetString("PROJECT_ID", ""), env.GetString("LOCATION_ID", ""), env.GetString("DEFAULT_QUEUE", ""))

	// Interface guards
	_ platform.GenericProvider        = (*AppEngineContextImpl)(nil)
	_ http.HttpRequestContextProvider = (*AppEngineContextImpl)(nil)

	_ platform.GenericProvider              = (*GoogleErrorReportingProviderImpl)(nil)
	_ errorreporting.ErrorReportingProvider = (*GoogleErrorReportingProviderImpl)(nil)

	_ platform.GenericProvider = (*StackdriverLoggingProviderImpl)(nil)
	_ logging.LoggingProvider  = (*StackdriverLoggingProviderImpl)(nil)
	_ metrics.MetricsProvider  = (*StackdriverLoggingProviderImpl)(nil)

	_ platform.GenericProvider = (*CloudTaskProviderImpl)(nil)
	_ tasks.HttpTaskProvider   = (*CloudTaskProviderImpl)(nil)
)

func init() {
	projectID := env.GetString("PROJECT_ID", "")

	// initialize logging
	lc, err := stackdriver_logging.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatal(err)
	}
	client = lc

}

func InitGoogleCloudPlatformProviders() {
	p, err := platform.InitPlatform(context.Background(), GoogleErrorReportingConfig, GoogleCloudTaskConfig, GoogleCloudLoggingConfig, GoogleCloudMetricsConfig, AppEngineContextConfig)
	if err != nil {
		log.Fatal(err)
	}
	platform.RegisterPlatform(p)
}

func NewAppEngineContextProvider(ID string) interface{} {
	return &AppEngineContextImpl{}
}

func (c *AppEngineContextImpl) NewHttpContext(req *h.Request) context.Context {
	return appengine.NewContext(req)
}

func (c *AppEngineContextImpl) Close() error {
	return nil
}

func NewStackdriverErrorReportingProvider(ID string) interface{} {
	projectID := env.GetString("PROJECT_ID", "")
	serviceName := env.GetString("SERVICE_NAME", "default")

	// initialize error reporting
	ec, err := stackdriver_error.NewClient(context.Background(), projectID, stackdriver_error.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("could not log error: %v", err)
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return &GoogleErrorReportingProviderImpl{
		client: ec,
	}
}

func (er *GoogleErrorReportingProviderImpl) ReportError(e error) {
	er.client.Report(stackdriver_error.Entry{Error: e})
}

func (c *GoogleErrorReportingProviderImpl) Close() error {
	return nil
}

func NewStackdriverLoggingProvider(ID string) interface{} {
	return &StackdriverLoggingProviderImpl{
		logger: client.Logger(ID),
	}
}

func (c *StackdriverLoggingProviderImpl) Close() error {
	return nil
}

func (l *StackdriverLoggingProviderImpl) Log(msg string, keyValuePairs ...string) {
	l.LogWithLevel(logging.Info, msg, keyValuePairs...)
}

func (l *StackdriverLoggingProviderImpl) LogWithLevel(lvl logging.Severity, msg string, keyValuePairs ...string) {
	e := stackdriver_logging.Entry{
		Payload:  msg,
		Severity: toSeverity(lvl),
	}

	n := len(keyValuePairs)
	if n > 0 {
		labels := make(map[string]string)
		if n == 1 {
			labels[keyValuePairs[0]] = ""
		} else {
			for i := 0; i < n/2; i++ {
				k := keyValuePairs[i*2]
				v := keyValuePairs[(i*2)+1]
				labels[k] = v
			}
			if n%2 == 1 {
				labels[keyValuePairs[n-1]] = ""
			}
		}
		e.Labels = labels
	}

	l.logger.Log(e)
}

// the metrics implementation is basically a logger.

// see https://pkg.go.dev/go.opentelemetry.io/otel/metric for inspiration

func NewStackdriverMetricsProvider(ID string) interface{} {
	metrics := env.GetString("METRICS_LOG_NAME", "metrics")
	return &StackdriverLoggingProviderImpl{
		logger: client.Logger(metrics),
	}
}

func (l *StackdriverLoggingProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
	l.LogWithLevel(logging.Info, metric, args...)
}

func NewCloudTasksProvider(ID string) interface{} {
	client, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		return nil
	}

	return &CloudTaskProviderImpl{
		client: client,
	}
}

func (c *CloudTaskProviderImpl) Close() error {
	return nil
}

func (t *CloudTaskProviderImpl) CreateHttpTask(ctx context.Context, task tasks.HttpTask) error {

	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   userAgentString,
	}
	if task.Token != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", task.Token)
	}

	req := &taskspb.CreateTaskRequest{
		Parent: workerQueue,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: toHttpMethod(task.Method),
					Url:        task.Request,
					Headers:    headers,
				},
			},
		},
	}

	if task.Payload != nil {
		// marshal the payload
		b, err := json.Marshal(task.Payload)
		if err != nil {
			return err
		}
		req.Task.GetHttpRequest().Body = b
	}

	_, err := t.client.CreateTask(ctx, req)
	return err
}

func toHttpMethod(m tasks.HttpMethod) taskspb.HttpMethod {
	switch m {
	case tasks.HttpMethodGet:
		return taskspb.HttpMethod_GET
	case tasks.HttpMethodPost:
		return taskspb.HttpMethod_POST
	case tasks.HttpMethodPut:
		return taskspb.HttpMethod_PUT
	case tasks.HttpMethodDelete:
		return taskspb.HttpMethod_DELETE
	}
	return taskspb.HttpMethod_GET
}

func toSeverity(severity logging.Severity) stackdriver_logging.Severity {
	switch severity {
	case logging.Info:
		return stackdriver_logging.Info
	case logging.Warn:
		return stackdriver_logging.Warning
	case logging.Error:
		return stackdriver_logging.Error
	case logging.Debug:
		return stackdriver_logging.Debug
	}
	return stackdriver_logging.Info
}
