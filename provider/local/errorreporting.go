package local

import (
	"log"

	"go.uber.org/zap"
)

type (
	ErrorReporter struct {
		log *zap.SugaredLogger
	}
)

var errorReporting *ErrorReporter

func init() {
	callerSkipConf := zap.AddCallerSkip(2)
	l, err := zap.NewProduction(callerSkipConf)

	if err != nil {
		log.Fatal(err)
	}
	er := ErrorReporter{
		log: l.Sugar(),
	}

	errorReporting = &er
}

func NewDefaultErrorReportingProvider(ID string) interface{} {
	return errorReporting
}

func (er *ErrorReporter) ReportError(e error) {
	er.log.Error(e)
}
