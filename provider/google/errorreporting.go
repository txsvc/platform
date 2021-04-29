package google

import (
	"context"
	"log"

	"cloud.google.com/go/errorreporting"
	"github.com/txsvc/platform/pkg/env"
)

type (
	ErrorReporter struct {
		errorClient *errorreporting.Client
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
		errorClient: ec,
	}
}

func (er *ErrorReporter) ReportError(e error) {
	er.errorClient.Report(errorreporting.Entry{Error: e})
}
