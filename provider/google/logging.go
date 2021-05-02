package google

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
	"github.com/txsvc/platform/v2/pkg/env"
	lp "github.com/txsvc/platform/v2/pkg/logging"
)

type (
	GoogleCloudLogger struct {
		logger *logging.Logger
	}

	/*
		logEvent struct {
			logger *logging.Logger
			evt    *logging.Entry
		}
	*/
)

var (
	client *logging.Client
	//logEvents chan (*logEvent)
)

func init() {
	projectID := env.GetString("PROJECT_ID", "")

	// initialize logging
	lc, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatal(err)
	}
	client = lc

	// initialize the logger queue
	//logEvents = make(chan *logEvent, 20) // with a backlog
	//go uploader()
}

/*
func uploader() {
	for {
		e := <-logEvents
		e.logger.Log(*e.evt)
	}
}
*/

func NewGoogleCloudLoggingProvider(ID string) interface{} {
	return &GoogleCloudLogger{
		logger: client.Logger(ID),
	}
}

func (l *GoogleCloudLogger) Log(msg string, keyValuePairs ...string) {
	l.LogWithLevel(lp.Info, msg, keyValuePairs...)
}

func (l *GoogleCloudLogger) LogWithLevel(lvl lp.Severity, msg string, keyValuePairs ...string) {
	e := logging.Entry{
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

	/*
		evt := &logEvent{
			logger: l.logger,
			evt:    &e,
		}
		logEvents <- evt
	*/

	l.logger.Log(e)
}

func toSeverity(severity lp.Severity) logging.Severity {
	switch severity {
	case lp.Info:
		return logging.Info
	case lp.Warn:
		return logging.Warning
	case lp.Error:
		return logging.Error
	case lp.Debug:
		return logging.Debug
	}
	return logging.Info
}
