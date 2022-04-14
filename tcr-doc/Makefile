.PHONY: default
default: doc ;

.PHONY: deps
deps:
	@go get

.PHONY: lint
lint:
	@golangci-lint run --build-tags=test_helper -v

.PHONY: vet
vet:
	@go vet --tags=test_helper ./...

.PHONY: doc
doc: tidy
	@echo "Generating cobra documentation"
	@go run .

.PHONY: tidy
tidy:
	@go mod tidy
