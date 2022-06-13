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
	@mkdir -p _test_results
	@gotestsum --format pkgname --junitfile _test_results/output.xml -- -covermode atomic -coverprofile=_test_results/coverage-all.out -tags=test_helper -v ./...
	@grep -v "_test_" _test_results/coverage-all.out >| _test_results/coverage.out

# TODO Find a way avoid replications between test and test-short targets
.PHONY: test-short
test-short:
	@mkdir -p _test_results
	@gotestsum --format pkgname --junitfile _test_results/output.xml -- -short -covermode atomic -coverprofile=_test_results/coverage-all.out -tags=test_helper -v ./...
	@grep -v "_test_" _test_results/coverage-all.out >| _test_results/coverage.out

.PHONY: cov
cov:
	@go tool cover -html=_test_results/coverage.out

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

.PHONY: download
download:
	@go mod download

.PHONY: install-tools
install-tools: download
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest