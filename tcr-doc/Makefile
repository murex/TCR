deps:
	@go get

lint:
	@golangci-lint run -v

vet:
	@go vet ./...

doc: tidy
	@echo "Generating cobra documentation"
	@go run .

tidy:
	@go mod tidy

.PHONY: deps lint vet doc tidy