.PHONY: default
default: build ;

.PHONY: setup
setup:
	@npm install

.PHONY: deps
deps:
	@npm update

.PHONY: lint
lint:
	@echo "TODO - investigate how to run angular/typescript linter"

.PHONY: test
test:
	@echo "TODO - investigate how to run angular/typescript tests"

.PHONY: cov
cov:
	@echo "TODO - investigate how to run angular/typescript coverage"

.PHONY: build
build: setup
	@npm run build-prod

.PHONY: run
run:
	@npm start

