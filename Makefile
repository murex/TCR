APP_MODULES = tcr-cli tcr-gui
LIB_MODULES = tcr-engine
DOC_MODULES = tcr-doc

MAIN_MODULES = $(APP_MODULES) $(DOC_MODULES)
BUILD_MODULES = $(APP_MODULES) $(LIB_MODULES)

ALL_MODULES = $(APP_MODULES) $(LIB_MODULES) $(DOC_MODULES)
.PHONY: $(ALL_MODULES)

# Module dependencies
tcr-cli: tcr-engine
tcr-gui: tcr-engine tcr-cli
tcr-doc: tcr-cli tcr-gui

.PHONY: deps
deps:
	@for mod in $(MAIN_MODULES); do \
		echo "- make $@ $$mod"; \
		$(MAKE) -C $$mod $@; \
	done

.PHONY: tidy vet
tidy vet:
	@for mod in $(ALL_MODULES); do \
		echo "- make $@ $$mod"; \
		$(MAKE) -C $$mod $@; \
	done

.PHONY: test
test:
	@for mod in $(BUID_MODULES); do \
		echo "- make $@ $$mod"; \
		$(MAKE) -C $$mod $@; \
	done

.PHONY: doc
doc:
	@for mod in $(DOC_MODULES); do \
		echo "- make $@ $$mod"; \
		$(MAKE) -C $$mod $@; \
	done