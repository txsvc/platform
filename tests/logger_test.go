package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/platform"
)

func TestDefaultLogger(t *testing.T) {
	lg := platform.Logger("platform-test-logs")
	assert.NotNil(t, lg)

	lg.Log("something happened")
}

func TestEntryWithParams(t *testing.T) {
	lg := platform.Logger("platform-test-logs")
	assert.NotNil(t, lg)

	lg.Log("something with parameters", "foo", "bar", "question", 42, "orphan", true)
}
