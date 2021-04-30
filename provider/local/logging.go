package local

import (
	"github.com/txsvc/platform/pkg/logging"
	"go.uber.org/zap"
)

type (
	Logger struct {
		lvl logging.Severity
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
		lvl: logging.Info,
		log: l.Sugar(),
	}

	return &logger
}

func (l *Logger) Log(msg string, keyValuePairs ...interface{}) {
	l.LogWithLevel(l.lvl, msg, keyValuePairs...)
}

func (l *Logger) LogWithLevel(lvl logging.Severity, msg string, keyValuePairs ...interface{}) {
	switch lvl {
	case logging.Info:
		l.log.Infow(msg, keyValuePairs...)
	case logging.Warn:
		l.log.Warnw(msg, keyValuePairs...)
	case logging.Error:
		l.log.Errorw(msg, keyValuePairs...)
	case logging.Debug:
		l.log.Debugw(msg, keyValuePairs...)
	}
}
