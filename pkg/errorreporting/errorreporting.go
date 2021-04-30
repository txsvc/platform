package errorreporting

type (
	ErrorReportingProvider interface {
		ReportError(error)
	}
)
