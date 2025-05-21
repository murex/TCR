.PHONY: default
default: doc ;

.PHONY: deps
deps:
	@go get -u -t tool ./...

.PHONY: lint
lint:
	@golangci-lint run --build-tags=test_helper -v

.PHONY: modernize
modernize:
	@go run -tags test_helper golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test -fix ./...

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
