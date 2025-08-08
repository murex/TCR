.PHONY: default
default: help ;

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Modules that can be built into an executable
BUILD_MODULES = webapp src
.PHONY: $(BUILD_MODULES)

# Documentation modules
DOC_MODULES = tcr-doc
.PHONY: $(DOC_MODULES)

# Modules with a main package
MAIN_MODULES = $(BUILD_MODULES) $(DOC_MODULES)
# Production modules
PROD_MODULES = $(BUILD_MODULES)

# Module dependencies
src: webapp
tcr-doc: src

.PHONY: prepare
prepare: deps tidy lint build doc test ## Convenience target for automating release preparation

.PHONY: deps tidy vet lint
deps: ## Run deps target on all modules
tidy: ## Run tidy target on all modules
vet: ## Run vet target on all modules
lint: ## Run lint target on all modules
deps tidy vet lint: $(MAIN_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: test test-short cov
test: ## Run tests for all production modules
test-short: ## Run short tests for all production modules
cov: ## Run tests with coverage for all production modules
test test-short cov: $(PROD_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: build download run
build: ## Build all production modules (webapp, src)
download: ## Download dependencies for all production modules
run: ## Run all production modules
build download run: $(BUILD_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: doc
doc: ## Generate command line documentation
doc: $(DOC_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: release
release: ## Create a release using goreleaser
	@goreleaser $(GORELEASER_ARGS)

.PHONY: snapshot
snapshot: ## Create a snapshot release using goreleaser with --snapshot
snapshot: GORELEASER_ARGS= --clean --snapshot
snapshot: release
