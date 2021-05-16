package provider

type (
	ErrorReportingProvider interface {
		ReportError(error)
	}
)
