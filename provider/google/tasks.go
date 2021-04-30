package google

import (
	"context"
	"encoding/json"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/txsvc/platform/pkg/env"
	"github.com/txsvc/platform/pkg/tasks"
)

type (
	CloudTasks struct {
		client *cloudtasks.Client
	}
)

var (
	// UserAgentString identifies any http request podops makes
	userAgentString string = "txsvc/platform 1.0.0"
	// workerQueue is the main worker queue for all the background tasks
	workerQueue string = fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetString("PROJECT_ID", ""), env.GetString("LOCATION_ID", ""), env.GetString("DEFAULT_QUEUE", ""))
)

func NewCloudTaskProvider(ID string) interface{} {
	client, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		return nil
	}

	return &CloudTasks{
		client: client,
	}
}

func (t *CloudTasks) CreateHttpTask(ctx context.Context, task tasks.HttpTask) error {

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
