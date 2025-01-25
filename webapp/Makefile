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
	@npm run lint

.PHONY: test
test:
	@npm test

.PHONY: cov
cov:
	@npm run coverage

.PHONY: tidy
tidy:
	@npm install

.PHONY: build
build: setup
	@npm run build-prod

.PHONY: run
run:
	@npm start

.PHONY: clean
clean:
	@npm run clean

.PHONY: install-tools
install-tools: setup
