.PHONY: default
default: build ;

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

CONFIG_PKG="github.com/murex/tcr/tcr-engine/settings"
.PHONY: build
build:
	@go build -ldflags "-s -w \
    		-X ${CONFIG_PKG}.BuildVersion="`git describe --tags`" \
            -X ${CONFIG_PKG}.BuildOs="`go env GOOS`" \
            -X ${CONFIG_PKG}.BuildArch="`go env GOARCH`" \
            -X ${CONFIG_PKG}.BuildCommit="`git rev-list --max-count=1 --tags`" \
            -X ${CONFIG_PKG}.BuildDate="`date -u +%FT%TZ`" \
            -X ${CONFIG_PKG}.BuildAuthor="`id -un`" \
            "

