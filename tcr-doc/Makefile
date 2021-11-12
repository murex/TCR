.PHONY: default
default: doc ;

.PHONY: deps
deps:
	@go get

.PHONY: lint
lint:
	@golangci-lint run -v

.PHONY: vet
vet:
	@go vet ./...

.PHONY: doc
doc: tidy
	@echo "Generating cobra documentation"
	@go run .

.PHONY: tidy
tidy:
	@go mod tidy
