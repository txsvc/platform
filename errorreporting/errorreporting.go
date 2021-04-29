package errorreporting

import (
	"log"

	"go.uber.org/zap"
)

type (
	ErrorReportingProvider interface {
		ReportError(error)
	}

	ErrorReporter struct {
		log *zap.SugaredLogger
	}
)

var errorReporter *ErrorReporter

func init() {
	callerSkipConf := zap.AddCallerSkip(2)
	l, err := zap.NewProduction(callerSkipConf)

	if err != nil {
		log.Fatal(err)
	}
	er := ErrorReporter{
		log: l.Sugar(),
	}

	errorReporter = &er
}

func NewDefaultErrorReporter(ID string) interface{} {
	return errorReporter
}

func (er *ErrorReporter) ReportError(e error) {
	er.log.Error(e)
}
