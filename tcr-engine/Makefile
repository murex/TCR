.PHONY: default
default: test ;

.PHONY: lint
lint:
	@golangci-lint run --build-tags=test_helper -v

.PHONY: vet
vet:
	@go vet --tags=test_helper ./...

.PHONY: test
test:
	@gotestsum --format pkgname --junitfile test-results/output.xml -- -covermode atomic -coverprofile=coverage-all.out -tags=test_helper -v ./...
	@grep -v "_test_" coverage-all.out >| coverage.out

.PHONY: cov
cov:
	@go tool cover -html=coverage.out

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -v ./...

