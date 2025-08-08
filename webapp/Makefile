.PHONY: default
default: help ;

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"; \
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Install dependencies
	@npm install

.PHONY: deps
deps: ## Update dependencies
	@npm update

.PHONY: lint
lint: ## Run linter
	@npm run lint

.PHONY: test
test: ## Run tests
	@npm test

.PHONY: cov
cov: ## Run tests with coverage
	@npm run coverage

.PHONY: tidy
tidy: ## Reinstall dependencies
	@npm install

.PHONY: build
build: setup ## Build for production
	@npm run build-prod

.PHONY: run
run: ## Start the application
	@npm start

.PHONY: clean
clean: ## Clean build artifacts
	@npm run clean
