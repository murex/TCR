.PHONY: default
default: build doc ;

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

# Convenience target for automating release preparation
.PHONY: prepare
prepare: deps tidy lint build doc test

.PHONY: deps tidy vet lint
deps tidy vet lint: $(MAIN_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: test test-short cov
test test-short cov: $(PROD_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: build download run
build download run: $(BUILD_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: doc
doc: $(DOC_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: release
release:
	@goreleaser $(GORELEASER_ARGS)

.PHONY: snapshot
snapshot: GORELEASER_ARGS= --clean --snapshot
snapshot: release
