.PHONY: test
test:
	@go test ./... -coverprofile coverage.out
	@go tool cover -html=coverage.out

.PHONY: read-test
read-test:
	@go run . read

.PHONY: write-test
write-test:
	@go run . write 1000

.PHONY: concurrency-test
concurrency-test:
	@go run . concurrency 1000
