.PHONY: default
default: build doc ;

# Modules that can be build into an executable
BUILD_MODULES = tcr-cli
# Library modules
LIB_MODULES = tcr-engine
# Documentation modules
DOC_MODULES = tcr-doc

# Modules with a main package
MAIN_MODULES = $(BUILD_MODULES) $(DOC_MODULES)
# Production modules
PROD_MODULES = $(BUILD_MODULES) $(LIB_MODULES)

ALL_MODULES = $(BUILD_MODULES) $(LIB_MODULES) $(DOC_MODULES)
.PHONY: $(ALL_MODULES)

# Module dependencies
tcr-cli: tcr-engine
tcr-doc: tcr-cli

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
tidy vet lint: $(ALL_MODULES)
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

.PHONY: build install-tools download
build install-tools download: $(BUILD_MODULES)
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
