package logging

import (
	"time"

	"go.uber.org/zap"
)

type (
	Severity int

	Entry struct {
		// Timestamp is the time of the entry. If zero, the current time is used.
		Timestamp time.Time

		// Severity is the entry's severity level.
		// The zero value is Default.
		Severity Severity

		// Payload must be either a string, or something that marshals via the
		// encoding/json package to a JSON object (and not any other type of JSON value).
		Payload interface{}

		// Labels optionally specifies key/value labels for the log entry.
		// The Logger.Log method takes ownership of this map. See Logger.CommonLabels
		// for more about labels.
		Labels map[string]string
	}

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

func LogEntry(payload interface{}, params ...interface{}) Entry {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Severity:  Info,
		Payload:   payload,
	}

	if len(params) > 0 {
		p := make(map[string]string)
		n := len(params) / 2
		for i := 0; i < n; i++ {
			key := params[i*2].(string)
			value := params[(i*2)+1].(string)
			p[key] = value
		}
		if len(params)%2 == 1 {
			key := params[len(params)-1].(string)
			p[key] = ""
		}
		e.Labels = p
	}
	return e
}

// the default logger implementation

func NewDefaultLogger(ID string) interface{} {
	l, err := zap.NewProduction()
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

/*
func (l *Logger) Log(e Entry) {
	var entry *lr.Entry
	if len(e.Labels) == 0 {
		entry = l.log.WithFields(nil)
	} else {
		f := make(lr.Fields)
		for k, v := range e.Labels {
			f[k] = v
		}
		entry = l.log.WithFields(f)
	}
	switch e.Severity {
	case Info:
		entry.Info(e.Payload)
	case Warn:
		entry.Warn(e.Payload)
	case Error:
		entry.Error(e.Payload)
	case Debug:
		fmt.Println("DEBUG")
		entry.Debug(e.Payload)
		fmt.Println("BUG")
	}
}

func (l *Logger) LogSync(ctx context.Context, e Entry) error {
	l.Log(e)
	return nil
}

func (l *Logger) Msg(msg string) {
	l.Log(Entry{Payload: msg})
}

func (l *Logger) ReportError(e error) {
	l.log.Error(e.Error())
}

func (l *Logger) Flush() error {
	return nil
}
*/
