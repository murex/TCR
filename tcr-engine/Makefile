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
	@RICHGO_FORCE_COLOR=1 richgo test -covermode atomic -coverprofile=coverage.out -tags=test_helper -v ./...

.PHONY: cov
cov:
	@go tool cover -html=coverage.out

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -v ./...

