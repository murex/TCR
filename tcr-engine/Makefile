.PHONY: default
default: test ;

.PHONY: deps
deps:
	@go get

.PHONY: lint
lint:
	@golangci-lint run --build-tags=test_helper -v

.PHONY: vet
vet:
	@go vet --tags=test_helper ./...

# When running outside of GitHub actions, e.g. locally, we run tests with the -short flag to prevent useless waiting times
ifeq ($(GITHUB_ACTIONS), true)
    SHORT_FLAG := ""
else
    SHORT_FLAG := "-short"
endif

.PHONY: test
test:
	@mkdir -p _test_results
	@gotestsum --format pkgname --junitfile _test_results/output.xml -- ${SHORT_FLAG} -covermode atomic -coverprofile=_test_results/coverage-all.out -tags=test_helper -v ./...
	@grep -v "_test_" _test_results/coverage-all.out >| _test_results/coverage.out

.PHONY: cov
cov:
	@go tool cover -html=_test_results/coverage.out

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -v ./...

.PHONY: download
download:
	@go mod download

.PHONY: install-tools
install-tools: download
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest
