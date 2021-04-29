package tests

import (
	"fmt"
	"testing"

	"github.com/txsvc/platform"
)

func TestErrorReporter(t *testing.T) {
	err := fmt.Errorf("something went wrong")

	platform.ReportError(err)
}
