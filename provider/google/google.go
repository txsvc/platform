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
	"github.com/txsvc/platform/v2/pkg/apis/provider"
	"github.com/txsvc/platform/v2/pkg/env"
	//"github.com/txsvc/platform/v2/tasks"
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
	GoogleErrorReportingConfig provider.ProviderConfig = provider.WithProvider("platform.google.errorreporting", provider.TypeErrorReporter, NewStackdriverErrorReportingProvider)
	GoogleCloudTaskConfig      provider.ProviderConfig = provider.WithProvider("platform.google.task", provider.TypeTask, NewCloudTasksProvider)
	GoogleCloudLoggingConfig   provider.ProviderConfig = provider.WithProvider("platform.google.logger", provider.TypeLogger, NewStackdriverLoggingProvider)
	GoogleCloudMetricsConfig   provider.ProviderConfig = provider.WithProvider("platform.google.metrics", provider.TypeMetrics, NewStackdriverMetricsProvider)
	// AppEngine
	AppEngineContextConfig provider.ProviderConfig = provider.WithProvider("platform.google.context", provider.TypeHttpContext, NewAppEngineContextProvider)

	client *stackdriver_logging.Client

	// UserAgentString identifies any http request podops makes
	userAgentString string = "txsvc/platform 1.0.0"
	// workerQueue is the main worker queue for all the background tasks
	workerQueue string = fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetString("PROJECT_ID", ""), env.GetString("LOCATION_ID", ""), env.GetString("DEFAULT_QUEUE", ""))

	// Interface guards
	_ provider.GenericProvider     = (*AppEngineContextImpl)(nil)
	_ provider.HttpContextProvider = (*AppEngineContextImpl)(nil)

	_ provider.GenericProvider        = (*GoogleErrorReportingProviderImpl)(nil)
	_ provider.ErrorReportingProvider = (*GoogleErrorReportingProviderImpl)(nil)

	_ provider.GenericProvider = (*StackdriverLoggingProviderImpl)(nil)
	_ provider.LoggingProvider = (*StackdriverLoggingProviderImpl)(nil)
	_ provider.MetricsProvider = (*StackdriverLoggingProviderImpl)(nil)

	_ provider.GenericProvider  = (*CloudTaskProviderImpl)(nil)
	_ provider.HttpTaskProvider = (*CloudTaskProviderImpl)(nil)
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

func NewAppEngineContextProvider() interface{} {
	return &AppEngineContextImpl{}
}

func (c *AppEngineContextImpl) NewHttpContext(req *h.Request) context.Context {
	return appengine.NewContext(req)
}

func (c *AppEngineContextImpl) Close() error {
	return nil
}

func NewStackdriverErrorReportingProvider() interface{} {
	projectID := env.GetString("PROJECT_ID", "")
	serviceName := env.GetString("SERVICE_NAME", "default")

	// initialize error reporting
	ec, err := stackdriver_error.NewClient(context.Background(), projectID, stackdriver_error.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("could not log error: %v", err)
		},
	})
	if err != nil || ec == nil {
		log.Fatal(err)
	}

	return &GoogleErrorReportingProviderImpl{
		client: ec,
	}
}

func (er *GoogleErrorReportingProviderImpl) ReportError(e error) {
	if e != nil {
		er.client.Report(stackdriver_error.Entry{Error: e})
	}
}

func (c *GoogleErrorReportingProviderImpl) Close() error {
	return nil
}

func NewStackdriverLoggingProvider() interface{} {
	return &StackdriverLoggingProviderImpl{
		logger: client.Logger("default"), // FIXME what should this be?
	}
}

func (c *StackdriverLoggingProviderImpl) Close() error {
	return nil
}

func (l *StackdriverLoggingProviderImpl) Log(msg string, keyValuePairs ...string) {
	l.LogWithLevel(provider.LevelInfo, msg, keyValuePairs...)
}

func (l *StackdriverLoggingProviderImpl) LogWithLevel(lvl provider.Severity, msg string, keyValuePairs ...string) {
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

func NewStackdriverMetricsProvider() interface{} {
	metrics := env.GetString("METRICS_LOG_NAME", "metrics")
	return &StackdriverLoggingProviderImpl{
		logger: client.Logger(metrics),
	}
}

func (l *StackdriverLoggingProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
	l.LogWithLevel(provider.LevelInfo, metric, args...)
}

func NewCloudTasksProvider() interface{} {
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

func (t *CloudTaskProviderImpl) CreateHttpTask(ctx context.Context, task provider.HttpTask) error {

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

func toHttpMethod(m provider.HttpMethod) taskspb.HttpMethod {
	switch m {
	case provider.HttpMethodGet:
		return taskspb.HttpMethod_GET
	case provider.HttpMethodPost:
		return taskspb.HttpMethod_POST
	case provider.HttpMethodPut:
		return taskspb.HttpMethod_PUT
	case provider.HttpMethodDelete:
		return taskspb.HttpMethod_DELETE
	}
	return taskspb.HttpMethod_GET
}

func toSeverity(severity provider.Severity) stackdriver_logging.Severity {
	switch severity {
	case provider.LevelInfo:
		return stackdriver_logging.Info
	case provider.LevelWarn:
		return stackdriver_logging.Warning
	case provider.LevelError:
		return stackdriver_logging.Error
	case provider.LevelDebug:
		return stackdriver_logging.Debug
	}
	return stackdriver_logging.Info
}
