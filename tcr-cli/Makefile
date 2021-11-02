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

SETTINGS_PKG="github.com/murex/tcr-engine/settings"
.PHONY: build
build:
	@go build -ldflags "-s -w \
    		-X ${SETTINGS_PKG}.BuildVersion="`git describe --tags`" \
            -X ${SETTINGS_PKG}.BuildOs="`go env GOOS`" \
            -X ${SETTINGS_PKG}.BuildArch="`go env GOARCH`" \
            -X ${SETTINGS_PKG}.BuildCommit="`git rev-list --max-count=1 --tags`" \
            -X ${SETTINGS_PKG}.BuildDate="`date -u +%FT%TZ`" \
            -X ${SETTINGS_PKG}.BuildAuthor="`id -un`" \
            "

