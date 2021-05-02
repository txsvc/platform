package google

import (
	"context"
	"log"

	"cloud.google.com/go/errorreporting"
	"github.com/txsvc/platform/v2/pkg/env"
)

type (
	ErrorReporter struct {
		client *errorreporting.Client
	}
)

func NewErrorReporter(ID string) interface{} {
	projectID := env.GetString("PROJECT_ID", "")
	serviceName := env.GetString("SERVICE_NAME", "default")

	// initialize error reporting
	ec, err := errorreporting.NewClient(context.Background(), projectID, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("could not log error: %v", err)
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return &ErrorReporter{
		client: ec,
	}
}

func (er *ErrorReporter) ReportError(e error) {
	er.client.Report(errorreporting.Entry{Error: e})
}
