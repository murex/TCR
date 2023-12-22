.PHONY: default
default: build ;

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
build:
	@ng build --output-path ../src/http/static/webapp --configuration production --extract-licenses false 

.PHONY: run
run:
	@ng serve -o

