# Modules that can be build into an executable
BUILD_MODULES = tcr-cli tcr-gui
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
tcr-gui: tcr-engine tcr-cli
tcr-doc: tcr-cli tcr-gui

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

.PHONY: test cov
test cov: $(PROD_MODULES)
	@for module in $^; do \
		echo "- make $@ $$module"; \
		$(MAKE) -C $$module $@; \
	done

.PHONY: build
build: $(BUILD_MODULES)
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