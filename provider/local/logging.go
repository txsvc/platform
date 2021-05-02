package local

import (
	"github.com/txsvc/platform/v2/pkg/logging"
	"go.uber.org/zap"
)

type (
	Logger struct {
		lvl logging.Severity
		log *zap.SugaredLogger
	}
)

// the default logger implementation

func NewDefaultLoggingProvider(ID string) interface{} {
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

func (l *Logger) Log(msg string, keyValuePairs ...string) {
	l.LogWithLevel(l.lvl, msg, keyValuePairs...)
}

func (l *Logger) LogWithLevel(lvl logging.Severity, msg string, keyValuePairs ...string) {

	if len(keyValuePairs) > 0 {
		params := make([]interface{}, len(keyValuePairs))
		for i := range keyValuePairs {
			params[i] = keyValuePairs[i]
		}

		switch lvl {
		case logging.Info:
			l.log.Infow(msg, params...)
		case logging.Warn:
			l.log.Warnw(msg, params...)
		case logging.Error:
			l.log.Errorw(msg, params...)
		case logging.Debug:
			l.log.Debugw(msg, params...)
		}
	} else {
		switch lvl {
		case logging.Info:
			l.log.Infow(msg)
		case logging.Warn:
			l.log.Warnw(msg)
		case logging.Error:
			l.log.Errorw(msg)
		case logging.Debug:
			l.log.Debugw(msg)
		}
	}
}
