.PHONY: deps
deps:
	@go get

.PHONY: lint
lint:
	@golangci-lint run -v

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test:
	@go test -coverprofile=cover.out -v ./...

.PHONY: cov
cov:
	@go tool cover -html=cover.out

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build .

.PHONY: release
release:
	@goreleaser $(GORELEASER_ARGS)

.PHONY: snapshot
snapshot: GORELEASER_ARGS= --rm-dist --snapshot
snapshot: release
