.PHONY: default
default: build doc ;

# Modules that can be built into an executable
BUILD_MODULES = src
.PHONY: $(BUILD_MODULES)

# Documentation modules
DOC_MODULES = tcr-doc
.PHONY: $(DOC_MODULES)

# Modules with a main package
MAIN_MODULES = $(BUILD_MODULES) $(DOC_MODULES)
# Production modules
PROD_MODULES = $(BUILD_MODULES)

# Module dependencies
tcr-doc: src

# Convenience target for automating release preparation
.PHONY: prepare
prepare: deps install-tools tidy lint build doc test

.PHONY: deps
deps: $(MAIN_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: tidy vet lint
tidy vet lint: $(MAIN_MODULES)
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

.PHONY: build install-tools download run
build install-tools download run: $(BUILD_MODULES)
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
snapshot: GORELEASER_ARGS= --rm-dist --snapshot
snapshot: release
