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
	@go build -ldflags '-s -w \
    		-X github.com/murex/tcr-engine/settings.BuildVersion='`git describe --tags`' \
            -X github.com/murex/tcr-engine/settings.BuildOs='`go env GOOS`' \
            -X github.com/murex/tcr-engine/settings.BuildArch='`go env GOARCH`' \
            -X github.com/murex/tcr-engine/settings.BuildCommit='`git rev-list --max-count=1 --tags`' \
            -X github.com/murex/tcr-engine/settings.BuildDate='`date -u +%FT%TZ`' \
            -X github.com/murex/tcr-engine/settings.BuildAuthor='`id -un`' \
            '
