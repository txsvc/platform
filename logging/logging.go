package logging

import (
	"go.uber.org/zap"
)

type (
	Severity int

	LoggingProvider interface {
		Log(string, ...interface{})
	}

	Logger struct {
		Lvl Severity
		log *zap.SugaredLogger
	}
)

const (
	Info Severity = iota
	Warn
	Error
	Debug
)

// Returns the name of a Severity level
func (l Severity) String() string {
	switch l {
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	case Debug:
		return "debug"
	default:
		panic("unsupported level")
	}
}

// the default logger implementation

func NewDefaultLogger(ID string) interface{} {
	callerSkipConf := zap.AddCallerSkip(1)

	l, err := zap.NewProduction(callerSkipConf)
	if err != nil {
		return nil
	}

	logger := Logger{
		Lvl: Info,
		log: l.Sugar(),
	}

	return &logger
}

func (l *Logger) Log(msg string, keyValuePairs ...interface{}) {
	switch l.Lvl {
	case Info:
		l.log.Infow(msg, keyValuePairs...)
	case Warn:
		l.log.Warnw(msg, keyValuePairs...)
	case Error:
		l.log.Errorw(msg, keyValuePairs...)
	case Debug:
		l.log.Debugw(msg, keyValuePairs...)
	}
}
