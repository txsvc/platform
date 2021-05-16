package provider

const (
	LevelInfo Severity = iota
	LevelWarn
	LevelError
	LevelDebug
)

type (
	Severity int

	// LoggingProvider defines a generic logging provider
	LoggingProvider interface {
		Log(string, ...string)
		LogWithLevel(Severity, string, ...string)
	}
)
