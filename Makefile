.PHONY: all
all: test

.PHONY: test
test:
	cd authentication && go test
	cd pkg/account && go test
	cd pkg/api && go test
	cd pkg/datastore && go test
	cd pkg/env && go test
	cd pkg/id && go test
	cd pkg/loader && go test
	cd pkg/netrc && go test
	cd pkg/timestamp && go test
	cd pkg/validate && go test
	cd provider/local && go test
	cd provider/google && go test
	go test

.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v 'tests\|google\|httpserver'` -coverprofile=coverage.txt -covermode=atomic
