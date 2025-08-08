.PHONY: default
default: help ;

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"; \
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: ## Download dependencies
	@go get -u -t tool ./...

.PHONY: lint
lint: ## Run linter
	@golangci-lint run --build-tags=test_helper -v

.PHONY: modernize
modernize: ## Run code modernizer utility
	@go run -tags test_helper golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test -fix ./...

.PHONY: vet
vet: ## Run go vet
	@go vet --tags=test_helper ./...

define RUN_TESTS
@mkdir -p _test_results
@go tool gotestsum --format pkgname --junitfile _test_results/output.xml --hide-summary=skipped -- $(1) -covermode atomic -coverprofile=_test_results/coverage-all.out -tags=test_helper -v ./...
@grep -v "_test_" _test_results/coverage-all.out >| _test_results/coverage.out
@sed -e 's:\(/tcr/\):\1src/:g' _test_results/coverage.out >| _test_results/coverage-sonar.out
endef

.PHONY: test
test: ## Run all tests
	$(call RUN_TESTS,)

.PHONY: test-short
test-short: ## Run short tests
	$(call RUN_TESTS,-short)

.PHONY: slow-tests
slow-tests: ## Run slowest tests report
	@go test -json -tags=test_helper ./... | go tool gotestsum tool slowest --threshold 500ms

.PHONY: cov
cov: ## Run tests with coverage report
	@go tool cover -html=_test_results/coverage.out

.PHONY: tidy
tidy: ## Run go mod tidy
	@go mod tidy

CONFIG_PKG="github.com/murex/tcr/settings"
export CGO_ENABLED=0
.PHONY: build
build: ## Build the tcr binary
	@go build -ldflags "-s -w \
    		-X ${CONFIG_PKG}.BuildVersion="`git describe --tags`" \
            -X ${CONFIG_PKG}.BuildOs="`go env GOOS`" \
            -X ${CONFIG_PKG}.BuildArch="`go env GOARCH`" \
            -X ${CONFIG_PKG}.BuildCommit="`git rev-list --max-count=1 --tags`" \
            -X ${CONFIG_PKG}.BuildDate="`date -u +%FT%TZ`" \
            -X ${CONFIG_PKG}.BuildAuthor="`id -un`" \
            "

.PHONY: run
run: ## Run locally built tcr on testdata/java
	@tcr-local

.PHONY: download
download: ## Download Go module dependencies
	@go mod download
