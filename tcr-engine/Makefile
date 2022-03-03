.PHONY: default
default: test ;

.PHONY: lint
lint:
	@golangci-lint run --build-tags=test_helper -v

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test:
	@RICHGO_FORCE_COLOR=1 richgo test -coverprofile=cover.out -tags=test_helper -v ./...

.PHONY: cov
cov:
	@go tool cover -html=cover.out

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -v ./...

