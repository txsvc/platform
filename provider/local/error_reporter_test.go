package local

import (
	"fmt"
	"testing"

	"github.com/txsvc/platform/v2"
)

func TestErrorReporter(t *testing.T) {
	InitDefaultProviders()

	err := fmt.Errorf("something went wrong")

	platform.ReportError(err)
}
