package local

import (
	"go.uber.org/zap"
)

const (
	Info Severity = iota
	Warn
	Error
	Debug
)

type (
	Severity int

	Logger struct {
		Lvl Severity
		log *zap.SugaredLogger
	}
)

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
