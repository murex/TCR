.PHONY: default
default: build ;

.PHONY: deps
deps:
	@go get

.PHONY: lint
lint:
	@golangci-lint run --build-tags=test_helper -v

.PHONY: vet
vet:
	@go vet --tags=test_helper ./...

.PHONY: test
test:
	@RICHGO_FORCE_COLOR=1 richgo test -covermode atomic -coverprofile=coverage.out -tags=test_helper -v ./...

.PHONY: cov
cov:
	@go tool cover -html=coverage.out

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

