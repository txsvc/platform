package logging

const (
	Info Severity = iota
	Warn
	Error
	Debug
)

type (
	Severity int

	// LoggingProvider defines a generic logging provider
	LoggingProvider interface {
		Log(string, ...string)
		LogWithLevel(Severity, string, ...string)
	}
)
