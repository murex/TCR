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

.PHONY: doc
doc: tidy ## Generate cobra documentation
	@echo "Generating cobra documentation"
	@go run .

.PHONY: tidy
tidy: ## Run go mod tidy
	@go mod tidy
