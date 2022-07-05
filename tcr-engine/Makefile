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

define RUN_TESTS =
@mkdir -p _test_results
@gotestsum --format pkgname --junitfile _test_results/output.xml -- $(1) -covermode atomic -coverprofile=_test_results/coverage-all.out -tags=test_helper -v ./...
@grep -v "_test_" _test_results/coverage-all.out >| _test_results/coverage.out
endef

.PHONY: test
test:
	$(call RUN_TESTS,)

.PHONY: test-short
test-short:
	$(call RUN_TESTS,-short)

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
