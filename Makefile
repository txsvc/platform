.PHONY: all
all: test

.PHONY: test
test:
	go test `go list ./... | grep -v 'tests\|google'` -coverprofile=coverage.txt -covermode=atomic
