TESTING_COVERAGE_FILE=testing_coverage.out

run:
	go run main.go

test:
	go test -cover -coverprofile=${TESTING_COVERAGE_FILE} ./...

test-report:
	go tool cover -html=${TESTING_COVERAGE_FILE}
